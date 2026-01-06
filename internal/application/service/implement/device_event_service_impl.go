package serviceimpl

import (
	"context"
	"log"
	"strconv"

	messagebroker "github.com/NiflheimDevs/dyslexics-clock/internal/domain/message_broker"
	"github.com/NiflheimDevs/dyslexics-clock/internal/domain/repository"
)

type DeviceEventService struct {
	alarmRepo        repository.AlarmRepo
	deviceRepo       repository.DeviceRepo
	requestPublisher messagebroker.RequestPublisher
}

func NewDeviceEventService(alarmRepo repository.AlarmRepo, requestPublisher messagebroker.RequestPublisher, deviceRepo repository.DeviceRepo) *DeviceEventService {
	return &DeviceEventService{
		alarmRepo:        alarmRepo,
		deviceRepo:       deviceRepo,
		requestPublisher: requestPublisher,
	}
}

func (s *DeviceEventService) HandleGetAlarmsMessage(ctx context.Context, deviceID string) error {
	log.Printf("Received get alarm device message for device ID: %s", deviceID)

	// Convert deviceID string to uint
	id, err := strconv.ParseUint(deviceID, 10, 64)
	if err != nil {
		log.Printf("Error converting device ID '%s' to uint: %v", deviceID, err)
		return err // Or a more specific error type
	}
	deviceUintID := uint(id)

	alarms, err := s.alarmRepo.GetAlarms(ctx, deviceUintID)
	if err != nil {
		log.Printf("Error getting alarms for device %s: %v", deviceID, err)
		return err
	}

	if len(alarms) == 0 {
		log.Printf("No alarms found for device %s. Publishing empty list.", deviceID)
	} else {
		log.Printf("Found %d alarms for device %s. Publishing list.", len(alarms), deviceID)
	}

	return s.requestPublisher.PublishAlarms(deviceID, alarms)
}

func (s *DeviceEventService) HandleGetColorMessage(ctx context.Context, deviceID string) error {
	log.Printf("Received get color device message for device ID: %s", deviceID)

	// Convert deviceID string to uint
	id, err := strconv.ParseUint(deviceID, 10, 64)
	if err != nil {
		log.Printf("Error converting device ID '%s' to uint: %v", deviceID, err)
		return err // Or a more specific error type
	}
	deviceUintID := uint(id)

	device, err := s.deviceRepo.GetDeviceById(ctx, deviceUintID)
	if err != nil {
		log.Printf("Error getting alarms for device %s: %v", deviceID, err)
		return err
	}

	return s.requestPublisher.PublishColor(deviceID, device.Color)
}

func (s *DeviceEventService) HandleGetVolumeMessage(ctx context.Context, deviceID string) error {
	log.Printf("Received get volume device message for device ID: %s", deviceID)

	// Convert deviceID string to uint
	id, err := strconv.ParseUint(deviceID, 10, 64)
	if err != nil {
		log.Printf("Error converting device ID '%s' to uint: %v", deviceID, err)
		return err // Or a more specific error type
	}
	deviceUintID := uint(id)

	device, err := s.deviceRepo.GetDeviceById(ctx, deviceUintID)
	if err != nil {
		log.Printf("Error getting alarms for device %s: %v", deviceID, err)
		return err
	}

	return s.requestPublisher.PublishVolume(deviceID, device.Volume)
}

func (s *DeviceEventService) HandleGetBrightnessMessage(ctx context.Context, deviceID string) error {
	log.Printf("Received get volume device message for device ID: %s", deviceID)

	// Convert deviceID string to uint
	id, err := strconv.ParseUint(deviceID, 10, 64)
	if err != nil {
		log.Printf("Error converting device ID '%s' to uint: %v", deviceID, err)
		return err // Or a more specific error type
	}
	deviceUintID := uint(id)

	device, err := s.deviceRepo.GetDeviceById(ctx, deviceUintID)
	if err != nil {
		log.Printf("Error getting alarms for device %s: %v", deviceID, err)
		return err
	}

	return s.requestPublisher.PublishBrightness(deviceID, device.Brightness)
}

func (s *DeviceEventService) HandleGetBirthdateMessage(ctx context.Context, deviceID string) error {
	log.Printf("Received get volume device message for device ID: %s", deviceID)

	// Convert deviceID string to uint
	id, err := strconv.ParseUint(deviceID, 10, 64)
	if err != nil {
		log.Printf("Error converting device ID '%s' to uint: %v", deviceID, err)
		return err // Or a more specific error type
	}
	deviceUintID := uint(id)

	device, err := s.deviceRepo.GetDeviceById(ctx, deviceUintID)
	if err != nil {
		log.Printf("Error getting alarms for device %s: %v", deviceID, err)
		return err
	}
	if device.Birthdate == nil {
		return nil
	}
	return s.requestPublisher.PublishBirthdate(deviceID, *device.Birthdate)
}
