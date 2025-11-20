package model

import "time"

type Alarm struct {
	ID            uint
	DeviceId      uint
	Time          time.Time      `json:"time"`
	IsRepeat      bool           `json:"is_repeat"`
	RepeatingDays []time.Weekday `json:"days"`
}
