package operation

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAccountOperations_Success(t *testing.T) {
	mockSvc := mocks.NewOperationService(t)
	uc := &UseCase{opSvc: mockSvc}

	expected := []models.Operation{{ID: 1, AccountID: 7, Name: "test"}}

	mockSvc.On("GetAccountOperations", mock.Anything, 7).
		Return(expected, nil).
		Once()

	ops, err := uc.GetAccountOperations(context.Background(), 7)
	assert.NoError(t, err)
	assert.Len(t, ops, 1)
	assert.Equal(t, 7, ops[0].AccountID)
	assert.Equal(t, "test", ops[0].Name)
}

func TestGetAccountOperations_Error(t *testing.T) {
	mockSvc := mocks.NewOperationService(t)
	uc := &UseCase{opSvc: mockSvc}

	mockSvc.On("GetAccountOperations", mock.Anything, 1).
		Return(nil, errors.New("db error")).
		Once()

	ops, err := uc.GetAccountOperations(context.Background(), 1)
	assert.Error(t, err)
	assert.Nil(t, ops)
}

func TestGetOperationByID_Success(t *testing.T) {
	mockSvc := mocks.NewOperationService(t)
	uc := &UseCase{opSvc: mockSvc}

	mockSvc.On("GetAccountOperations", mock.Anything, 10).
		Return([]models.Operation{
			{ID: 1, AccountID: 10, Name: "A"},
			{ID: 2, AccountID: 10, Name: "B"},
		}, nil).
		Once()

	op, err := uc.GetOperationByID(context.Background(), 10, 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, op.ID)
	assert.Equal(t, "B", op.Name)
}

func TestGetOperationByID_NotFound(t *testing.T) {
	mockSvc := mocks.NewOperationService(t)
	uc := &UseCase{opSvc: mockSvc}

	mockSvc.On("GetAccountOperations", mock.Anything, 10).
		Return([]models.Operation{{ID: 1}}, nil).
		Once()

	op, err := uc.GetOperationByID(context.Background(), 10, 999)
	assert.Error(t, err)
	assert.Equal(t, 0, op.ID)
}

func TestGetOperationByID_ErrorFromService(t *testing.T) {
	mockSvc := mocks.NewOperationService(t)
	uc := &UseCase{opSvc: mockSvc}

	mockSvc.On("GetAccountOperations", mock.Anything, 10).
		Return(nil, errors.New("db error")).
		Once()

	op, err := uc.GetOperationByID(context.Background(), 10, 1)
	assert.Error(t, err)
	assert.Equal(t, 0, op.ID)
}

func TestCreateOperation_Success(t *testing.T) {
	mockSvc := mocks.NewOperationService(t)
	uc := &UseCase{opSvc: mockSvc}

	req := models.CreateOperationRequest{Name: "test", CreatedAt: time.Now()}
	expected := models.Operation{ID: 42, AccountID: 5, Name: "test", CreatedAt: req.CreatedAt}

	mockSvc.On("CreateOperation", mock.Anything, req, 5).
		Return(expected, nil).
		Once()

	op, err := uc.CreateOperation(context.Background(), req, 5)
	assert.NoError(t, err)
	assert.Equal(t, 42, op.ID)
	assert.Equal(t, "test", op.Name)
	assert.Equal(t, req.CreatedAt, op.CreatedAt)
}

func TestCreateOperation_Error(t *testing.T) {
	mockSvc := mocks.NewOperationService(t)
	uc := &UseCase{opSvc: mockSvc}

	req := models.CreateOperationRequest{}
	mockSvc.On("CreateOperation", mock.Anything, req, 1).
		Return(models.Operation{}, errors.New("create failed")).
		Once()

	op, err := uc.CreateOperation(context.Background(), req, 1)
	assert.Error(t, err)
	assert.Equal(t, 0, op.ID)
}

func TestUpdateOperation_Success(t *testing.T) {
	mockSvc := mocks.NewOperationService(t)
	uc := &UseCase{opSvc: mockSvc}

	req := models.UpdateOperationRequest{Name: ptr("new")}
	expected := models.Operation{ID: 2, AccountID: 1, Name: "updated"}

	mockSvc.On("UpdateOperation", mock.Anything, req, 1, 2).
		Return(expected, nil).
		Once()

	op, err := uc.UpdateOperation(context.Background(), req, 1, 2)
	assert.NoError(t, err)
	assert.Equal(t, "updated", op.Name)
	assert.Equal(t, 2, op.ID)
}

func TestUpdateOperation_Error(t *testing.T) {
	mockSvc := mocks.NewOperationService(t)
	uc := &UseCase{opSvc: mockSvc}

	req := models.UpdateOperationRequest{}
	mockSvc.On("UpdateOperation", mock.Anything, req, 1, 2).
		Return(models.Operation{}, errors.New("update failed")).
		Once()

	op, err := uc.UpdateOperation(context.Background(), req, 1, 2)
	assert.Error(t, err)
	assert.Equal(t, 0, op.ID)
}

func TestDeleteOperation_Success(t *testing.T) {
	mockSvc := mocks.NewOperationService(t)
	uc := &UseCase{opSvc: mockSvc}

	expected := models.Operation{ID: 99, AccountID: 1}

	mockSvc.On("DeleteOperation", mock.Anything, 1, 99).
		Return(expected, nil).
		Once()

	op, err := uc.DeleteOperation(context.Background(), 1, 99)
	assert.NoError(t, err)
	assert.Equal(t, 99, op.ID)
	assert.Equal(t, 1, op.AccountID)
}

func TestDeleteOperation_Error(t *testing.T) {
	mockSvc := mocks.NewOperationService(t)
	uc := &UseCase{opSvc: mockSvc}

	mockSvc.On("DeleteOperation", mock.Anything, 1, 2).
		Return(models.Operation{}, errors.New("delete failed")).
		Once()

	op, err := uc.DeleteOperation(context.Background(), 1, 2)
	assert.Error(t, err)
	assert.Equal(t, 0, op.ID)
}

func ptr[T any](v T) *T {
	return &v
}
