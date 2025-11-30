package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/NiflheimDevs/dyslexics-clock/bootstrap"
	"github.com/NiflheimDevs/dyslexics-clock/internal/application/dto"
	"github.com/NiflheimDevs/dyslexics-clock/internal/application/service"
	"github.com/NiflheimDevs/dyslexics-clock/internal/domain/model"
	"github.com/NiflheimDevs/dyslexics-clock/internal/domain/pkg"
	"github.com/go-chi/chi/v5"
)

type AlarmHandler struct {
	Constants    *bootstrap.Constants
	AlarmService service.AlarmService
	Validator    pkg.Validator
}

func NewAlarmHandler(
	constants *bootstrap.Constants,
	alarmService service.AlarmService,
	validator pkg.Validator,
) *AlarmHandler {
	return &AlarmHandler{
		Constants:    constants,
		AlarmService: alarmService,
		Validator:    validator,
	}
}

func (ah *AlarmHandler) GetAlarms(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	deviceID := ctx.Value(ah.Constants.Context.DeviceID).(uint)

	alarms, err := ah.AlarmService.GetAlarms(ctx, deviceID)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(alarms)
}

func (ah *AlarmHandler) UpdateAlarm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	deviceID := ctx.Value(ah.Constants.Context.DeviceID).(uint)

	alarmIDString := chi.URLParam(r, "id")
	alarmID, err := strconv.ParseUint(alarmIDString, 10, 32)
	if err != nil {
		panic(err)
	}

	req := Validated[dto.UpdateAlarm](ah.Validator, r)

	err = ah.AlarmService.UpdateAlarm(ctx, uint(alarmID), deviceID, &req)

	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ah *AlarmHandler) DeleteAlarm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	deviceID := ctx.Value(ah.Constants.Context.DeviceID).(uint)

	alarmIDString := chi.URLParam(r, "id")
	alarmID, err := strconv.ParseUint(alarmIDString, 10, 32)
	if err != nil {
		panic(err)
	}

	err = ah.AlarmService.DeleteAlarmById(ctx, uint(alarmID), deviceID)

	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ah *AlarmHandler) CreateAlarm(w http.ResponseWriter, r *http.Request) {
	type CreateAlarmRequest struct {
		Time          time.Time      `json:"time"`
		IsRepeat      bool           `json:"is_repeat"`
		RepeatingDays []time.Weekday `json:"days"`
	}
	ctx := r.Context()
	deviceID := ctx.Value(ah.Constants.Context.DeviceID).(uint)

	req := Validated[CreateAlarmRequest](ah.Validator, r)

	alarmModel := model.Alarm{
		DeviceId:      deviceID,
		Time:          req.Time,
		IsRepeat:      req.IsRepeat,
		RepeatingDays: req.RepeatingDays,
	}

	err := ah.AlarmService.InsertAlarm(ctx, &alarmModel)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusNoContent)
}
