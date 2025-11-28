package repository

import (
	"database/sql"
	"errors"

	serviceerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/errors"
	"github.com/lib/pq"
)

func MapPgCategoryError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return serviceerrors.ErrCategoryNotFound
	}

	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return err
	}

	switch pqErr.Code {
	case UniqueViolation:
		return serviceerrors.ErrCategoryExists
	case NotNullViolation:
		return serviceerrors.ErrInvalidData
	case CheckViolation:
		return serviceerrors.ErrInvalidData
	default:
		return err
	}
}

func MapPgOperationError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return serviceerrors.ErrOperationNotFound
	}

	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return err
	}

	switch pqErr.Code {
	case NotNullViolation:
		return serviceerrors.ErrInvalidData
	case CheckViolation:
		return serviceerrors.ErrInvalidData
	default:
		return err
	}
}

func MapPgAccountError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return serviceerrors.ErrAccountNotFound
	}

	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return err
	}

	switch pqErr.Code {
	case NotNullViolation:
		return serviceerrors.ErrInvalidData
	case CheckViolation:
		return serviceerrors.ErrNegativeBalance
	default:
		return err
	}
}
