package service

import (
	"context"

	"github.com/NiflheimDevs/dyslexics-clock/internal/application/dto"
	"github.com/NiflheimDevs/dyslexics-clock/internal/domain/model"
)

type AlarmService interface {
	GetAlarms(ctx context.Context, DeviceId uint) ([]model.Alarm, error)
	InsertAlarm(ctx context.Context, alarm *model.Alarm) error
	DeleteAlarmById(ctx context.Context, alarmID uint, deviceID uint) error
	UpdateAlarm(ctx context.Context, alarmID uint, deviceID uint ,updateAlarm *dto.UpdateAlarm) error
}
