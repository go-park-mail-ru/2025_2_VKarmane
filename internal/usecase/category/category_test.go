package category

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUseCase_CreateCategory_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockCategoryService(ctrl)
	uc := NewUseCase(mockSvc)

	req := models.CreateCategoryRequest{
		Name:        "Food",
		Description: "Food expenses",
	}
	expectedCat := models.Category{ID: 1, Name: "Food", Description: "Food expenses"}

	mockSvc.EXPECT().CreateCategory(gomock.Any(), req, 1).Return(expectedCat, nil)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	result, err := uc.CreateCategory(ctx, req, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedCat, result)
}

func TestUseCase_CreateCategory_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockCategoryService(ctrl)
	uc := NewUseCase(mockSvc)

	req := models.CreateCategoryRequest{Name: "Food"}

	mockSvc.EXPECT().CreateCategory(gomock.Any(), req, 1).Return(models.Category{}, errors.New("db error"))

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	result, err := uc.CreateCategory(ctx, req, 1)
	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestUseCase_GetCategoriesByUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockCategoryService(ctrl)
	uc := NewUseCase(mockSvc)

	expectedCats := []models.CategoryWithStats{
		{Category: models.Category{ID: 1, Name: "Food"}, OperationsCount: 0},
	}

	mockSvc.EXPECT().GetCategoriesByUser(gomock.Any(), 1).Return(expectedCats, nil)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	result, err := uc.GetCategoriesByUser(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedCats, result)
}

func TestUseCase_GetCategoryByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockCategoryService(ctrl)
	uc := NewUseCase(mockSvc)

	expectedCat := models.CategoryWithStats{Category: models.Category{ID: 1, Name: "Food"}, OperationsCount: 0}

	mockSvc.EXPECT().GetCategoryByID(gomock.Any(), 1, 1).Return(expectedCat, nil)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	result, err := uc.GetCategoryByID(ctx, 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedCat, result)
}

func TestUseCase_UpdateCategory_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockCategoryService(ctrl)
	uc := NewUseCase(mockSvc)

	name := "Updated"
	req := models.UpdateCategoryRequest{Name: &name}
	expectedCat := models.Category{ID: 1, Name: "Updated"}

	mockSvc.EXPECT().UpdateCategory(gomock.Any(), req, 1, 1).Return(expectedCat, nil)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	result, err := uc.UpdateCategory(ctx, req, 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedCat, result)
}

func TestUseCase_DeleteCategory_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockCategoryService(ctrl)
	uc := NewUseCase(mockSvc)

	mockSvc.EXPECT().DeleteCategory(gomock.Any(), 1, 1).Return(nil)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	err := uc.DeleteCategory(ctx, 1, 1)
	assert.NoError(t, err)
}

