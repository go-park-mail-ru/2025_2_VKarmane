package operation

import (
	"context"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type Clock struct {
	FixedTime time.Time
}

func (f Clock) Now() time.Time {
	return f.FixedTime
}

const testUserID = 42

func contextWithUserID() context.Context {
	return context.WithValue(context.Background(), middleware.UserIDKey, testUserID)
}

func TestService_GetAccountOperations(t *testing.T) {
	ctx := contextWithUserID()
	accID := 10

	mockSvc := mocks.NewOperationService(t)

	expectedOps := []models.Operation{
		{ID: 1, AccountID: accID, Name: "TestOp"},
	}

	mockSvc.On("GetAccountOperations", mock.Anything, accID).Return(expectedOps, nil).Once()

	service := mockSvc

	result, err := service.GetAccountOperations(ctx, accID)
	assert.NoError(t, err)
	assert.Equal(t, expectedOps, result)

	mockSvc.On("GetAccountOperations", mock.Anything, 999).Return([]models.Operation{}, assert.AnError).Once()

	result, err = service.GetAccountOperations(ctx, 999)
	assert.Error(t, err)
	assert.Empty(t, result)

	mockSvc.AssertExpectations(t)
}

func TestService_CreateOperation(t *testing.T) {
	clockInstance := Clock{FixedTime: time.Now()}
	ctx := contextWithUserID()
	accID := 1

	mockSvc := mocks.NewOperationService(t)

	req := models.CreateOperationRequest{
		AccountID:   accID,
		CategoryID:  2,
		Type:        models.OperationExpense,
		Name:        "Lunch",
		Description: "Food",
		Sum:         250,
		CreatedAt:   clockInstance.Now(),
	}

	expectedOp := models.Operation{
		ID:        123,
		AccountID: accID,
		Name:      "Lunch",
		Status:    models.OperationFinished,
		CreatedAt: clockInstance.Now(),
	}

	mockSvc.On("CreateOperation", mock.Anything, req, accID).Return(expectedOp, nil).Once()

	service := mockSvc

	result, err := service.CreateOperation(ctx, req, accID)
	assert.NoError(t, err)
	assert.Equal(t, expectedOp.ID, result.ID)
	assert.Equal(t, "Lunch", result.Name)
	assert.Equal(t, models.OperationFinished, result.Status)
	assert.WithinDuration(t, expectedOp.CreatedAt, result.CreatedAt, time.Second)

	mockSvc.AssertExpectations(t)
}

func TestService_UpdateOperation(t *testing.T) {
	clockInstance := Clock{FixedTime: time.Now()}
	ctx := contextWithUserID()
	accID := 5
	opID := 42
	newName := "Updated"
	newSum := float64(1000)
	created_at := clockInstance.Now()

	mockSvc := mocks.NewOperationService(t)

	req := models.UpdateOperationRequest{
		Name:      &newName,
		Sum:       &newSum,
		CreatedAt: &created_at,
	}

	expectedOp := models.Operation{
		ID:        opID,
		AccountID: accID,
		Name:      newName,
		Sum:       newSum,
		CreatedAt: created_at,
	}

	mockSvc.On("UpdateOperation", mock.Anything, req, accID, opID).Return(expectedOp, nil).Once()

	result, err := mockSvc.UpdateOperation(ctx, req, accID, opID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", result.Name)
	assert.Equal(t, newSum, result.Sum)

	mockSvc.On("UpdateOperation", mock.Anything, req, 999, opID).Return(models.Operation{}, assert.AnError).Once()

	result, err = mockSvc.UpdateOperation(ctx, req, 999, opID)
	assert.Error(t, err)
	assert.Empty(t, result)

	mockSvc.AssertExpectations(t)
}

func TestService_DeleteOperation(t *testing.T) {
	ctx := contextWithUserID()
	accID := 3
	opID := 99

	mockSvc := mocks.NewOperationService(t)

	expectedOp := models.Operation{
		ID:        opID,
		AccountID: accID,
	}

	mockSvc.On("DeleteOperation", mock.Anything, accID, opID).Return(expectedOp, nil).Once()

	result, err := mockSvc.DeleteOperation(ctx, accID, opID)
	assert.NoError(t, err)
	assert.Equal(t, opID, result.ID)
	mockSvc.On("DeleteOperation", mock.Anything, 999, opID).Return(models.Operation{}, assert.AnError).Once()

	result, err = mockSvc.DeleteOperation(ctx, 999, opID)
	assert.Error(t, err)
	assert.Empty(t, result)

	mockSvc.AssertExpectations(t)
}
