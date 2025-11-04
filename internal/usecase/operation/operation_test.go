package operation

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

func TestUseCase_GetAccountOperations_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockOperationService(ctrl)
	uc := NewUseCase(mockSvc)

	expectedOps := []models.Operation{
		{ID: 1, Name: "Test"},
	}

	mockSvc.EXPECT().GetAccountOperations(gomock.Any(), 1).Return(expectedOps, nil)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	result, err := uc.GetAccountOperations(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedOps, result)
}

func TestUseCase_GetAccountOperations_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockOperationService(ctrl)
	uc := NewUseCase(mockSvc)

	mockSvc.EXPECT().GetAccountOperations(gomock.Any(), 1).Return(nil, errors.New("database error"))

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	result, err := uc.GetAccountOperations(ctx, 1)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUseCase_CreateOperation_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockOperationService(ctrl)
	uc := NewUseCase(mockSvc)

	req := models.CreateOperationRequest{
		Name: "Test",
		Sum:  100,
	}
	expectedOp := models.Operation{ID: 1, Name: "Test", Sum: 100}

	mockSvc.EXPECT().CreateOperation(gomock.Any(), req, 1).Return(expectedOp, nil)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	result, err := uc.CreateOperation(ctx, req, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedOp, result)
}

func TestUseCase_UpdateOperation_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockOperationService(ctrl)
	uc := NewUseCase(mockSvc)

	name := "Updated"
	req := models.UpdateOperationRequest{Name: &name}
	expectedOp := models.Operation{ID: 1, Name: "Updated"}

	mockSvc.EXPECT().UpdateOperation(gomock.Any(), req, 1, 1).Return(expectedOp, nil)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	result, err := uc.UpdateOperation(ctx, req, 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedOp, result)
}

func TestUseCase_DeleteOperation_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockOperationService(ctrl)
	uc := NewUseCase(mockSvc)

	expectedOp := models.Operation{ID: 1, Status: models.OperationReverted}

	mockSvc.EXPECT().DeleteOperation(gomock.Any(), 1, 1).Return(expectedOp, nil)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	result, err := uc.DeleteOperation(ctx, 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedOp, result)
}

func TestUseCase_GetOperationByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockOperationService(ctrl)
	uc := NewUseCase(mockSvc)

	expectedOp := models.Operation{ID: 1, Name: "Test"}

	mockSvc.EXPECT().GetOperationByID(gomock.Any(), 1, 1).Return(expectedOp, nil)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	result, err := uc.GetOperationByID(ctx, 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedOp, result)
}
