package errors

import "errors"

var (
	ErrBudgetNotFound = errors.New("BUDGET_NOT_FOUND")
	ErrBudgetExists = errors.New("BUDGET_EXISTS")
	ErrForbidden = errors.New("FORBIDDEN")
	ErrInavlidData = errors.New("INVALID_DATA")
)


