package driver

import (
	"fmt"
	"log"

	"github.com/NiflheimDevs/dyslexics-clock/bootstrap"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// TODO: 
func ConnectMosquitto(di *bootstrap.Di) mqtt.Client {
	client := mqtt.NewClient(mqtt.NewClientOptions().
		SetClientID(di.Env.MQTT.ClientID).
		AddBroker(fmt.Sprintf("tcp://%s:%s", di.Env.MQTT.Address, di.Env.MQTT.Port)).
		SetAutoReconnect(true).
		SetCleanSession(true).SetOnConnectHandler(func(c mqtt.Client) {
		log.Println("MQTT connected")
	}).
		SetConnectionLostHandler(func(c mqtt.Client, err error) {
			log.Printf("MQTT connection lost: %v", err)
		}))
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client
}
