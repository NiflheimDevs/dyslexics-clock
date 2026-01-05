package model

type Device struct {
	ID         uint
	Volume     uint
	Password   []byte
	Username   string
	Color      string
	Brightness uint
	Timezone   string
}
