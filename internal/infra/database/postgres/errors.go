package repositoryimpl

import (
	"context"
	"errors"

	derror "github.com/NiflheimDevs/dyslexics-clock/internal/domain/error"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func NormalizeDBError(err error, message string) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return derror.New(derror.ErrTypeNotFound, message, err)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return derror.New(derror.ErrTypeConflict, message, err)

		case "23503": // foreign key
			return derror.New(derror.ErrTypeForeignKey, message, err)

		default:
			return derror.New(derror.ErrTypeDB, message, err)
		}
	}

	if errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, context.Canceled) {
		return derror.New(derror.ErrTypeTimeout, message, err)
	}

	return derror.New(derror.ErrTypeInternal, message, err)
}
