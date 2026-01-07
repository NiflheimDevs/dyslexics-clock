package model

import "time"

type Device struct {
	ID         uint
	Volume     uint
	Password   []byte
	Username   string
	Color      string
	Birthdate  *time.Time
	Brightness uint
	Timezone   string
}
