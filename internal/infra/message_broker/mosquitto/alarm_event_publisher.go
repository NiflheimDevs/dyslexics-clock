package mosquitto

import (
	"encoding/json"
	"fmt"
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
	return p.client.Publish(topic, 1, false, payload).Error()
}

func (p *AlarmEventPublisher) PublishUpdate(deviceID string, alarm *model.Alarm) error {
	topic := fmt.Sprintf("devices/%s/alarms/update", deviceID)
	payload, _ := json.Marshal(alarm)
	return p.client.Publish(topic, 1, false, payload).Error()
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
	return p.client.Publish(topic, 1, false, payload).Error()
}
