package dto

import "time"

type UpdateAlarm struct {
	Time          time.Time      `json:"time"`
	IsRepeat      bool           `json:"is_repeat"`
	RepeatingDays []time.Weekday `json:"days"`
}
