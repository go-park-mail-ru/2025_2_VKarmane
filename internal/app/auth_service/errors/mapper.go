package errors

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"google.golang.org/grpc/codes"
)

var ErrorMap = map[error]struct {
	Code codes.Code
	Msg string
} {
	ErrLoginExists: {Code: codes.AlreadyExists, Msg: string(models.ErrCodeLoginExists)},
    ErrEmailExists: {Code: codes.AlreadyExists, Msg: string(models.ErrCodeEmailExists)},
	ErrInvalidCredentials: {Code: codes.Unauthenticated, Msg: string(models.ErrCodeInvalidCredentials)},
	ErrUserNotFound: {Code: codes.NotFound, Msg: string(models.ErrCodeUserNotFound)},
	ErrForbidden: {Code: codes.PermissionDenied, Msg: string(models.ErrCodeForbidden)},
	ErrUserExists: {Code: codes.AlreadyExists, Msg: string(models.ErrCodeUserExists)},
}