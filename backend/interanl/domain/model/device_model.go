package model

type Device struct {
	ID       uint
	Password []byte
	Username string
	Alarms   []Alarm
	Color    string
	Timezone string
}
