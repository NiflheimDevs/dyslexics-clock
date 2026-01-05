package service

import (
	"context"

	"github.com/NiflheimDevs/dyslexics-clock/internal/domain/model"
)

type DeviceService interface {
	Login(ctx context.Context, username string, password string) (string, error)
	GetDeviceColor(ctx context.Context, id uint) (string, error)
	UpdateDeviceColor(ctx context.Context, id uint, newColor string) error
	GetDeviceById(ctx context.Context, id uint) (*model.Device, error)
	UpdateDeviceVolume(ctx context.Context, id uint, newVolume uint) error
	UpdateDeviceBrightness(ctx context.Context, id uint, newBrightness uint) error
}
