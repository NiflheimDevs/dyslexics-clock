package messagebroker

import "github.com/NiflheimDevs/dyslexics-clock/internal/domain/model"

type AlarmEventPublisher interface {
	PublishCreate(deviceID string, alarm *model.Alarm) error
	PublishUpdate(deviceID string, alarm *model.Alarm) error
	PublishDelete(deviceID string, alarmID string) error
	PublishRing(deviceID string) error
	PublishSnooze(deviceID string) error
	PublishSilent(deviceID string) error
}
