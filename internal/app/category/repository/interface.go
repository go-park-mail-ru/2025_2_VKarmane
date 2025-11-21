package category

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/dto"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category dto.CategoryDB) (int, error)
	GetCategoriesByUser(ctx context.Context, userID int) ([]dto.CategoryDB, error)
	GetCategoryByID(ctx context.Context, userID, categoryID int) (dto.CategoryDB, error)
	UpdateCategory(ctx context.Context, category dto.CategoryDB) error
	DeleteCategory(ctx context.Context, userID, categoryID int) error
	GetCategoryStats(ctx context.Context, userID, categoryID int) (int, error)
}
