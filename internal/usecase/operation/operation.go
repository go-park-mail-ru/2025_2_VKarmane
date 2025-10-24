package operation

import (
	"context"
	"fmt"

	pkgErrors "github.com/pkg/errors"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/operation"
)

type UseCase struct {
	opSvc operation.OperationService
}

func NewUseCase(store repository.Repository) *UseCase {
	postgresStore := store.(*repository.PostgresStore)

	accountRepoAdapter := operation.NewPostgresAccountRepositoryAdapter(postgresStore)
	operationRepoAdapter := operation.NewPostgresOperationRepositoryAdapter(postgresStore)
	opService := operation.NewService(accountRepoAdapter, operationRepoAdapter)

	return &UseCase{
		opSvc: opService,
	}
}

func (uc *UseCase) GetAccountOperations(ctx context.Context, accID int) ([]models.Operation, error) {
	log := logger.FromContext(ctx)

	opsData, err := uc.opSvc.GetAccountOperations(ctx, accID)
	if err != nil {
		if log != nil {
			log.Error("Failed to get ops for acc", "error", err, "account_id", accID)
		}

		return nil, pkgErrors.Wrap(err, "operation.GetAccountOperations")
	}

	return opsData, nil
}

func (uc *UseCase) GetOperationByID(ctx context.Context, accID int, opID int) (models.Operation, error) {
	log := logger.FromContext(ctx)

	opData, err := uc.opSvc.GetOperationByID(ctx, accID, opID)
	if err != nil {
		if log != nil {
			log.Error("Failed to get operation by ID", "error", err, "account_id", accID, "operation_id", opID)
		}

		return models.Operation{}, pkgErrors.Wrap(err, "operation.GetOperationByID")
	}

	return opData, nil
}

func (uc *UseCase) CreateOperation(ctx context.Context, req models.CreateOperationRequest, accID int) (models.Operation, error) {
	log := logger.FromContext(ctx)

	// Валидация запроса
	if req.Sum <= 0 {
		return models.Operation{}, fmt.Errorf("invalid sum: %s", models.ErrCodeInvalidAmount)
	}

	if req.Name == "" {
		return models.Operation{}, fmt.Errorf("name is required: %s", models.ErrCodeMissingFields)
	}

	opData, err := uc.opSvc.CreateOperation(ctx, req, accID)
	if err != nil {
		if log != nil {
			log.Error("Failed to create operation", "error", err, "account_id", accID)
		}

		return models.Operation{}, pkgErrors.Wrap(err, "operation.CreateOperation")
	}

	return opData, nil
}

func (uc *UseCase) UpdateOperation(ctx context.Context, req models.UpdateOperationRequest, accID int, opID int) (models.Operation, error) {
	log := logger.FromContext(ctx)

	// Валидация запроса
	if req.Sum != nil && *req.Sum <= 0 {
		return models.Operation{}, fmt.Errorf("invalid sum: %s", models.ErrCodeInvalidAmount)
	}

	opData, err := uc.opSvc.UpdateOperation(ctx, req, accID, opID)
	if err != nil {
		if log != nil {
			log.Error("Failed to update operation", "error", err, "account_id", accID, "operation_id", opID)
		}

		return models.Operation{}, pkgErrors.Wrap(err, "operation.UpdateOperation")
	}

	return opData, nil
}

func (uc *UseCase) DeleteOperation(ctx context.Context, accID int, opID int) (models.Operation, error) {
	log := logger.FromContext(ctx)

	opData, err := uc.opSvc.DeleteOperation(ctx, accID, opID)
	if err != nil {
		if log != nil {
			log.Error("Failed to delete operation", "error", err, "account_id", accID, "operation_id", opID)
		}

		return models.Operation{}, pkgErrors.Wrap(err, "operation.DeleteOperation")
	}

	return opData, nil
}
