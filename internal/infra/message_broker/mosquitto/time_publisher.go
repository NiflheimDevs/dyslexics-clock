package mosquitto

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type TimePublisher struct {
	client mqtt.Client
	qos    byte
}

func NewTimePublisher(client mqtt.Client) *TimePublisher {
	return &TimePublisher{
		client: client,
		qos:    0, 
	}
}

func (p *TimePublisher) PublishTime(time int64) error {
	topic := "devices/time"
	return p.client.Publish(topic, 1, false, time).Error()
}
