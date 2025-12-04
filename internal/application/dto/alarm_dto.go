package dto

import "time"

type UpdateAlarm struct {
	Time          *time.Time     `json:"time,omitempty"`
	IsRepeat      *bool          `json:"is_repeat,omitempty"`
	RepeatingDays []time.Weekday `json:"days,omitempty"`
}
