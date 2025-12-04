package repositoryimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/NiflheimDevs/dyslexics-clock/internal/application/dto"
	derror "github.com/NiflheimDevs/dyslexics-clock/internal/domain/error"
	"github.com/NiflheimDevs/dyslexics-clock/internal/domain/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AlarmRepo struct {
	DB *pgxpool.Pool
}

func NewAlarmRepo(db *pgxpool.Pool) *AlarmRepo {
	return &AlarmRepo{
		DB: db,
	}
}

func (d *AlarmRepo) InsertAlarm(ctx context.Context, alarm *model.Alarm) error {
	query := `INSERT INTO alarms (device_id, time, is_repeat, days) VALUES ($1, $2, $3, $4)`
	alarmsBytes, err := json.Marshal(alarm.RepeatingDays)
	if err != nil {
		return NormalizeDBError(err, "failed to insert alarm")
	}
	_, err = d.DB.Exec(ctx, query, alarm.DeviceId, alarm.Time, alarm.IsRepeat, alarmsBytes)
	return NormalizeDBError(err, "failed to insert alarm")
}

func (d *AlarmRepo) DeleteAlarmById(ctx context.Context, alarmID uint, deviceID uint) (int64, error) {
	query := `DELETE FROM alarms WHERE id = $1 AND device_id = $2`
	res, err := d.DB.Exec(ctx, query, alarmID, deviceID)

	return res.RowsAffected(), NormalizeDBError(err, "failed to delete alarm")
}

func(d *AlarmRepo) buildUpdateQuery(alarmID uint, deviceID uint ,input *dto.UpdateAlarm, alarmDayBytes []byte) (string, []any) {
    var sets []string
    var args []any
    argIndex := 1
    
    if input.Time != nil {
        sets = append(sets, fmt.Sprintf("time = $%d", argIndex))
        args = append(args, *input.Time)
        argIndex++
    }
    
    if input.IsRepeat != nil {
        sets = append(sets, fmt.Sprintf("is_repeat = $%d", argIndex))
        args = append(args, *input.IsRepeat)
        argIndex++
    }
    
    if input.RepeatingDays != nil {
        sets = append(sets, fmt.Sprintf("repeating_days = $%d", argIndex))
        args = append(args, alarmDayBytes)
        argIndex++
    }
    
    if len(sets) == 0 {
        return "", nil // No updates
    }
    
    query := fmt.Sprintf(
        "UPDATE alarms SET %s WHERE id = $%d AND device_id = $%d",
        strings.Join(sets, ", "),
        argIndex,
		argIndex+1,
    )
    args = append(args, alarmID)
    args = append(args, deviceID)
    
    return query, args
}
func (d *AlarmRepo) UpdateAlarm(ctx context.Context, alarmID uint, deviceID uint, updateAlarm *dto.UpdateAlarm) (int64, error) {
	// query := `UPDATE alarms SET time = $2, is_repeat = $3, days = $4 WHERE id = $1 AND device_id = $5`
	alarmsBytes, err := json.Marshal(updateAlarm.RepeatingDays)
	if err != nil {
		return -1, derror.New(derror.ErrTypeBadRequest, "invalid repeating days", err)
	}
	
	query, args := d.buildUpdateQuery(alarmID, deviceID, updateAlarm, alarmsBytes)
	// res, err := d.DB.Exec(ctx, query, alarmID, updateAlarm.Time, updateAlarm.IsRepeat, alarmsBytes, deviceID)
	res, err := d.DB.Exec(ctx, query, args...)
	return res.RowsAffected(), NormalizeDBError(err, "failed to update alarm")
}

func (d *AlarmRepo) GetAlarms(ctx context.Context, DeviceId uint) ([]model.Alarm, error) {
	query := `SELECT id, device_id, time, is_repeat, days FROM alarms where device_id = $1`
	var result []model.Alarm
	rows, err := d.DB.Query(ctx, query, DeviceId)
	if err != nil {
		return nil, NormalizeDBError(err, "failed to get alarms")
	}
	defer rows.Close()
	for rows.Next() {
		var days []byte
		var item model.Alarm
		err = rows.Scan(&item.ID, &item.DeviceId, &item.Time, &item.IsRepeat, &days)
		if err != nil {
			return nil, NormalizeDBError(err, "failed to get alarm")
		}
		err = json.Unmarshal(days, &item.RepeatingDays)
		if err != nil {
			return nil, NormalizeDBError(err, "failed to get alarm repeating days")
		}
		result = append(result, item)
	}
	return result, nil
}

func (d *AlarmRepo) GetAlarmById(ctx context.Context, DeviceId uint) (*model.Alarm, error) {
	query := `SELECT id, device_id, time, is_repeat, days FROM alarms where device_id = $1`
	var days []byte
	var result model.Alarm
	err := d.DB.QueryRow(ctx, query, DeviceId).Scan(&result.ID, &result.DeviceId, &result.Time, &result.IsRepeat, &days)
	if err != nil {
		return nil, NormalizeDBError(err, "failed to get alarm")
	}
	err = json.Unmarshal(days, &result.RepeatingDays)
	if err != nil {
		return nil, NormalizeDBError(err, "failed to get alarm repeating days")
	}
	return &result, nil
}
