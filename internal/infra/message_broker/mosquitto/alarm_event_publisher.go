package mosquitto

import (
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
	return p.client.Publish(topic, 1, false, alarm).Error()
}

func (p *AlarmEventPublisher) PublishUpdate(deviceID string, alarm *model.Alarm) error {
	topic := fmt.Sprintf("devices/%s/alarms/update", deviceID)
	return p.client.Publish(topic, 1, false, alarm).Error()
}

func (p *AlarmEventPublisher) PublishDelete(deviceID string, alarmID string) error {
	topic := fmt.Sprintf("devices/%s/alarms/delete", deviceID)

	payload := struct {
		ID        string    `json:"id"`
		DeletedAt time.Time `json:"deleted_at"`
	}{
		ID:        alarmID,
		DeletedAt: time.Now().UTC(),
	}

	return p.client.Publish(topic, 1, false, payload).Error()
}
