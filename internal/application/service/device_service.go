package service

import (
	"context"
)

type DeviceService interface {
	Login(ctx context.Context, username string, password string) (string, error)
	GetDeviceColor(ctx context.Context, id uint) (string, error)
	UpdateDeviceColor(ctx context.Context, id uint, newColor string) error
}
