package category

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type combinedRepo struct {
	catRepo *mocks.MockCategoryRepository
}

func (c *combinedRepo) CreateCategory(ctx context.Context, category models.Category) (models.Category, error) {
	return c.catRepo.CreateCategory(ctx, category)
}

func (c *combinedRepo) GetCategoriesByUser(ctx context.Context, userID int) ([]models.Category, error) {
	return c.catRepo.GetCategoriesByUser(ctx, userID)
}

func (c *combinedRepo) GetCategoryByID(ctx context.Context, userID, categoryID int) (models.Category, error) {
	return c.catRepo.GetCategoryByID(ctx, userID, categoryID)
}

func (c *combinedRepo) UpdateCategory(ctx context.Context, category models.Category) error {
	return c.catRepo.UpdateCategory(ctx, category)
}

func (c *combinedRepo) DeleteCategory(ctx context.Context, userID, categoryID int) error {
	return c.catRepo.DeleteCategory(ctx, userID, categoryID)
}

func (c *combinedRepo) GetCategoryStats(ctx context.Context, userID, categoryID int) (int, error) {
	return c.catRepo.GetCategoryStats(ctx, userID, categoryID)
}

func TestService_CreateCategory_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCatRepo := mocks.NewMockCategoryRepository(ctrl)
	repo := &combinedRepo{catRepo: mockCatRepo}
	svc := NewService(repo)

	req := models.CreateCategoryRequest{
		Name:        "Food",
		Description: "Food expenses",
	}
	expectedCat := models.Category{ID: 1, Name: "Food", LogoHashedID: "c1dfd96eea8cc2b62785275bca38ac261256e278"}

	mockCatRepo.EXPECT().CreateCategory(gomock.Any(), gomock.Any()).Return(expectedCat, nil)

	result, err := svc.CreateCategory(context.Background(), req, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedCat.ID, result.ID)
	assert.Equal(t, expectedCat.Name, result.Name)
}

func TestService_GetCategoriesByUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCatRepo := mocks.NewMockCategoryRepository(ctrl)
	repo := &combinedRepo{catRepo: mockCatRepo}
	svc := NewService(repo)

	categories := []models.Category{
		{ID: 1, Name: "Food"},
	}

	mockCatRepo.EXPECT().GetCategoriesByUser(gomock.Any(), 1).Return(categories, nil)
	mockCatRepo.EXPECT().GetCategoryStats(gomock.Any(), 1, 1).Return(5, nil)

	result, err := svc.GetCategoriesByUser(context.Background(), 1)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 1, result[0].ID)
	assert.Equal(t, 5, result[0].OperationsCount)
}

func TestService_GetCategoryByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCatRepo := mocks.NewMockCategoryRepository(ctrl)
	repo := &combinedRepo{catRepo: mockCatRepo}
	svc := NewService(repo)

	category := models.Category{ID: 1, Name: "Food"}

	mockCatRepo.EXPECT().GetCategoryByID(gomock.Any(), 1, 1).Return(category, nil)
	mockCatRepo.EXPECT().GetCategoryStats(gomock.Any(), 1, 1).Return(3, nil)

	result, err := svc.GetCategoryByID(context.Background(), 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, 3, result.OperationsCount)
}

func TestService_UpdateCategory_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCatRepo := mocks.NewMockCategoryRepository(ctrl)
	repo := &combinedRepo{catRepo: mockCatRepo}
	svc := NewService(repo)

	name := "Updated"
	req := models.UpdateCategoryRequest{Name: &name}
	existing := models.Category{ID: 1, Name: "Old"}

	mockCatRepo.EXPECT().GetCategoryByID(gomock.Any(), 1, 1).Return(existing, nil)
	mockCatRepo.EXPECT().UpdateCategory(gomock.Any(), gomock.Any()).Return(nil)

	result, err := svc.UpdateCategory(context.Background(), req, 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", result.Name)
}

func TestService_DeleteCategory_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCatRepo := mocks.NewMockCategoryRepository(ctrl)
	repo := &combinedRepo{catRepo: mockCatRepo}
	svc := NewService(repo)

	mockCatRepo.EXPECT().DeleteCategory(gomock.Any(), 1, 1).Return(nil)

	err := svc.DeleteCategory(context.Background(), 1, 1)
	assert.NoError(t, err)
}

func TestService_GetCategoriesByUser_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCatRepo := mocks.NewMockCategoryRepository(ctrl)
	repo := &combinedRepo{catRepo: mockCatRepo}
	svc := NewService(repo)

	mockCatRepo.EXPECT().GetCategoriesByUser(gomock.Any(), 1).Return(nil, errors.New("db error"))

	result, err := svc.GetCategoriesByUser(context.Background(), 1)
	assert.Error(t, err)
	assert.Empty(t, result)
}
