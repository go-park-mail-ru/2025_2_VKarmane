package category

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type Service struct {
	repo interface {
		CreateCategory(ctx context.Context, category models.Category) (models.Category, error)
		GetCategoriesByUser(ctx context.Context, userID int) ([]models.Category, error)
		GetCategoryByID(ctx context.Context, userID, categoryID int) (models.Category, error)
		UpdateCategory(ctx context.Context, category models.Category) error
		DeleteCategory(ctx context.Context, userID, categoryID int) error
		GetCategoryStats(ctx context.Context, userID, categoryID int) (int, error)
	}
}

func NewService(repo interface {
	CreateCategory(ctx context.Context, category models.Category) (models.Category, error)
	GetCategoriesByUser(ctx context.Context, userID int) ([]models.Category, error)
	GetCategoryByID(ctx context.Context, userID, categoryID int) (models.Category, error)
	UpdateCategory(ctx context.Context, category models.Category) error
	DeleteCategory(ctx context.Context, userID, categoryID int) error
	GetCategoryStats(ctx context.Context, userID, categoryID int) (int, error)
}) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateCategory(ctx context.Context, req models.CreateCategoryRequest, userID int) (models.Category, error) {
	logoHashedID := req.LogoHashedID
	if logoHashedID == "" {
		logoHashedID = "c1dfd96eea8cc2b62785275bca38ac261256e278"
	}

	category := models.Category{
		UserID:       userID,
		Name:         req.Name,
		Description:  req.Description,
		LogoHashedID: logoHashedID,
	}

	createdCategory, err := s.repo.CreateCategory(ctx, category)
	if err != nil {
		return models.Category{}, fmt.Errorf("failed to create category: %w", err)
	}

	return createdCategory, nil
}

func (s *Service) GetCategoriesByUser(ctx context.Context, userID int) ([]models.CategoryWithStats, error) {
	categories, err := s.repo.GetCategoriesByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	var categoriesWithStats []models.CategoryWithStats
	for _, category := range categories {
		stats, err := s.repo.GetCategoryStats(ctx, userID, category.ID)
		if err != nil {
			stats = 0
		}

		categoriesWithStats = append(categoriesWithStats, models.CategoryWithStats{
			Category:        category,
			OperationsCount: stats,
		})
	}

	return categoriesWithStats, nil
}

func (s *Service) GetCategoryByID(ctx context.Context, userID, categoryID int) (models.CategoryWithStats, error) {
	category, err := s.repo.GetCategoryByID(ctx, userID, categoryID)
	if err != nil {
		return models.CategoryWithStats{}, fmt.Errorf("failed to get category: %w", err)
	}

	stats, err := s.repo.GetCategoryStats(ctx, userID, categoryID)
	if err != nil {
		stats = 0
	}

	return models.CategoryWithStats{
		Category:        category,
		OperationsCount: stats,
	}, nil
}

func (s *Service) UpdateCategory(ctx context.Context, req models.UpdateCategoryRequest, userID, categoryID int) (models.Category, error) {
	existingCategory, err := s.repo.GetCategoryByID(ctx, userID, categoryID)
	if err != nil {
		return models.Category{}, fmt.Errorf("failed to get category: %w", err)
	}

	if req.Name != nil {
		existingCategory.Name = *req.Name
	}
	if req.Description != nil {
		existingCategory.Description = *req.Description
	}
	if req.LogoHashedID != nil {
		existingCategory.LogoHashedID = *req.LogoHashedID
	}

	err = s.repo.UpdateCategory(ctx, existingCategory)
	if err != nil {
		return models.Category{}, fmt.Errorf("failed to update category: %w", err)
	}

	return existingCategory, nil
}

func (s *Service) DeleteCategory(ctx context.Context, userID, categoryID int) error {
	err := s.repo.DeleteCategory(ctx, userID, categoryID)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}
