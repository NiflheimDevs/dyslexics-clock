package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/NiflheimDevs/dyslexics-clock/bootstrap"
	"github.com/NiflheimDevs/dyslexics-clock/internal/application/service"
	"github.com/NiflheimDevs/dyslexics-clock/internal/domain/pkg"
)

type DeviceHandler struct {
	Constants        *bootstrap.Constants
	DeviceService    service.DeviceService
	PublisherService service.PublisherService
	Validator        pkg.Validator
}

func NewDeviceHandler(
	constants *bootstrap.Constants,
	deviceService service.DeviceService,
	validator pkg.Validator,
	publisherService service.PublisherService,
) *DeviceHandler {
	return &DeviceHandler{
		Constants:        constants,
		DeviceService:    deviceService,
		Validator:        validator,
		PublisherService: publisherService,
	}
}

func (dh *DeviceHandler) Login(w http.ResponseWriter, r *http.Request) {

	type LoginRequest struct {
		// DeviceID uuid.UUID `json:"device_id" validator:"required"`
		Username string `json:"username" validator:"required"`
		Password string `json:"password" validator:"required,password"`
	}

	type LoginResponse struct {
		Token string `json:"token"`
	}

	req := Validated[LoginRequest](dh.Validator, r)
	ctx := r.Context()
	token, err := dh.DeviceService.Login(ctx, req.Username, req.Password)
	if err != nil {
		panic(err)
	}
	respond := LoginResponse{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respond)
}

func (dh *DeviceHandler) GetDevice(w http.ResponseWriter, r *http.Request) {
	type GetDeviceResponse struct {
		Color      string     `json:"color"`
		Volume     uint       `json:"volume"`
		Brightness uint       `json:"brightness"`
		Birthdate  *time.Time `json:"birthdate"`
	}
	ctx := r.Context()
	deviceID := ctx.Value(dh.Constants.Context.DeviceID).(uint)

	device, err := dh.DeviceService.GetDeviceById(ctx, deviceID)
	if err != nil {
		panic(err)
	}

	respond := GetDeviceResponse{
		Color:      device.Color,
		Volume:     device.Volume,
		Brightness: device.Brightness,
		Birthdate:  device.Birthdate,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respond)
}

func (dh *DeviceHandler) GetColor(w http.ResponseWriter, r *http.Request) {
	type GetColorResponse struct {
		Color string `json:"color"`
	}
	ctx := r.Context()
	deviceID := ctx.Value(dh.Constants.Context.DeviceID).(uint)

	color, err := dh.DeviceService.GetDeviceColor(ctx, deviceID)
	if err != nil {
		panic(err)
	}

	respond := GetColorResponse{
		Color: color,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respond)
}

func (dh *DeviceHandler) GetBrightness(w http.ResponseWriter, r *http.Request) {
	type GetBrightnessResponse struct {
		Brightness uint `json:"brightness"`
	}
	ctx := r.Context()
	deviceID := ctx.Value(dh.Constants.Context.DeviceID).(uint)

	device, err := dh.DeviceService.GetDeviceById(ctx, deviceID)
	if err != nil {
		panic(err)
	}

	respond := GetBrightnessResponse{
		Brightness: device.Brightness,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respond)
}

func (dh *DeviceHandler) UpdateColor(w http.ResponseWriter, r *http.Request) {
	type UpdateColorRequest struct {
		Color string `json:"color" validator:"required"`
	}

	ctx := r.Context()
	deviceID := ctx.Value(dh.Constants.Context.DeviceID).(uint)

	req := Validated[UpdateColorRequest](dh.Validator, r)

	err := dh.DeviceService.UpdateDeviceColor(ctx, deviceID, req.Color)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (dh *DeviceHandler) UpdateVolume(w http.ResponseWriter, r *http.Request) {
	type UpdateVolumeRequest struct {
		Volume uint `json:"volume" validator:"required,max=30,min=1"`
	}

	ctx := r.Context()
	deviceID := ctx.Value(dh.Constants.Context.DeviceID).(uint)

	req := Validated[UpdateVolumeRequest](dh.Validator, r)

	err := dh.DeviceService.UpdateDeviceVolume(ctx, deviceID, req.Volume)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (dh *DeviceHandler) Ring(w http.ResponseWriter, r *http.Request) {
	deviceID := r.Context().Value(dh.Constants.Context.DeviceID).(uint)
	err := dh.PublisherService.Ring(deviceID)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (dh *DeviceHandler) Snooze(w http.ResponseWriter, r *http.Request) {
	deviceID := r.Context().Value(dh.Constants.Context.DeviceID).(uint)
	err := dh.PublisherService.Snooze(deviceID)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (dh *DeviceHandler) Silent(w http.ResponseWriter, r *http.Request) {
	deviceID := r.Context().Value(dh.Constants.Context.DeviceID).(uint)
	err := dh.PublisherService.Silent(deviceID)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (dh *DeviceHandler) UpdateBrightness(w http.ResponseWriter, r *http.Request) {
	type UpdateBrightnessRequest struct {
		Brightness uint `json:"brightness" validator:"required,max=255,min=0"`
	}

	ctx := r.Context()
	deviceID := ctx.Value(dh.Constants.Context.DeviceID).(uint)

	req := Validated[UpdateBrightnessRequest](dh.Validator, r)

	err := dh.DeviceService.UpdateDeviceBrightness(ctx, deviceID, req.Brightness)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (dh *DeviceHandler) UpdateBirthdate(w http.ResponseWriter, r *http.Request) {
	type UpdateBirthdateRequest struct {
		Birthdate time.Time `json:"brightness" validator:"required"`
	}

	ctx := r.Context()
	deviceID := ctx.Value(dh.Constants.Context.DeviceID).(uint)

	req := Validated[UpdateBirthdateRequest](dh.Validator, r)

	err := dh.DeviceService.UpdateDeviceBirthdate(ctx, deviceID, req.Birthdate)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusNoContent)
}
