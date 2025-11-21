package category

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/category/service"
)

type UseCase struct {
	categorySvc category.CategoryService
}

func NewUseCase(categorySvc category.CategoryService) *UseCase {
	return &UseCase{
		categorySvc: categorySvc,
	}
}

func (uc *UseCase) CreateCategory(ctx context.Context, req models.CreateCategoryRequest, userID int) (models.Category, error) {
	log := logger.FromContext(ctx)
	category, err := uc.categorySvc.CreateCategory(ctx, req, userID)
	if err != nil {
		log.Error("Failed to create category", "error", err, "user_id", userID)
		return models.Category{}, fmt.Errorf("category.CreateCategory: %w", err)
	}

	return category, nil
}

func (uc *UseCase) GetCategoriesByUser(ctx context.Context, userID int) ([]models.CategoryWithStats, error) {
	log := logger.FromContext(ctx)
	categories, err := uc.categorySvc.GetCategoriesByUser(ctx, userID)
	if err != nil {
		log.Error("Failed to get categories for user", "error", err, "user_id", userID)
		return nil, fmt.Errorf("category.GetCategoriesByUser: %w", err)
	}

	return categories, nil
}

func (uc *UseCase) GetCategoryByID(ctx context.Context, userID, categoryID int) (models.CategoryWithStats, error) {
	log := logger.FromContext(ctx)
	category, err := uc.categorySvc.GetCategoryByID(ctx, userID, categoryID)
	if err != nil {
		log.Error("Failed to get category by ID", "error", err, "user_id", userID, "category_id", categoryID)
		return models.CategoryWithStats{}, fmt.Errorf("category.GetCategoryByID: %w", err)
	}

	return category, nil
}

func (uc *UseCase) UpdateCategory(ctx context.Context, req models.UpdateCategoryRequest, userID, categoryID int) (models.Category, error) {
	log := logger.FromContext(ctx)
	category, err := uc.categorySvc.UpdateCategory(ctx, req, userID, categoryID)
	if err != nil {
		log.Error("Failed to update category", "error", err, "user_id", userID, "category_id", categoryID)
		return models.Category{}, fmt.Errorf("category.UpdateCategory: %w", err)
	}

	return category, nil
}

func (uc *UseCase) DeleteCategory(ctx context.Context, userID, categoryID int) error {
	log := logger.FromContext(ctx)
	err := uc.categorySvc.DeleteCategory(ctx, userID, categoryID)
	if err != nil {
		log.Error("Failed to delete category", "error", err, "user_id", userID, "category_id", categoryID)
		return fmt.Errorf("category.DeleteCategory: %w", err)
	}

	return nil
}
