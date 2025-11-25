package errors

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"google.golang.org/grpc/codes"
)

var ErrorMap = map[error]struct {
	Code codes.Code
	Msg string
} {
	ErrBudgetExists: {Code: codes.AlreadyExists, Msg: string(models.ErrCodeBudgetExists)},
	ErrBudgetNotFound: {Code: codes.NotFound, Msg: string(models.ErrCodeBudgetNotFound)},
	ErrForbidden: {Code: codes.PermissionDenied, Msg: string(models.ErrCodeForbidden)},
	ErrInavlidData: {Code: codes.InvalidArgument, Msg: string(models.ErrCodeInvalidData)},
}