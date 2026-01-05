package repository

import (
	"context"

	"github.com/NiflheimDevs/dyslexics-clock/internal/application/dto"
	"github.com/NiflheimDevs/dyslexics-clock/internal/domain/model"
)

type DeviceRepo interface {
	GetDeviceByUsername(ctx context.Context, username string) (*dto.LoginDto, error)
	GetDeviceById(ctx context.Context, Id uint) (*model.Device, error)
	UpdateColor(ctx context.Context, Id uint, color string) error
	UpdateVolume(ctx context.Context, Id uint, volume uint) error
	UpdateBrightness(ctx context.Context, Id uint, brightness uint) error
}
