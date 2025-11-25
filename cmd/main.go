package main

import (
	"log"

	"github.com/NiflheimDevs/dyslexics-clock/internal/mqtt"
)

func main() {
	// 1. Connect to MQTT
	mqttClient, err := mqtt.New("tcp://localhost:1883", "backend-server")
	if err != nil {
		log.Fatal("MQTT connect error:", err)
	}

	// 2. Init services
	//deviceService := device.New()
	//alarmService := alarm.New()

	// 3. Register MQTT subscribers
	sub := mqtt.NewSubscriber(mqttClient /* alarmService, deviceService */)
	sub.RegisterHandlers()

	select {}
}
