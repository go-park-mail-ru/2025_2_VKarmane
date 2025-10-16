package operation

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type OperationService interface {
	GetAccountOperations(ctx context.Context, accID int) ([]models.Operation, error)
	CreateOperation(ctx context.Context, req models.CreateOperationRequest, accID int) (models.Operation, error)
	UpdateOperation(ctx context.Context, req models.UpdateOperationRequest, accID, opID int) (models.Operation, error)
	DeleteOperation(ctx context.Context, accID, opID int) (models.Operation, error)
}

type OperationRepository interface {
	GetOperationsByAccount(ctx context.Context, accountID int) ([]models.Operation, error)
	GetOperationByID(ctx context.Context, accID int, opID int) (models.Operation, error)
	CreateOperation(ctx context.Context, op models.Operation) (models.Operation, error)
	UpdateOperation(ctx context.Context, req models.UpdateOperationRequest, accID int, opID int) (models.Operation, error)
	DeleteOperation(ctx context.Context, accID int, opID int) (models.Operation, error)
}

type AccountRepository interface {
	GetAccountsByUser(ctx context.Context, userID int) ([]models.Account, error)
}
