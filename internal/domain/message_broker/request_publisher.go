package messagebroker

import "github.com/NiflheimDevs/dyslexics-clock/internal/domain/model"

type RequestPublisher interface {
	PublishAlarms(deviceID string, alarms []model.Alarm) error
	PublishColor(deviceID string, color string) error
	PublishVolume(deviceID string, volume uint) error
	PublishBrightness(deviceID string, brightness uint) error
}
