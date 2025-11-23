package repositoryimpl

import (
	"context"
	"encoding/json"

	"github.com/NiflheimDevs/dyslexics-clock/internal/application/dto"
	"github.com/NiflheimDevs/dyslexics-clock/internal/domain/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AlarmRepo struct {
	DB *pgxpool.Pool
}

func (d *AlarmRepo) InsertAlarm(ctx context.Context, alarm *model.Alarm) error {
	query := `INSERT INTO alarms (time, is_repeat, days) VALUES ($1, $2, $3)`
	alarmsBytes, err := json.Marshal(alarm.RepeatingDays)
	if err != nil {
		return err
	}
	_, err = d.DB.Exec(ctx, query, alarm.Time, alarm.IsRepeat, alarmsBytes)
	return err
}

func (d *AlarmRepo) DeleteAlarmById(ctx context.Context, id uint) error {
	query := `DELETE alarms WHERE id = $1`
	_, err := d.DB.Exec(ctx, query, id)
	return err
}

func (d *AlarmRepo) UpdateAlarm(ctx context.Context, id uint, updateAlarm *dto.UpdateAlarm) error {
	query := `UPDATE alarms SET time = $2, is_repeat = $3, days = $4 WHERE id = $1`
	alarmsBytes, err := json.Marshal(updateAlarm.RepeatingDays)
	if err != nil {
		return err
	}
	_, err = d.DB.Exec(ctx, query, id, updateAlarm.Time, updateAlarm.IsRepeat, alarmsBytes)
	return err
}

func (d *AlarmRepo) GetAlarms(ctx context.Context, DeviceId uint) ([]model.Alarm, error) {
	query := `SELECT id, device_id, time, is_repeat, days FROM alarms where device_id = $1`
	var result []model.Alarm
	rows, err := d.DB.Query(ctx, query, DeviceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var days []byte
		var item model.Alarm
		err = rows.Scan(&item.ID, &item.DeviceId, &item.Time, &item.IsRepeat, &days)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(days, &item.RepeatingDays)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, err
}

func (d *AlarmRepo) GetAlarmById(ctx context.Context, DeviceId uint) (*model.Alarm, error) {
	query := `SELECT id, device_id, time, is_repeat, days FROM alarms where device_id = $1`
	var days []byte
	var result model.Alarm
	err := d.DB.QueryRow(ctx, query, DeviceId).Scan(&result.ID, &result.DeviceId, &result.Time, &result.IsRepeat, &days)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(days, &result.RepeatingDays)
	if err != nil {
		return nil, err
	}
	return &result, err
}
