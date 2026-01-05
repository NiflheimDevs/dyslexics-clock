package mosquitto

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/NiflheimDevs/dyslexics-clock/internal/domain/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type RequestPublisher struct {
	client mqtt.Client
	qos    byte
}

func NewRequestPublisher(client mqtt.Client) *RequestPublisher {
	return &RequestPublisher{
		client: client,
		qos:    1, // at-least-once delivery
	}
}

func (p *RequestPublisher) PublishAlarms(deviceID string, alarms []model.Alarm) error {
	topic := fmt.Sprintf("devices/%s/alarms", deviceID)

	payload, err := json.Marshal(alarms)
	if err != nil {
		log.Printf("Error marshalling alarms for device %s: %v", deviceID, err)
		return err
	}

	token := p.client.Publish(topic, p.qos, false, payload)
	token.Wait() // Wait for the publication to complete
	if token.Error() != nil {
		log.Printf("Error publishing alarm list for device %s to topic %s: %v", deviceID, topic, token.Error())
	} else {
		log.Printf("Published %d alarms for device %s to topic %s", len(alarms), deviceID, topic)
	}

	return token.Error()
}

func (p *RequestPublisher) PublishColor(deviceID string, color string) error {
	topic := fmt.Sprintf("devices/%s/color", deviceID)

	token := p.client.Publish(topic, p.qos, false, color)
	token.Wait() // Wait for the publication to complete
	if token.Error() != nil {
		log.Printf("Error publishing color for device %s to topic %s: %v", deviceID, topic, token.Error())
	} else {
		log.Printf("Published %s color for device %s to topic %s", color, deviceID, topic)
	}
	return token.Error()
}

func (p *RequestPublisher) PublishVolume(deviceID string, volume uint) error {
	topic := fmt.Sprintf("devices/%s/volume", deviceID)

	token := p.client.Publish(topic, p.qos, false, fmt.Append(nil, volume))
	token.Wait() // Wait for the publication to complete
	if token.Error() != nil {
		log.Printf("Error publishing volume for device %s to topic %s: %v", deviceID, topic, token.Error())
	} else {
		log.Printf("Published %dvolume for device %s to topic %s", volume, deviceID, topic)
	}
	return token.Error()
}
	
func (p *RequestPublisher) PublishBrightness(deviceID string, brightness uint) error {
	topic := fmt.Sprintf("devices/%s/volume", deviceID)

	token := p.client.Publish(topic, p.qos, false, fmt.Append(nil, brightness))
	token.Wait() // Wait for the publication to complete
	if token.Error() != nil {
		log.Printf("Error publishing brightness for device %s to topic %s: %v", deviceID, topic, token.Error())
	} else {
		log.Printf("Published %d brightness for device %s to topic %s", brightness, deviceID, topic)
	}
	return token.Error()
}
	
