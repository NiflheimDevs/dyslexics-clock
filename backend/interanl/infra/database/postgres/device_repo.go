package repositoryimpl

import (
	"context"
	"encoding/json"

	"github.com/NiflheimDevs/wclock/interanl/application/dto"
	"github.com/NiflheimDevs/wclock/interanl/domain/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DeviceRepo struct {
	DB *pgxpool.Pool
}

func (d *DeviceRepo) GetDeviceByUsername(ctx context.Context, username string) (*dto.LoginDto, error) {
	query := `SELECT id, username, password FROM devices WHERE username = $1`
	var result dto.LoginDto
	err := d.DB.QueryRow(ctx, query, username).Scan(&result.Id, &result.Username, &result.Password)
	return &result, err
}

func (d *DeviceRepo) GetDeviceById(ctx context.Context, Id uint) (*model.Device, error) {
	query := `SELECT id, username, password, alarms, color FROM devices WHERE Id = $1`
	var alarms []byte
	var result model.Device
	err := d.DB.QueryRow(ctx, query, Id).Scan(&result.ID, &result.Username, &result.Password, &alarms, &result.Color)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(alarms, &result.Alarms)
	return &result, err
}

func (d *DeviceRepo) UpdateColor(ctx context.Context, Id uint, color string) error {
	query := `UPDATE devices SET color = $1 WHERE Id = $2`
	_, err := d.DB.Exec(ctx, query, color, Id)
	return err
}
