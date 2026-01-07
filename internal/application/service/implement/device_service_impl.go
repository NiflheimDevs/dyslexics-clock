package serviceimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/NiflheimDevs/dyslexics-clock/internal/application/service"
	derror "github.com/NiflheimDevs/dyslexics-clock/internal/domain/error"
	"github.com/NiflheimDevs/dyslexics-clock/internal/domain/model"
	"github.com/NiflheimDevs/dyslexics-clock/internal/domain/pkg"
	repository "github.com/NiflheimDevs/dyslexics-clock/internal/domain/repository"
)

type DeviceService struct {
	SecretSauce      pkg.SecretSauce
	DeviceRepo       repository.DeviceRepo
	PublisherService service.PublisherService
	JWTService       service.JWT
}

func NewDeviceService(deviceRepo repository.DeviceRepo,
	secretSauce pkg.SecretSauce,
	jwtService service.JWT,
	p service.PublisherService,
) *DeviceService {
	return &DeviceService{
		DeviceRepo:       deviceRepo,
		SecretSauce:      secretSauce,
		JWTService:       jwtService,
		PublisherService: p,
	}
}

func (d *DeviceService) Login(ctx context.Context, username string, password string) (string, error) {
	deviceInfo, err := d.DeviceRepo.GetDeviceByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if err := d.SecretSauce.SauceReferee(deviceInfo.Password, password); err != nil {
		panic(derror.New(derror.ErrTypeNotFound, "user not found", err))
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

func (d *DeviceService) GetDeviceById(ctx context.Context, id uint) (*model.Device, error) {
	return d.DeviceRepo.GetDeviceById(ctx, id)
}

func (d *DeviceService) UpdateDeviceColor(ctx context.Context, id uint, newColor string) error {
	err := d.DeviceRepo.UpdateColor(ctx, id, newColor)
	if err != nil {
		return err
	}
	go d.PublisherService.PublishColor(fmt.Sprint(id), newColor)
	return nil
}

func (d *DeviceService) UpdateDeviceVolume(ctx context.Context, id uint, newVolume uint) error {
	err := d.DeviceRepo.UpdateVolume(ctx, id, newVolume)
	if err != nil {
		return err
	}
	go d.PublisherService.PublishVolume(fmt.Sprint(id), newVolume)
	return nil
}

func (d *DeviceService) UpdateDeviceBrightness(ctx context.Context, id uint, newBrightness uint) error {
	err := d.DeviceRepo.UpdateBrightness(ctx, id, newBrightness)
	if err != nil {
		return err
	}
	go d.PublisherService.PublishBrightness(fmt.Sprint(id), newBrightness)
	return nil
}

func (d *DeviceService) UpdateDeviceBirthdate(ctx context.Context, id uint, newBirthdate time.Time) error {
	err := d.DeviceRepo.UpdateBirthdate(ctx, id, newBirthdate)
	if err != nil {
		return err
	}
	go d.PublisherService.PublishBirthdate(fmt.Sprint(id), newBirthdate)
	return nil
}
