package account

import "errors"

var (
	ErrAccountNotFound     = errors.New("account not found")
	ErrUniqueViolation     = errors.New("unique constraint violation")
	ErrForeignKeyViolation = errors.New("foreign key constraint violation")
	ErrNotNullViolation    = errors.New("not null constraint violation")
	ErrCheckViolation      = errors.New("check constraint violation")
)
