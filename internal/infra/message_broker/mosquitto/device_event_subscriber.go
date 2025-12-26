package mosquitto

import (
	"context"
	"log"
	"time"

	"github.com/NiflheimDevs/dyslexics-clock/internal/application/service"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const deviceRegistrationTopic = "devices/register" // Topic where devices post their IDs

type DeviceEventSubscriber struct {
	client             mqtt.Client
	deviceEventService service.DeviceEventService
	qos                byte
}

func NewDeviceEventSubscriber(client mqtt.Client, deviceEventService service.DeviceEventService) *DeviceEventSubscriber {
	return &DeviceEventSubscriber{
		client:             client,
		deviceEventService: deviceEventService,
		qos:                1, // at-least-once delivery
	}
}

func (s *DeviceEventSubscriber) Subscribe() error {
	log.Printf("Subscribing to MQTT topic: %s", deviceRegistrationTopic)

	token := s.client.Subscribe(deviceRegistrationTopic, s.qos, func(client mqtt.Client, msg mqtt.Message) {
		deviceID := string(msg.Payload())
		log.Printf("Received message on topic '%s': Device ID '%s'", msg.Topic(), deviceID)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Set a timeout for handling the message
		defer cancel()

		if err := s.deviceEventService.HandleGetAlarmsMessage(ctx, deviceID); err != nil {
			log.Printf("Error handling device message for device ID '%s': %v", deviceID, err)
		}
	})

	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}

	log.Printf("Successfully subscribed to topic: %s", deviceRegistrationTopic)
	return nil
}
