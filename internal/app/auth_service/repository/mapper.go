package user

import (
	"database/sql"
	"errors"

	serviceerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/errors"
	"github.com/lib/pq"
)

func MapPgError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return serviceerrors.ErrUserNotFound
	}

	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return err
	}

	switch pqErr.Code {
	case UniqueViolation:
		return mapUniqueViolation(pqErr)
	case CheckViolation:
		return serviceerrors.ErrInvalidCredentials
	default:
		return err
	}
}

func mapUniqueViolation(pqErr *pq.Error) error {
	switch pqErr.Constraint {
	case "user_user_login_key":
		return serviceerrors.ErrLoginExists
	case "user_email_key":
		return serviceerrors.ErrEmailExists
	}

	return pqErr
}
