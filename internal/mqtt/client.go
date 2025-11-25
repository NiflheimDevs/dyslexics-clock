package mqtt

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	Client mqtt.Client
}

func New(broker string, clientID string) (*Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetAutoReconnect(true)
	opts.SetConnectTimeout(5 * time.Second)

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("mqtt connection error: %v", token.Error())
	}

	return &Client{Client: client}, nil
}

func (m *Client) Publish(topic string, payload []byte) error {
	token := m.Client.Publish(topic, 0, false, payload)
	token.Wait()
	return token.Error()
}

func (m *Client) Subscribe(topic string, handler mqtt.MessageHandler) error {
	token := m.Client.Subscribe(topic, 0, handler)
	token.Wait()
	return token.Error()
}

func (m *Client) Close() {
	m.Client.Disconnect(250)
}
