package messagebroker

import (
	"time"

	"github.com/NiflheimDevs/dyslexics-clock/internal/domain/model"
)

type RequestPublisher interface {
	PublishAlarms(deviceID string, alarms []model.Alarm) error
	PublishColor(deviceID string, color string) error
	PublishVolume(deviceID string, volume uint) error
	PublishBrightness(deviceID string, brightness uint) error
	PublishBirthdate(deviceID string, birthdate time.Time) error
}
