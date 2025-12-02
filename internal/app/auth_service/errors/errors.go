package errors

import "errors"

var (
	ErrUserNotFound       = errors.New("USER_NOT_FOUND")
	ErrUserExists         = errors.New("USER_EXISTS")
	ErrInvalidCredentials = errors.New("INVALID_CREDENTIALS")
	ErrLoginExists        = errors.New("LOGIN_EXISTS")
	ErrEmailExists        = errors.New("EMAIL_EXISTS")
	ErrForbidden          = errors.New("FORBIDDEN")
)
