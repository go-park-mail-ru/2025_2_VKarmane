package operation

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/account"
	opRepo "github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/operation"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/operation"
)

type UseCase struct {
	opSvc OperationService
}

func NewUseCase(store *repository.Store) *UseCase {
	accountRepo := account.NewRepository(store.Accounts, store.UserAccounts)
	opRepo := opRepo.NewRepository(store.Operations)
	opService := operation.NewService(accountRepo, opRepo)

	return &UseCase{
		opSvc: opService,
	}
}

func (uc *UseCase) GetAccountOperations(ctx context.Context, accID int) ([]models.Operation, error) {
	opsData, err := uc.opSvc.GetAccountOperations(ctx, accID)
	if err != nil {
		if log := logger.FromContext(ctx); log != nil {
			log.Error("Failed to get ops for acc", "error", err, "user_id", accID)
		}

		return nil, fmt.Errorf("operation.GetUserOperations: %w", err)
	}

	return opsData, nil
}

func (uc *UseCase) GetOperationByID(ctx context.Context, accID int, opID int) (models.Operation, error) {
	opsData, err := uc.opSvc.GetAccountOperations(ctx, accID)
	if err != nil {
		if log := logger.FromContext(ctx); log != nil {
			log.Error("Failed to get op for acc", "error", err, "user_id", accID)
		}

		return models.Operation{}, fmt.Errorf("operation.GetOperationByID: %w", err)
	}

	for _, op := range opsData {
		if op.ID == opID {
			return op, nil
		}
	}

	if log := logger.FromContext(ctx); log != nil {
		log.Warn("Op not found", "acc_id", accID, "operation_id", opID)
	}

	return models.Operation{}, fmt.Errorf("operation.GetOperationByID: %s", models.ErrCodeTransactionNotFound)
}

func (uc *UseCase) CreateOperation(ctx context.Context, req models.CreateOperationRequest, accID int) (models.Operation, error) {
	op, err := uc.opSvc.CreateOperation(ctx, req, accID)
	if err != nil {
		if log := logger.FromContext(ctx); log != nil {
			log.Error("Failed to create op for acc", "error", err, "acc_id", accID)
		}

		return models.Operation{}, fmt.Errorf("operation.CreateOperation: %w", err)
	}

	return op, nil
}

func (uc *UseCase) UpdateOperation(ctx context.Context, req models.UpdateOperationRequest, accID int, opID int) (models.Operation, error) {
	op, err := uc.opSvc.UpdateOperation(ctx, req, accID, opID)
	if err != nil {
		if log := logger.FromContext(ctx); log != nil {
			log.Error("Failed to update op for acc", "error", err, "user_id", accID)
		}

		return models.Operation{}, fmt.Errorf("operation.UpdateOperation: %w", err)
	}

	return op, nil
}

func (uc *UseCase) DeleteOperation(ctx context.Context, accID int, opID int) (models.Operation, error) {
	op, err := uc.opSvc.DeleteOperation(ctx, accID, opID)
	if err != nil {
		if log := logger.FromContext(ctx); log != nil {
			log.Error("Failed to delete op for acc", "error", err, "acc_id", accID)
		}

		return models.Operation{}, fmt.Errorf("operation.DeleteOperation: %w", err)
	}

	return op, nil
}
