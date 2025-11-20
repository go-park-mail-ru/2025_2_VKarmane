package budget

import (
	"database/sql"
	"errors"

	serviceerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/errors"
	"github.com/lib/pq"
)

func MapPgError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return serviceerrors.ErrBudgetNotFound
	}

	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return err
	}

	switch pqErr.Code {
	case UniqueViolation:
		return serviceerrors.ErrBudgetExists
	case NotNullViolation:
		return serviceerrors.ErrInavlidData
	case CheckViolation:
		return serviceerrors.ErrInavlidData
	default:
		return err
	}
}
