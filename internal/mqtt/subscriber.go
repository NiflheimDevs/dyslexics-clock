package mqtt

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/NiflheimDevs/dyslexics-clock/internal/domain/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Subscriber struct {
	Client *Client
	// AlarmService  *alarm.Service
	// DeviceService *device.Service
}

func NewSubscriber(m *Client /* a *alarm.Service, d *device.Service */) *Subscriber {
	return &Subscriber{Client: m /* AlarmService: a, DeviceService: d*/}
}

func (s *Subscriber) RegisterHandlers() {
	s.Client.Subscribe("device/+/color/update", s.onColorUpdate)
	s.Client.Subscribe("device/+/alarm/+", s.onAlarmEvent)

	log.Println("MQTT subscriptions registered")
}

func (s *Subscriber) onColorUpdate(c mqtt.Client, m mqtt.Message) {
	parts := strings.Split(m.Topic(), "/")
	deviceID := parts[1]
	color := string(m.Payload())

	log.Println("Color update received:", deviceID, color)

	// save to DB
	//id, _ := strconv.Atoi(deviceID)
	//s.DeviceService.UpdateColor(uint(id), color)

	// forward to device
	s.Client.Publish("server/"+deviceID+"/color", []byte(color))
}

func (s *Subscriber) onAlarmEvent(c mqtt.Client, m mqtt.Message) {
	parts := strings.Split(m.Topic(), "/")
	deviceID := parts[1]
	action := parts[3]

	log.Println("Alarm event:", deviceID, action)

	switch action {

	case "add":
		var a model.Alarm
		json.Unmarshal(m.Payload(), &a)
		id, _ := strconv.Atoi(deviceID)
		a.DeviceId = uint(id)

		//s.AlarmService.Create(a)
		s.Client.Publish("server/"+deviceID+"/alarm/add", m.Payload())

	case "update":
		var a model.Alarm
		json.Unmarshal(m.Payload(), &a)
		//s.AlarmService.Update(a)
		s.Client.Publish("server/"+deviceID+"/alarm/update", m.Payload())

	case "delete":
		var data struct{ ID uint }
		json.Unmarshal(m.Payload(), &data)
		//s.AlarmService.Delete(data.ID)
		s.Client.Publish("server/"+deviceID+"/alarm/delete", m.Payload())
	}
}
