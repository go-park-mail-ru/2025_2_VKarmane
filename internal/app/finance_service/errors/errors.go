package errors

import "errors"

var (
	ErrAccountNotFound   = errors.New("account not found")
	ErrOperationNotFound = errors.New("operation not found")
	ErrCategoryNotFound  = errors.New("category not found")
	ErrCategoryExists    = errors.New("category already exists")
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidData       = errors.New("invalid data")
	ErrNegativeBalance   = errors.New("negative balance")
	ErrPrivateAccount    = errors.New("private account")
	ErrSharingExists     = errors.New("sharing exists")
	ErrUserNotFound      = errors.New("user not found")
)
