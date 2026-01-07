package messagebroker

type DeviceEventSubscriber interface {
	Subscribe() error
}
