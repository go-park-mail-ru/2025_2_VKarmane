package errors

import "errors"

var (
	ErrAccountNotFound   = errors.New("account not found")
	ErrOperationNotFound = errors.New("operation not found")
	ErrCategoryNotFound  = errors.New("category not found")
	ErrCategoryExists    = errors.New("category already exists")
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidData       = errors.New("invalid data")
)
