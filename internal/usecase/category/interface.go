package category

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type CategoryUseCase interface {
	CreateCategory(ctx context.Context, req models.CreateCategoryRequest, userID int) (models.Category, error)
	GetCategoriesByUser(ctx context.Context, userID int) ([]models.CategoryWithStats, error)
	GetCategoryByID(ctx context.Context, userID, categoryID int) (models.CategoryWithStats, error)
	UpdateCategory(ctx context.Context, req models.UpdateCategoryRequest, userID, categoryID int) (models.Category, error)
	DeleteCategory(ctx context.Context, userID, categoryID int) error
}
