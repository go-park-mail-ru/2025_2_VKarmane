package operation

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type combinedRepo struct {
	opRepo  *mocks.MockOperationRepository
	accRepo *mocks.MockAccountRepository
}

func (c *combinedRepo) GetAccountsByUser(ctx context.Context, userID int) ([]models.Account, error) {
	return c.accRepo.GetAccountsByUser(ctx, userID)
}

func (c *combinedRepo) GetOperationsByAccount(ctx context.Context, accountID int) ([]models.OperationInList, error) {
	return c.opRepo.GetOperationsByAccount(ctx, accountID)
}

func (c *combinedRepo) GetOperationByID(ctx context.Context, accID int, opID int) (models.Operation, error) {
	return c.opRepo.GetOperationByID(ctx, accID, opID)
}

func (c *combinedRepo) CreateOperation(ctx context.Context, op models.Operation) (models.Operation, error) {
	return c.opRepo.CreateOperation(ctx, op)
}

func (c *combinedRepo) UpdateOperation(ctx context.Context, req models.UpdateOperationRequest, accID int, opID int) (models.Operation, error) {
	return c.opRepo.UpdateOperation(ctx, req, accID, opID)
}

func (c *combinedRepo) DeleteOperation(ctx context.Context, accID int, opID int) (models.Operation, error) {
	return c.opRepo.DeleteOperation(ctx, accID, opID)
}

func createContextWithUser(userID int) context.Context {
	ctx := context.Background()
	ctx = logger.WithLogger(ctx, logger.NewSlogLogger())
	ctx = context.WithValue(ctx, middleware.UserIDKey, userID)
	return ctx
}

func TestService_GetAccountOperations_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOpRepo := mocks.NewMockOperationRepository(ctrl)
	mockAccRepo := mocks.NewMockAccountRepository(ctrl)

	repo := &combinedRepo{opRepo: mockOpRepo, accRepo: mockAccRepo}
	svc := NewService(repo, clock.RealClock{})

	expectedOps := []models.OperationInList{
		{ID: 1, AccountID: 1, Name: "Test"},
	}

	accounts := []models.Account{{ID: 1}}
	mockAccRepo.EXPECT().GetAccountsByUser(gomock.Any(), 1).Return(accounts, nil)
	mockOpRepo.EXPECT().GetOperationsByAccount(gomock.Any(), 1).Return(expectedOps, nil)

	ctx := createContextWithUser(1)
	result, err := svc.GetAccountOperations(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedOps, result)
}

func TestService_GetAccountOperations_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOpRepo := mocks.NewMockOperationRepository(ctrl)
	mockAccRepo := mocks.NewMockAccountRepository(ctrl)

	repo := &combinedRepo{opRepo: mockOpRepo, accRepo: mockAccRepo}
	svc := NewService(repo, clock.RealClock{})

	accounts := []models.Account{{ID: 1}}
	mockAccRepo.EXPECT().GetAccountsByUser(gomock.Any(), 1).Return(accounts, nil)
	mockOpRepo.EXPECT().GetOperationsByAccount(gomock.Any(), 1).Return(nil, errors.New("db error"))

	ctx := createContextWithUser(1)
	result, err := svc.GetAccountOperations(ctx, 1)
	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestService_CreateOperation_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOpRepo := mocks.NewMockOperationRepository(ctrl)
	mockAccRepo := mocks.NewMockAccountRepository(ctrl)

	repo := &combinedRepo{opRepo: mockOpRepo, accRepo: mockAccRepo}
	fixedClock := clock.FixedClock{FixedTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}
	svc := NewService(repo, fixedClock)

	req := models.CreateOperationRequest{
		AccountID: 1,
		Name:      "Test",
		Sum:       100,
		Type:      models.OperationExpense,
	}

	expectedOp := models.Operation{
		ID:        1,
		AccountID: 1,
		Name:      "Test",
		Sum:       100,
		Type:      models.OperationExpense,
		Status:    models.OperationFinished,
	}

	accounts := []models.Account{{ID: 1}}
	mockAccRepo.EXPECT().GetAccountsByUser(gomock.Any(), 1).Return(accounts, nil)
	mockOpRepo.EXPECT().CreateOperation(gomock.Any(), gomock.Any()).Return(expectedOp, nil)

	ctx := createContextWithUser(1)
	result, err := svc.CreateOperation(ctx, req, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedOp.ID, result.ID)
	assert.Equal(t, expectedOp.Name, result.Name)
}

func TestService_GetOperationByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOpRepo := mocks.NewMockOperationRepository(ctrl)
	mockAccRepo := mocks.NewMockAccountRepository(ctrl)

	repo := &combinedRepo{opRepo: mockOpRepo, accRepo: mockAccRepo}
	svc := NewService(repo, clock.RealClock{})

	expectedOp := models.Operation{ID: 1, AccountID: 1, Name: "Test"}

	accounts := []models.Account{{ID: 1}}
	mockAccRepo.EXPECT().GetAccountsByUser(gomock.Any(), 1).Return(accounts, nil)
	mockOpRepo.EXPECT().GetOperationByID(gomock.Any(), 1, 1).Return(expectedOp, nil)

	ctx := createContextWithUser(1)
	result, err := svc.GetOperationByID(ctx, 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedOp, result)
}

func TestService_UpdateOperation_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOpRepo := mocks.NewMockOperationRepository(ctrl)
	mockAccRepo := mocks.NewMockAccountRepository(ctrl)

	repo := &combinedRepo{opRepo: mockOpRepo, accRepo: mockAccRepo}
	svc := NewService(repo, clock.RealClock{})

	name := "Updated"
	req := models.UpdateOperationRequest{Name: &name}
	expectedOp := models.Operation{ID: 1, Name: "Updated"}

	accounts := []models.Account{{ID: 1}}
	mockAccRepo.EXPECT().GetAccountsByUser(gomock.Any(), 1).Return(accounts, nil)
	mockOpRepo.EXPECT().UpdateOperation(gomock.Any(), req, 1, 1).Return(expectedOp, nil)

	ctx := createContextWithUser(1)
	result, err := svc.UpdateOperation(ctx, req, 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedOp, result)
}

func TestService_DeleteOperation_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOpRepo := mocks.NewMockOperationRepository(ctrl)
	mockAccRepo := mocks.NewMockAccountRepository(ctrl)

	repo := &combinedRepo{opRepo: mockOpRepo, accRepo: mockAccRepo}
	svc := NewService(repo, clock.RealClock{})

	expectedOp := models.Operation{ID: 1, Status: models.OperationReverted}

	accounts := []models.Account{{ID: 1}}
	mockAccRepo.EXPECT().GetAccountsByUser(gomock.Any(), 1).Return(accounts, nil)
	mockOpRepo.EXPECT().DeleteOperation(gomock.Any(), 1, 1).Return(expectedOp, nil)

	ctx := createContextWithUser(1)
	result, err := svc.DeleteOperation(ctx, 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedOp, result)
}

func TestService_CheckAccountOwnership_True(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOpRepo := mocks.NewMockOperationRepository(ctrl)
	mockAccRepo := mocks.NewMockAccountRepository(ctrl)

	repo := &combinedRepo{opRepo: mockOpRepo, accRepo: mockAccRepo}
	svc := NewService(repo, clock.RealClock{})

	accounts := []models.Account{{ID: 1}, {ID: 2}}
	mockAccRepo.EXPECT().GetAccountsByUser(gomock.Any(), 1).Return(accounts, nil)

	ctx := createContextWithUser(1)
	result := svc.CheckAccountOwnership(ctx, 1)
	assert.True(t, result)
}

func TestService_CheckAccountOwnership_False(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOpRepo := mocks.NewMockOperationRepository(ctrl)
	mockAccRepo := mocks.NewMockAccountRepository(ctrl)

	repo := &combinedRepo{opRepo: mockOpRepo, accRepo: mockAccRepo}
	svc := NewService(repo, clock.RealClock{})

	accounts := []models.Account{{ID: 1}, {ID: 2}}
	mockAccRepo.EXPECT().GetAccountsByUser(gomock.Any(), 1).Return(accounts, nil)

	ctx := createContextWithUser(1)
	result := svc.CheckAccountOwnership(ctx, 999)
	assert.False(t, result)
}
