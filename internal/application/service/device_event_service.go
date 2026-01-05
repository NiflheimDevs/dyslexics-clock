package service

import "context"

type DeviceEventService interface {
	HandleGetAlarmsMessage(ctx context.Context, deviceID string) error
	HandleGetColorMessage(ctx context.Context, deviceID string) error
	HandleGetVolumeMessage(ctx context.Context, deviceID string) error
	HandleGetBrightnessMessage(ctx context.Context, deviceID string) error
}
