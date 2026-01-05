package mosquitto

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/NiflheimDevs/dyslexics-clock/internal/domain/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type AlarmEventPublisher struct {
	client mqtt.Client
	qos    byte
}

func NewAlarmEventPublisher(client mqtt.Client) *AlarmEventPublisher {
	return &AlarmEventPublisher{
		client: client,
		qos:    1, // at-least-once delivery
	}
}

func (p *AlarmEventPublisher) PublishCreate(deviceID string, alarm *model.Alarm) error {
	topic := fmt.Sprintf("devices/%s/alarms/create", deviceID)
	payload, _ := json.Marshal(alarm)
	token := p.client.Publish(topic, 1, false, payload)
	token.Wait()
	err := token.Error()
	if err != nil {
		log.Printf("Error publishing create for device %s to topic %s: %v", deviceID, topic, err)
	}
	return err
}

func (p *AlarmEventPublisher) PublishUpdate(deviceID string, alarm *model.Alarm) error {
	topic := fmt.Sprintf("devices/%s/alarms/update", deviceID)
	payload, _ := json.Marshal(alarm)
	token := p.client.Publish(topic, 1, false, payload)
	token.Wait()
	err := token.Error()
	if err != nil {
		log.Printf("Error publishing update for device %s to topic %s: %v", deviceID, topic, err)
	}
	return err
}

func (p *AlarmEventPublisher) PublishDelete(deviceID string, alarmID string) error {
	topic := fmt.Sprintf("devices/%s/alarms/delete", deviceID)

	payloadStruct := struct {
		ID        string    `json:"id"`
		DeletedAt time.Time `json:"deleted_at"`
	}{
		ID:        alarmID,
		DeletedAt: time.Now().UTC(),
	}

	payload, _ := json.Marshal(payloadStruct)
	token := p.client.Publish(topic, 1, false, payload)
	token.Wait()
	err := token.Error()
	if err != nil {
		log.Printf("Error publishing delete for device %s to topic %s: %v", deviceID, topic, err)
	}
	return err
}

func (p *AlarmEventPublisher) PublishRing(deviceID string) error {
	topic := fmt.Sprintf("devices/%s/ring", deviceID)
	token := p.client.Publish(topic, 1, false, nil)
	token.Wait()
	err := token.Error()
	if err != nil {
		log.Printf("Error publishing ring for device %s to topic %s: %v", deviceID, topic, err)
	}
	return err
}

func (p *AlarmEventPublisher) PublishSnooze(deviceID string) error {
	topic := fmt.Sprintf("devices/%s/snooze", deviceID)
	token := p.client.Publish(topic, 1, false, nil)
	token.Wait()
	err := token.Error()
	if err != nil {
		log.Printf("Error publishing silent for device %s to topic %s: %v", deviceID, topic, err)
	}
	return err
}

func (p *AlarmEventPublisher) PublishSilent(deviceID string) error {
	topic := fmt.Sprintf("devices/%s/silent", deviceID)
	token := p.client.Publish(topic, 1, false, nil)
	token.Wait()
	err := token.Error()
	if err != nil {
		log.Printf("Error publishing silent for device %s to topic %s: %v", deviceID, topic, err)
	}
	return err
}
