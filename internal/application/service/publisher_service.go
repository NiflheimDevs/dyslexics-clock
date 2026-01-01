package service

import "github.com/NiflheimDevs/dyslexics-clock/internal/domain/model"

type PublisherService interface {
	PublishTime() error
	PublishAlarm(deviceID uint, alarm *model.Alarm) error
	PublishAlarmUpdate(deviceID uint, alarm *model.Alarm) error
	PublishAlarmDelete(deviceID string, alarmID string) error
}
