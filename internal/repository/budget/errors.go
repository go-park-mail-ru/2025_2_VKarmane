package budget

import "errors"



var (
	ErrBudgetNotFound         = errors.New("budget not found")
	ErrActiveBudgetExists     = errors.New("active budget already exists for this category")
	ErrForeignKeyViolation    = errors.New("foreign key violation")
	ErrNotNullViolation       = errors.New("not null violation")
	ErrCheckViolation         = errors.New("check constraint violation")
	ErrUniqueViolation        = errors.New("unique constraint violation")
)
