package category

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, req models.CreateCategoryRequest, userID int) (models.Category, error)
	GetCategoriesByUser(ctx context.Context, userID int) ([]models.CategoryWithStats, error)
	GetCategoryByID(ctx context.Context, userID, categoryID int) (models.CategoryWithStats, error)
	UpdateCategory(ctx context.Context, req models.UpdateCategoryRequest, userID, categoryID int) (models.Category, error)
	DeleteCategory(ctx context.Context, userID, categoryID int) error
}

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category models.Category) (models.Category, error)
	GetCategoriesByUser(ctx context.Context, userID int) ([]models.Category, error)
	GetCategoryByID(ctx context.Context, userID, categoryID int) (models.Category, error)
	UpdateCategory(ctx context.Context, category models.Category) error
	DeleteCategory(ctx context.Context, userID, categoryID int) error
	GetCategoryStats(ctx context.Context, userID, categoryID int) (int, error)
}
