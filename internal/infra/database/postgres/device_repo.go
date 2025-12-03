package repositoryimpl

import (
	"context"

	"github.com/NiflheimDevs/dyslexics-clock/internal/application/dto"
	"github.com/NiflheimDevs/dyslexics-clock/internal/domain/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DeviceRepo struct {
	DB *pgxpool.Pool
}

func NewDeviceRepo(db *pgxpool.Pool) *DeviceRepo {
	return &DeviceRepo{
		DB: db,
	}
}

func (d *DeviceRepo) GetDeviceByUsername(ctx context.Context, username string) (*dto.LoginDto, error) {
	query := `SELECT id, username, password FROM devices WHERE username = $1`
	var result dto.LoginDto
	err := d.DB.QueryRow(ctx, query, username).Scan(&result.Id, &result.Username, &result.Password)
	return &result, NormalizeDBError(err, "failed to get device")
}

func (d *DeviceRepo) GetDeviceById(ctx context.Context, Id uint) (*model.Device, error) {
	query := `SELECT id, username, password, color FROM devices WHERE Id = $1`
	var result model.Device
	err := d.DB.QueryRow(ctx, query, Id).Scan(&result.ID, &result.Username, &result.Password, &result.Color)
	if err != nil {
		return nil, NormalizeDBError(err, "failed to get device")
	}
	return &result, NormalizeDBError(err, "failed to get device")
}

func (d *DeviceRepo) UpdateColor(ctx context.Context, Id uint, color string) error {
	query := `UPDATE devices SET color = $1 WHERE Id = $2`
	_, err := d.DB.Exec(ctx, query, color, Id)
	return NormalizeDBError(err, "failed to update color")
}
