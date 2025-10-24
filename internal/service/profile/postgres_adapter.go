package profile

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
)

type PostgresProfileRepositoryAdapter struct {
	store repository.Repository
}

func NewPostgresProfileRepositoryAdapter(store repository.Repository) *PostgresProfileRepositoryAdapter {
	return &PostgresProfileRepositoryAdapter{store: store}
}

func (a *PostgresProfileRepositoryAdapter) GetUserByID(ctx context.Context, id int) (models.User, error) {
	return a.store.GetUserByID(ctx, id)
}

func (a *PostgresProfileRepositoryAdapter) UpdateUser(ctx context.Context, user models.User) error {
	return a.store.UpdateUser(ctx, user)
}
