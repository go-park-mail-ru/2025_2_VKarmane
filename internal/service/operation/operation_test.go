package operation

import (
	"context"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/stretchr/testify/require"
)

type mockAccountRepo struct {
	getAccountsByUser func(ctx context.Context, userID int) []models.Account
}

func (m *mockAccountRepo) GetAccountsByUser(ctx context.Context, userID int) []models.Account {
	if m.getAccountsByUser != nil {
		return m.getAccountsByUser(ctx, userID)
	}
	return []models.Account{}
}

type mockOperationRepo struct {
	getOpsByAccountCalled bool
	createCalled          bool
	updateCalled          bool
	deleteCalled          bool

	getOpsByAccount func(ctx context.Context, accID int) []models.Operation
	createOp        func(ctx context.Context, op models.Operation) models.Operation
	updateOp        func(ctx context.Context, req models.UpdateOperationRequest, accID, opID int) models.Operation
	deleteOp        func(ctx context.Context, accID, opID int) models.Operation
}

func (m *mockOperationRepo) GetOperationsByAccount(ctx context.Context, accID int) []models.Operation {
	m.getOpsByAccountCalled = true
	if m.getOpsByAccount != nil {
		return m.getOpsByAccount(ctx, accID)
	}
	return nil
}

func (m *mockOperationRepo) GetOperationByID(ctx context.Context, accID int, opID int) models.Operation {
	if m.getOpsByAccount != nil {
		ops := m.getOpsByAccount(ctx, accID)
		for _, op := range ops {
			if op.ID == opID {
				return op
			}
		}
	}
	return models.Operation{}
}

func (m *mockOperationRepo) CreateOperation(ctx context.Context, op models.Operation) models.Operation {
	m.createCalled = true
	if m.createOp != nil {
		return m.createOp(ctx, op)
	}
	return op
}

func (m *mockOperationRepo) UpdateOperation(ctx context.Context, req models.UpdateOperationRequest, accID, opID int) models.Operation {
	m.updateCalled = true
	if m.updateOp != nil {
		return m.updateOp(ctx, req, accID, opID)
	}
	return models.Operation{}
}

func (m *mockOperationRepo) DeleteOperation(ctx context.Context, accID, opID int) models.Operation {
	m.deleteCalled = true
	if m.deleteOp != nil {
		return m.deleteOp(ctx, accID, opID)
	}
	return models.Operation{}
}

const testUserID = 42

func contextWithUserID() context.Context {
	return context.WithValue(context.Background(), middleware.UserIDKey, testUserID)
}


func TestGetAccountOperations_SuccessAndForbidden(t *testing.T) {
	accID := 10

	mockAccRepo := &mockAccountRepo{
		getAccountsByUser: func(ctx context.Context, userID int) []models.Account {
			if userID == testUserID {
				return []models.Account{{ID: accID}}
			}
			return []models.Account{}
		},
	}

	mockOpRepo := &mockOperationRepo{
		getOpsByAccount: func(ctx context.Context, accID int) []models.Operation {
			return []models.Operation{{ID: 1, AccountID: accID, Name: "TestOp"}}
		},
	}

	service := NewService(mockAccRepo, mockOpRepo)

	ops, err := service.GetAccountOperations(contextWithUserID(), accID)
	require.NoError(t, err)
	require.Len(t, ops, 1)
	require.Equal(t, "TestOp", ops[0].Name)

	ops, err = service.GetAccountOperations(contextWithUserID(), 999)
	require.Error(t, err)
	require.Contains(t, err.Error(), "forbidden")
	require.Len(t, ops, 0)
}

func TestCreateOperation_Success(t *testing.T) {
	accID := 1

	mockAccRepo := &mockAccountRepo{
		getAccountsByUser: func(ctx context.Context, userID int) []models.Account {
			return []models.Account{{ID: accID}}
		},
	}

	mockOpRepo := &mockOperationRepo{
		createOp: func(ctx context.Context, op models.Operation) models.Operation {
			op.ID = 123
			return op
		},
	}

	service := NewService(mockAccRepo, mockOpRepo)

	req := models.CreateOperationRequest{
		AccountID:  accID,
		CategoryID: 2,
		Type:       models.OperationExpense,
		Name:       "Lunch",
		Description: "Food",
		Sum:        250,
	}

	op, err := service.CreateOperation(contextWithUserID(), req, accID)
	require.NoError(t, err)
	require.True(t, mockOpRepo.createCalled)
	require.Equal(t, 123, op.ID)
	require.Equal(t, "Lunch", op.Name)
	require.Equal(t, models.OperationFinished, op.Status)
	require.WithinDuration(t, time.Now(), op.CreatedAt, time.Second)
}

func TestUpdateOperation_SuccessAndForbidden(t *testing.T) {
	accID := 5
	opID := 42
	newName := "Updated"
	newSum := float64(1000)

	mockAccRepo := &mockAccountRepo{
		getAccountsByUser: func(ctx context.Context, userID int) []models.Account {
			return []models.Account{{ID: accID}}
		},
	}

	mockOpRepo := &mockOperationRepo{
		updateOp: func(ctx context.Context, req models.UpdateOperationRequest, accID, opID int) models.Operation {
			return models.Operation{
				ID:        opID,
				AccountID: accID,
				Name:      *req.Name,
				Sum:       *req.Sum,
			}
		},
	}

	service := NewService(mockAccRepo, mockOpRepo)

	req := models.UpdateOperationRequest{
		Name: &newName,
		Sum:  &newSum,
	}
	op, err := service.UpdateOperation(contextWithUserID(), req, accID, opID)
	require.NoError(t, err)
	require.Equal(t, "Updated", op.Name)

	op, err = service.UpdateOperation(contextWithUserID(), req, 999, opID)
	require.Error(t, err)
	require.Contains(t, err.Error(), "forbidden")
}

func TestDeleteOperation_SuccessAndForbidden(t *testing.T) {
	accID := 3
	opID := 99

	mockAccRepo := &mockAccountRepo{
		getAccountsByUser: func(ctx context.Context, userID int) []models.Account {
			return []models.Account{{ID: accID}}
		},
	}

	mockOpRepo := &mockOperationRepo{
		deleteOp: func(ctx context.Context, accID, opID int) models.Operation {
			return models.Operation{ID: opID, AccountID: accID}
		},
	}

	service := NewService(mockAccRepo, mockOpRepo)

	op, err := service.DeleteOperation(contextWithUserID(), accID, opID)
	require.NoError(t, err)
	require.Equal(t, opID, op.ID)

	op, err = service.DeleteOperation(contextWithUserID(), 999, opID)
	require.Error(t, err)
	require.Contains(t, err.Error(), "forbidden")
}
