package operation

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
)

type PostgresAccountRepositoryAdapter struct {
	store *repository.PostgresStore
}

func NewPostgresAccountRepositoryAdapter(store *repository.PostgresStore) *PostgresAccountRepositoryAdapter {
	return &PostgresAccountRepositoryAdapter{
		store: store,
	}
}

func (a *PostgresAccountRepositoryAdapter) GetAccountsByUser(ctx context.Context, userID int) ([]models.Account, error) {
	return a.store.AccountRepo.GetAccountsByUser(ctx, userID)
}

type PostgresOperationRepositoryAdapter struct {
	store *repository.PostgresStore
}

func NewPostgresOperationRepositoryAdapter(store *repository.PostgresStore) *PostgresOperationRepositoryAdapter {
	return &PostgresOperationRepositoryAdapter{
		store: store,
	}
}

func (a *PostgresOperationRepositoryAdapter) GetOperationsByAccount(ctx context.Context, accountID int) ([]models.Operation, error) {
	return a.store.OperationRepo.GetOperationsByAccount(ctx, accountID)
}

func (a *PostgresOperationRepositoryAdapter) GetOperationByID(ctx context.Context, accID int, opID int) (models.Operation, error) {
	return a.store.OperationRepo.GetOperationByID(ctx, accID, opID)
}

func (a *PostgresOperationRepositoryAdapter) CreateOperation(ctx context.Context, op models.Operation) (models.Operation, error) {
	return a.store.OperationRepo.CreateOperation(ctx, op)
}

func (a *PostgresOperationRepositoryAdapter) UpdateOperation(ctx context.Context, req models.UpdateOperationRequest, accID int, opID int) (models.Operation, error) {
	return a.store.OperationRepo.UpdateOperation(ctx, req, accID, opID)
}

func (a *PostgresOperationRepositoryAdapter) DeleteOperation(ctx context.Context, accID int, opID int) (models.Operation, error) {
	return a.store.OperationRepo.DeleteOperation(ctx, accID, opID)
}
