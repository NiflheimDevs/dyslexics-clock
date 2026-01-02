package service

import "github.com/NiflheimDevs/dyslexics-clock/internal/domain/model"

type PublisherService interface {
	PublishTime() error
	PublishAlarm(deviceID string, alarm *model.Alarm) error
	PublishAlarmUpdate(deviceID string, alarm *model.Alarm) error
	PublishAlarmDelete(deviceID string, alarmID string) error
	PublishVolume(deviceID string, volume uint) error
	PublishColor(deviceID string, color string) error
}
