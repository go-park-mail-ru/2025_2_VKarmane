package category

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/dto"
)

type PostgresCategoryRepositoryAdapter struct {
	store repository.Repository
}

func NewPostgresCategoryRepositoryAdapter(store repository.Repository) *PostgresCategoryRepositoryAdapter {
	return &PostgresCategoryRepositoryAdapter{store: store}
}

func (a *PostgresCategoryRepositoryAdapter) CreateCategory(ctx context.Context, category models.Category) (models.Category, error) {
	postgresStore := a.store.(*repository.PostgresStore)

	var description *string
	if category.Description != "" {
		description = &category.Description
	}

	categoryDB := dto.CategoryDB{
		UserID:       category.UserID,
		Name:         category.Name,
		Description:  description,
		LogoHashedID: category.LogoHashedID,
	}

	id, err := postgresStore.CategoryRepo.CreateCategory(ctx, categoryDB)
	if err != nil {
		return models.Category{}, err
	}

	category.ID = id
	return category, nil
}

func (a *PostgresCategoryRepositoryAdapter) GetCategoriesByUser(ctx context.Context, userID int) ([]models.Category, error) {
	return a.store.GetCategoriesByUser(ctx, userID)
}

func (a *PostgresCategoryRepositoryAdapter) GetCategoryByID(ctx context.Context, userID, categoryID int) (models.Category, error) {
	return a.store.GetCategoryByID(ctx, userID, categoryID)
}

func (a *PostgresCategoryRepositoryAdapter) UpdateCategory(ctx context.Context, category models.Category) error {
	postgresStore := a.store.(*repository.PostgresStore)

	var description *string
	if category.Description != "" {
		description = &category.Description
	}

	categoryDB := dto.CategoryDB{
		ID:           category.ID,
		UserID:       category.UserID,
		Name:         category.Name,
		Description:  description,
		LogoHashedID: category.LogoHashedID,
		CreatedAt:    category.CreatedAt,
		UpdatedAt:    category.UpdatedAt,
	}

	return postgresStore.CategoryRepo.UpdateCategory(ctx, categoryDB)
}

func (a *PostgresCategoryRepositoryAdapter) DeleteCategory(ctx context.Context, userID, categoryID int) error {
	postgresStore := a.store.(*repository.PostgresStore)
	return postgresStore.CategoryRepo.DeleteCategory(ctx, userID, categoryID)
}

func (a *PostgresCategoryRepositoryAdapter) GetCategoryStats(ctx context.Context, userID, categoryID int) (int, error) {
	postgresStore := a.store.(*repository.PostgresStore)
	return postgresStore.CategoryRepo.GetCategoryStats(ctx, userID, categoryID)
}
