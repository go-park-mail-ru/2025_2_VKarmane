package operation

import (
	"context"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const testUserID = 42

func contextWithUserID() context.Context {
	return context.WithValue(context.Background(), middleware.UserIDKey, testUserID)
}

func TestService_GetAccountOperations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := contextWithUserID()
	accID := 10

	mockSvc := mocks.NewMockOperationService(ctrl)

	expectedOps := []models.Operation{
		{ID: 1, AccountID: accID, Name: "TestOp"},
	}

	mockSvc.EXPECT().GetAccountOperations(gomock.Any(), accID).Return(expectedOps, nil)

	service := mockSvc

	result, err := service.GetAccountOperations(ctx, accID)
	assert.NoError(t, err)
	assert.Equal(t, expectedOps, result)

	mockSvc.EXPECT().GetAccountOperations(gomock.Any(), 999).Return([]models.Operation{}, assert.AnError)

	result, err = service.GetAccountOperations(ctx, 999)
	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestService_CreateOperation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := contextWithUserID()
	accID := 1

	mockSvc := mocks.NewMockOperationService(ctrl)

	categoryID := 2
	req := models.CreateOperationRequest{
		AccountID:   accID,
		CategoryID:  &categoryID,
		Type:        models.OperationExpense,
		Name:        "Lunch",
		Description: "Food",
		Sum:         250,
	}

	expectedOp := models.Operation{
		ID:        123,
		AccountID: accID,
		Name:      "Lunch",
		Status:    models.OperationFinished,
	}

	mockSvc.EXPECT().CreateOperation(gomock.Any(), req, accID).Return(expectedOp, nil)

	service := mockSvc

	result, err := service.CreateOperation(ctx, req, accID)
	assert.NoError(t, err)
	assert.Equal(t, expectedOp.ID, result.ID)
	assert.Equal(t, "Lunch", result.Name)
	assert.Equal(t, models.OperationFinished, result.Status)
	assert.WithinDuration(t, expectedOp.CreatedAt, result.CreatedAt, time.Second)
}

func TestService_UpdateOperation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := contextWithUserID()
	accID := 5
	opID := 42
	newName := "Updated"
	newSum := float64(1000)

	mockSvc := mocks.NewMockOperationService(ctrl)

	req := models.UpdateOperationRequest{
		Name: &newName,
		Sum:  &newSum,
	}

	expectedOp := models.Operation{
		ID:        opID,
		AccountID: accID,
		Name:      newName,
		Sum:       newSum,
	}

	mockSvc.EXPECT().UpdateOperation(gomock.Any(), req, accID, opID).Return(expectedOp, nil)

	result, err := mockSvc.UpdateOperation(ctx, req, accID, opID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", result.Name)
	assert.Equal(t, newSum, result.Sum)

	mockSvc.EXPECT().UpdateOperation(gomock.Any(), req, 999, opID).Return(models.Operation{}, assert.AnError)

	result, err = mockSvc.UpdateOperation(ctx, req, 999, opID)
	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestService_DeleteOperation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := contextWithUserID()
	accID := 3
	opID := 99

	mockSvc := mocks.NewMockOperationService(ctrl)

	expectedOp := models.Operation{
		ID:        opID,
		AccountID: accID,
	}

	mockSvc.EXPECT().DeleteOperation(gomock.Any(), accID, opID).Return(expectedOp, nil)

	result, err := mockSvc.DeleteOperation(ctx, accID, opID)
	assert.NoError(t, err)
	assert.Equal(t, opID, result.ID)
	
	mockSvc.EXPECT().DeleteOperation(gomock.Any(), 999, opID).Return(models.Operation{}, assert.AnError)

	result, err = mockSvc.DeleteOperation(ctx, 999, opID)
	assert.Error(t, err)
	assert.Empty(t, result)
}
