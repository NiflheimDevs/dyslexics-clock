package mosquitto

import (
	"fmt"
	"log"

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
	token := p.client.Publish(topic, 1, false, fmt.Append(nil, time))
	token.Wait()
	err := token.Error()
	if err != nil {
		log.Printf("Error publishing time to topic %s: %v", topic, err)
	}
	return err
}
