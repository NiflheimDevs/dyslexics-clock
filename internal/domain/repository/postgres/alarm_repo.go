package repository

import (
	"context"

	"github.com/NiflheimDevs/dyslexics-clock/internal/application/dto"
	"github.com/NiflheimDevs/dyslexics-clock/internal/domain/model"
)

type AlarmRepo interface {
	InsertAlarm(ctx context.Context, alarm *model.Alarm) error
	DeleteAlarmById(ctx context.Context, alarmID uint, deviceID uint) (int64, error)
	UpdateAlarm(ctx context.Context, alarmID uint, deviceID uint ,updateAlarm *dto.UpdateAlarm) (int64, error)
	GetAlarms(ctx context.Context, DeviceId uint) ([]model.Alarm, error)
	GetAlarmById(ctx context.Context, DeviceId uint) (*model.Alarm, error)
}
