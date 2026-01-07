package messagebroker

type TimePublisher interface {
	PublishTime(time int64) error
}
