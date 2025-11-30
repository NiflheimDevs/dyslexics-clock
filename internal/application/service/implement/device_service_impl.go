package serviceimpl

import (
	"context"

	"github.com/NiflheimDevs/dyslexics-clock/internal/application/service"
	derror "github.com/NiflheimDevs/dyslexics-clock/internal/domain/error"
	"github.com/NiflheimDevs/dyslexics-clock/internal/domain/pkg"
	repository "github.com/NiflheimDevs/dyslexics-clock/internal/domain/repository/postgres"
)

type DeviceService struct {
	SecretSauce pkg.SecretSauce
	DeviceRepo  repository.DeviceRepo
	JWTService  service.JWT
}

func NewDeviceService(deviceRepo repository.DeviceRepo,
	secretSauce pkg.SecretSauce,
	jwtService service.JWT,
) *DeviceService {
	return &DeviceService{
		DeviceRepo:  deviceRepo,
		SecretSauce: secretSauce,
		JWTService:  jwtService,
	}
}

func (d *DeviceService) Login(ctx context.Context, username string, password string) (string, error) {
	deviceInfo, err := d.DeviceRepo.GetDeviceByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if err := d.SecretSauce.SauceReferee(deviceInfo.Password, password); err != nil {
		panic(derror.New(derror.ErrTypeUnauthorized, "invalid password", err))
	}
	token, _ := d.JWTService.GenerateToken(deviceInfo.Id)
	return token, nil
}

func (d *DeviceService) GetDeviceColor(ctx context.Context, id uint) (string, error) {
	device, err := d.DeviceRepo.GetDeviceById(ctx, id)
	if err != nil {
		return "", err 
	}
	return device.Color, nil
}

func (d *DeviceService) UpdateDeviceColor(ctx context.Context, id uint, newColor string) error {
	return d.DeviceRepo.UpdateColor(ctx, id, newColor)
}
