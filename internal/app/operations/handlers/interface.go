package operation

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type OperationUseCase interface {
	GetAccountOperations(ctx context.Context, accID int) ([]models.OperationInList, error)
	GetOperationByID(ctx context.Context, accID int, opID int) (models.Operation, error)
	UpdateOperation(ctx context.Context, req models.UpdateOperationRequest, accID int, opID int) (models.Operation, error)
	DeleteOperation(ctx context.Context, accID int, opID int) (models.Operation, error)
	CreateOperation(ctx context.Context, req models.CreateOperationRequest, accID int) (models.Operation, error)
}
