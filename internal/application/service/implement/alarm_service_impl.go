package serviceimpl

import (
	"context"
	"fmt"

	"github.com/NiflheimDevs/dyslexics-clock/internal/application/dto"
	"github.com/NiflheimDevs/dyslexics-clock/internal/application/service"
	derror "github.com/NiflheimDevs/dyslexics-clock/internal/domain/error"
	"github.com/NiflheimDevs/dyslexics-clock/internal/domain/model"
	repository "github.com/NiflheimDevs/dyslexics-clock/internal/domain/repository"
)

type AlarmService struct {
	AlarmRepo repository.AlarmRepo
	Publisher service.PublisherService
}

func NewAlarmService(AlarmRepo repository.AlarmRepo, p service.PublisherService) *AlarmService {
	return &AlarmService{
		AlarmRepo: AlarmRepo,
		Publisher: p,
	}
}

func (a *AlarmService) GetAlarms(ctx context.Context, DeviceId uint) ([]model.Alarm, error) {
	return a.AlarmRepo.GetAlarms(ctx, DeviceId)
}

func (a *AlarmService) InsertAlarm(ctx context.Context, alarm *model.Alarm) error {
	err := a.AlarmRepo.InsertAlarm(ctx, alarm)
	if err != nil {
		return err
	}

	go a.Publisher.PublishAlarm(fmt.Sprint(alarm.DeviceId), alarm)
	return nil
}

func (a *AlarmService) DeleteAlarmById(ctx context.Context, alarmID uint, deviceID uint) error {
	rowsAffected, err := a.AlarmRepo.DeleteAlarmById(ctx, alarmID, deviceID)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return derror.New(derror.ErrTypeNotFound, "alarm not found", nil)
	}

	go a.Publisher.PublishAlarmDelete(fmt.Sprint(deviceID), fmt.Sprint(alarmID))
	return nil
}

func (a *AlarmService) UpdateAlarm(ctx context.Context, alarmID uint, deviceID uint, updateAlarm *dto.UpdateAlarm) (*model.Alarm, error) {
	alarm, err := a.AlarmRepo.UpdateAlarm(ctx, alarmID, deviceID, updateAlarm)
	if err != nil {
		return nil, err
	}
	if alarm == nil {
		return nil, derror.New(derror.ErrTypeNotFound, "alarm not found", nil)
	}

	go a.Publisher.PublishAlarmUpdate(fmt.Sprint(deviceID), alarm)
	return alarm, nil
}
