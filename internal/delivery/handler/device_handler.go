package handler

import (
	"encoding/json"
	"net/http"

	"github.com/NiflheimDevs/dyslexics-clock/bootstrap"
	"github.com/NiflheimDevs/dyslexics-clock/internal/application/service"
	"github.com/NiflheimDevs/dyslexics-clock/internal/domain/pkg"
)

type DeviceHandler struct {
	Constants     *bootstrap.Constants
	DeviceService service.DeviceService
	Validator     pkg.Validator
}

func NewDeviceHandler(
	constants *bootstrap.Constants,
	deviceService service.DeviceService,
	validator pkg.Validator,
) *DeviceHandler {
	return &DeviceHandler{
		Constants:     constants,
		DeviceService: deviceService,
		Validator:     validator,
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

