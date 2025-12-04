package model

import "time"

type Alarm struct {
	ID            uint			 `json:"id"`
	DeviceId      uint			 `json:"-"`
	Time          time.Time      `json:"time"`
	IsRepeat      bool           `json:"is_repeat"`
	RepeatingDays []time.Weekday `json:"days"`
}
