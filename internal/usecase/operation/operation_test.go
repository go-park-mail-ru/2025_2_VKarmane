package operation

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/stretchr/testify/require"
)

type MockOperationService struct {
	GetAccountOps func(ctx context.Context, accID int) ([]models.Operation, error)
	CreateOp      func(ctx context.Context, req models.CreateOperationRequest, accID int) (models.Operation, error)
	UpdateOp      func(ctx context.Context, req models.UpdateOperationRequest, accID, opID int) (models.Operation, error)
	DeleteOp      func(ctx context.Context, accID, opID int) (models.Operation, error)
}

func (m *MockOperationService) GetAccountOperations(ctx context.Context, accID int) ([]models.Operation, error) {
	return m.GetAccountOps(ctx, accID)
}

func (m *MockOperationService) CreateOperation(ctx context.Context, req models.CreateOperationRequest, accID int) (models.Operation, error) {
	return m.CreateOp(ctx, req, accID)
}

func (m *MockOperationService) UpdateOperation(ctx context.Context, req models.UpdateOperationRequest, accID, opID int) (models.Operation, error) {
	return m.UpdateOp(ctx, req, accID, opID)
}

func (m *MockOperationService) DeleteOperation(ctx context.Context, accID, opID int) (models.Operation, error) {
	return m.DeleteOp(ctx, accID, opID)
}

func TestGetAccountOperations_Success(t *testing.T) {
	mockSvc := &MockOperationService{
		GetAccountOps: func(ctx context.Context, accID int) ([]models.Operation, error) {
			return []models.Operation{{ID: 1, AccountID: accID, Name: "test"}}, nil
		},
	}
	uc := &UseCase{opSvc: mockSvc}

	ops, err := uc.GetAccountOperations(context.Background(), 7)
	require.NoError(t, err)
	require.Len(t, ops, 1)
	require.Equal(t, 7, ops[0].AccountID)
	require.Equal(t, "test", ops[0].Name)
}

func TestGetAccountOperations_Error(t *testing.T) {
	mockSvc := &MockOperationService{
		GetAccountOps: func(ctx context.Context, accID int) ([]models.Operation, error) {
			return nil, errors.New("db error")
		},
	}
	uc := &UseCase{opSvc: mockSvc}

	ops, err := uc.GetAccountOperations(context.Background(), 1)
	require.Error(t, err)
	require.Nil(t, ops)
}

func TestGetOperationByID_Success(t *testing.T) {
	mockSvc := &MockOperationService{
		GetAccountOps: func(ctx context.Context, accID int) ([]models.Operation, error) {
			return []models.Operation{
				{ID: 1, AccountID: accID, Name: "A"},
				{ID: 2, AccountID: accID, Name: "B"},
			}, nil
		},
	}
	uc := &UseCase{opSvc: mockSvc}

	op, err := uc.GetOperationByID(context.Background(), 10, 2)
	require.NoError(t, err)
	require.Equal(t, 2, op.ID)
	require.Equal(t, "B", op.Name)
}

func TestGetOperationByID_NotFound(t *testing.T) {
	mockSvc := &MockOperationService{
		GetAccountOps: func(ctx context.Context, accID int) ([]models.Operation, error) {
			return []models.Operation{{ID: 1}}, nil
		},
	}
	uc := &UseCase{opSvc: mockSvc}

	op, err := uc.GetOperationByID(context.Background(), 10, 999)
	require.Error(t, err)
	require.Equal(t, 0, op.ID)
}

func TestGetOperationByID_ErrorFromService(t *testing.T) {
	mockSvc := &MockOperationService{
		GetAccountOps: func(ctx context.Context, accID int) ([]models.Operation, error) {
			return nil, errors.New("db error")
		},
	}
	uc := &UseCase{opSvc: mockSvc}

	op, err := uc.GetOperationByID(context.Background(), 10, 1)
	require.Error(t, err)
	require.Equal(t, 0, op.ID)
}

func TestCreateOperation_Success(t *testing.T) {
	mockSvc := &MockOperationService{
		CreateOp: func(ctx context.Context, req models.CreateOperationRequest, accID int) (models.Operation, error) {
			return models.Operation{ID: 42, AccountID: accID, Name: req.Name}, nil
		},
	}
	uc := &UseCase{opSvc: mockSvc}

	req := models.CreateOperationRequest{Name: "test"}
	op, err := uc.CreateOperation(context.Background(), req, 5)
	require.NoError(t, err)
	require.Equal(t, 42, op.ID)
	require.Equal(t, "test", op.Name)
}

func TestCreateOperation_Error(t *testing.T) {
	mockSvc := &MockOperationService{
		CreateOp: func(ctx context.Context, req models.CreateOperationRequest, accID int) (models.Operation, error) {
			return models.Operation{}, errors.New("create failed")
		},
	}
	uc := &UseCase{opSvc: mockSvc}

	op, err := uc.CreateOperation(context.Background(), models.CreateOperationRequest{}, 1)
	require.Error(t, err)
	require.Equal(t, 0, op.ID)
}

func TestUpdateOperation_Success(t *testing.T) {
	mockSvc := &MockOperationService{
		UpdateOp: func(ctx context.Context, req models.UpdateOperationRequest, accID, opID int) (models.Operation, error) {
			return models.Operation{ID: opID, AccountID: accID, Name: "updated"}, nil
		},
	}
	uc := &UseCase{opSvc: mockSvc}

	req := models.UpdateOperationRequest{Name: ptr("new")}
	op, err := uc.UpdateOperation(context.Background(), req, 1, 2)
	require.NoError(t, err)
	require.Equal(t, "updated", op.Name)
	require.Equal(t, 2, op.ID)
}

func TestUpdateOperation_Error(t *testing.T) {
	mockSvc := &MockOperationService{
		UpdateOp: func(ctx context.Context, req models.UpdateOperationRequest, accID, opID int) (models.Operation, error) {
			return models.Operation{}, errors.New("update failed")
		},
	}
	uc := &UseCase{opSvc: mockSvc}

	op, err := uc.UpdateOperation(context.Background(), models.UpdateOperationRequest{}, 1, 2)
	require.Error(t, err)
	require.Equal(t, 0, op.ID)
}

func TestDeleteOperation_Success(t *testing.T) {
	mockSvc := &MockOperationService{
		DeleteOp: func(ctx context.Context, accID, opID int) (models.Operation, error) {
			return models.Operation{ID: opID, AccountID: accID}, nil
		},
	}
	uc := &UseCase{opSvc: mockSvc}

	op, err := uc.DeleteOperation(context.Background(), 1, 99)
	require.NoError(t, err)
	require.Equal(t, 99, op.ID)
	require.Equal(t, 1, op.AccountID)
}

func TestDeleteOperation_Error(t *testing.T) {
	mockSvc := &MockOperationService{
		DeleteOp: func(ctx context.Context, accID, opID int) (models.Operation, error) {
			return models.Operation{}, errors.New("delete failed")
		},
	}
	uc := &UseCase{opSvc: mockSvc}

	op, err := uc.DeleteOperation(context.Background(), 1, 2)
	require.Error(t, err)
	require.Equal(t, 0, op.ID)
}

func ptr[T any](v T) *T {
	return &v
}
