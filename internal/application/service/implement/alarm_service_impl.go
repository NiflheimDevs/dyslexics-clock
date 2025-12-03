package serviceimpl

import (
	"context"

	"github.com/NiflheimDevs/dyslexics-clock/internal/application/dto"
	derror "github.com/NiflheimDevs/dyslexics-clock/internal/domain/error"
	"github.com/NiflheimDevs/dyslexics-clock/internal/domain/model"
	repository "github.com/NiflheimDevs/dyslexics-clock/internal/domain/repository/postgres"
)

type AlarmService struct {
	AlarmRepo repository.AlarmRepo
}

func NewAlarmService(AlarmRepo repository.AlarmRepo) *AlarmService {
	return &AlarmService{
		AlarmRepo: AlarmRepo,
	}
}

func (a *AlarmService) GetAlarms(ctx context.Context, DeviceId uint) ([]model.Alarm, error) {
	return a.AlarmRepo.GetAlarms(ctx, DeviceId)
}

func (a *AlarmService) InsertAlarm(ctx context.Context, alarm *model.Alarm) error {
	return a.AlarmRepo.InsertAlarm(ctx, alarm)
}

func (a *AlarmService) DeleteAlarmById(ctx context.Context, alarmID uint, deviceID uint) error {
	rowsAffected, err := a.AlarmRepo.DeleteAlarmById(ctx, alarmID, deviceID)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return derror.New(derror.ErrTypeNotFound, "alarm not found", nil)
	}
	return nil
}

func (a *AlarmService) UpdateAlarm(ctx context.Context, alarmID uint, deviceID uint, updateAlarm *dto.UpdateAlarm) error {
	rowsAffected, err := a.AlarmRepo.UpdateAlarm(ctx, alarmID, deviceID, updateAlarm)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return derror.New(derror.ErrTypeNotFound, "alarm not found", nil)
	}
	return nil
}
