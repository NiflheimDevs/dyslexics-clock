package mosquitto

import (
	"context"
	"log"
	"time"

	"github.com/NiflheimDevs/dyslexics-clock/internal/application/service"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const deviceGetAllAlarmsTopic = "devices/alarms" // Topic where devices post their IDs for getting all alarms
const deviceGetColorTopic = "devices/color"      // Topic where devices post their IDs for getting their color
const deviceGetVolumeTopic = "devices/volume"    // Topic where devices post their IDs for getting their volume
const deviceGetBrightness = "devices/brightness"
const deviceGetBirthdate = "devices/birthdate"

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
	log.Printf("Subscribing to MQTT topic: %s", deviceGetAllAlarmsTopic)

	token := s.client.Subscribe(deviceGetAllAlarmsTopic, s.qos, func(client mqtt.Client, msg mqtt.Message) {
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

	log.Printf("Successfully subscribed to topic: %s", deviceGetAllAlarmsTopic)

	log.Printf("Subscribing to MQTT topic: %s", deviceGetColorTopic)

	token = s.client.Subscribe(deviceGetColorTopic, s.qos, func(client mqtt.Client, msg mqtt.Message) {
		deviceID := string(msg.Payload())
		log.Printf("Received message on topic '%s': Device ID '%s'", msg.Topic(), deviceID)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Set a timeout for handling the message
		defer cancel()

		if err := s.deviceEventService.HandleGetColorMessage(ctx, deviceID); err != nil {
			log.Printf("Error handling device message for device ID '%s': %v", deviceID, err)
		}
	})

	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}

	log.Printf("Successfully subscribed to topic: %s", deviceGetColorTopic)

	log.Printf("Subscribing to MQTT topic: %s", deviceGetVolumeTopic)

	token = s.client.Subscribe(deviceGetVolumeTopic, s.qos, func(client mqtt.Client, msg mqtt.Message) {
		deviceID := string(msg.Payload())
		log.Printf("Received message on topic '%s': Device ID '%s'", msg.Topic(), deviceID)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Set a timeout for handling the message
		defer cancel()

		if err := s.deviceEventService.HandleGetVolumeMessage(ctx, deviceID); err != nil {
			log.Printf("Error handling device message for device ID '%s': %v", deviceID, err)
		}
	})

	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}

	log.Printf("Successfully subscribed to topic: %s", deviceGetVolumeTopic)

	log.Printf("Subscribing to MQTT topic: %s", deviceGetBrightness)

	token = s.client.Subscribe(deviceGetBrightness, s.qos, func(client mqtt.Client, msg mqtt.Message) {
		deviceID := string(msg.Payload())
		log.Printf("Received message on topic '%s': Device ID '%s'", msg.Topic(), deviceID)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Set a timeout for handling the message
		defer cancel()

		if err := s.deviceEventService.HandleGetBrightnessMessage(ctx, deviceID); err != nil {
			log.Printf("Error handling device message for device ID '%s': %v", deviceID, err)
		}
	})

	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}

	log.Printf("Successfully subscribed to topic: %s", deviceGetBrightness)
	log.Printf("Subscribing to MQTT topic: %s", deviceGetBirthdate)

	token = s.client.Subscribe(deviceGetBirthdate, s.qos, func(client mqtt.Client, msg mqtt.Message) {
		deviceID := string(msg.Payload())
		log.Printf("Received message on topic '%s': Device ID '%s'", msg.Topic(), deviceID)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Set a timeout for handling the message
		defer cancel()

		if err := s.deviceEventService.HandleGetBirthdateMessage(ctx, deviceID); err != nil {
			log.Printf("Error handling device message for device ID '%s': %v", deviceID, err)
		}
	})

	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}

	log.Printf("Successfully subscribed to topic: %s", deviceGetBirthdate)
	return nil
}
