package errors

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"google.golang.org/grpc/codes"
)

var ErrorMap = map[error]struct {
	Code codes.Code
	Msg  string
}{
	ErrCategoryExists:    {Code: codes.AlreadyExists, Msg: string(models.ErrCodeCategoryExists)},
	ErrCategoryNotFound:  {Code: codes.NotFound, Msg: string(models.ErrCodeCategoryNotFound)},
	ErrOperationNotFound: {Code: codes.NotFound, Msg: string(models.ErrCodeTransactionNotFound)},
	ErrAccountNotFound:   {Code: codes.NotFound, Msg: string(models.ErrCodeAccountNotFound)},
	ErrForbidden:         {Code: codes.PermissionDenied, Msg: string(models.ErrCodeForbidden)},
	ErrInvalidData:       {Code: codes.InvalidArgument, Msg: string(models.ErrCodeInvalidData)},
}
