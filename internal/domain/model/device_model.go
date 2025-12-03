package model

type Device struct {
	ID       uint
	Password []byte
	Username string
	Color    string
	Timezone string
}
