package budget

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
)

type PostgresBudgetRepositoryAdapter struct {
	store repository.Repository
}

func NewPostgresBudgetRepositoryAdapter(store repository.Repository) *PostgresBudgetRepositoryAdapter {
	return &PostgresBudgetRepositoryAdapter{
		store: store,
	}
}

func (a *PostgresBudgetRepositoryAdapter) GetBudgetsByUser(ctx context.Context, userID int) ([]models.Budget, error) {
	return a.store.GetBudgetsByUser(ctx, userID)
}

type PostgresAccountRepositoryAdapter struct {
	store repository.Repository
}

func NewPostgresAccountRepositoryAdapter(store repository.Repository) *PostgresAccountRepositoryAdapter {
	return &PostgresAccountRepositoryAdapter{
		store: store,
	}
}

func (a *PostgresAccountRepositoryAdapter) GetAccountsByUser(ctx context.Context, userID int) ([]models.Account, error) {
	return a.store.GetAccountsByUser(ctx, userID)
}

type PostgresOperationRepositoryAdapter struct {
	store repository.Repository
}

func NewPostgresOperationRepositoryAdapter(store repository.Repository) *PostgresOperationRepositoryAdapter {
	return &PostgresOperationRepositoryAdapter{
		store: store,
	}
}

func (a *PostgresOperationRepositoryAdapter) GetOperationsByAccount(ctx context.Context, accountID int) ([]models.Operation, error) {
	return a.store.GetOperationsByAccount(ctx, accountID)
}
