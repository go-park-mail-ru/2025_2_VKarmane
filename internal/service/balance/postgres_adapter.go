package balance

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
)

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
