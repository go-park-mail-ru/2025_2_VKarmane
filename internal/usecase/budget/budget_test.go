package budget

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	budgeterrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/budget"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUseCase_GetBudgetsForUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockBudgetService(ctrl)
	uc := &UseCase{budgetSvc: mockSvc}
	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())

	tests := []struct {
		name           string
		userID         int
		mockBudgets    []models.Budget
		mockError      error
		expectedResult []models.Budget
		expectedError  error
	}{
		{
			name:   "successful get budgets",
			userID: 1,
			mockBudgets: []models.Budget{
				{ID: 1, UserID: 1, Amount: 1000, Actual: 200},
				{ID: 2, UserID: 1, Amount: 500, Actual: 100},
			},
			mockError: nil,
			expectedResult: []models.Budget{
				{ID: 1, UserID: 1, Amount: 1000, Actual: 200},
				{ID: 2, UserID: 1, Amount: 500, Actual: 100},
			},
			expectedError: nil,
		},
		{
			name:           "service error",
			userID:         1,
			mockBudgets:    nil,
			mockError:      errors.New("db fail"),
			expectedResult: nil,
			expectedError:  errors.New("budget.GetBudgetsForUser"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc.EXPECT().GetBudgetsForUser(gomock.Any(), tt.userID).Return(tt.mockBudgets, tt.mockError)
			result, err := uc.GetBudgetsForUser(ctx, tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestUseCase_GetBudgetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockBudgetService(ctrl)
	uc := &UseCase{budgetSvc: mockSvc}
	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())

	tests := []struct {
		name           string
		userID         int
		budgetID       int
		mockBudgets    []models.Budget
		mockError      error
		expectedBudget models.Budget
		expectedError  error
	}{
		{
			name:     "successful get by ID",
			userID:   1,
			budgetID: 2,
			mockBudgets: []models.Budget{
				{ID: 1, UserID: 1, Amount: 1000, Actual: 200},
				{ID: 2, UserID: 1, Amount: 500, Actual: 100},
			},
			mockError:      nil,
			expectedBudget: models.Budget{ID: 2, UserID: 1, Amount: 500, Actual: 100},
			expectedError:  nil,
		},
		{
			name:     "budget not found",
			userID:   1,
			budgetID: 999,
			mockBudgets: []models.Budget{
				{ID: 1, UserID: 1, Amount: 1000, Actual: 200},
			},
			mockError:      nil,
			expectedBudget: models.Budget{},
			expectedError:  errors.New("budget.GetBudgetByID"),
		},
		{
			name:           "service error",
			userID:         1,
			budgetID:       1,
			mockBudgets:    nil,
			mockError:      errors.New("db fail"),
			expectedBudget: models.Budget{},
			expectedError:  errors.New("budget.GetBudgetByID"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc.EXPECT().GetBudgetsForUser(gomock.Any(), tt.userID).Return(tt.mockBudgets, tt.mockError)
			result, err := uc.GetBudgetByID(ctx, tt.userID, tt.budgetID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBudget, result)
			}
		})
	}
}

func TestUseCase_CreateBudget(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockBudgetService(ctrl)
	uc := &UseCase{budgetSvc: mockSvc}
	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())

	tests := []struct {
		name           string
		req            models.CreateBudgetRequest
		userID         int
		mockBudget     models.Budget
		mockError      error
		expectedBudget models.Budget
		expectedError  error
	}{
		{
			name:           "success",
			req:            models.CreateBudgetRequest{CategoryID: 1, Amount: 100},
			userID:         1,
			mockBudget:     models.Budget{ID: 10, UserID: 1, Amount: 100},
			mockError:      nil,
			expectedBudget: models.Budget{ID: 10, UserID: 1, Amount: 100},
			expectedError:  nil,
		},
		{
			name:           "active budget exists",
			req:            models.CreateBudgetRequest{CategoryID: 1, Amount: 100},
			userID:         1,
			mockBudget:     models.Budget{},
			mockError:      budgeterrors.ErrActiveBudgetExists,
			expectedBudget: models.Budget{},
			expectedError:  errors.New("budget.CreateBudget"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc.EXPECT().CreateBudget(gomock.Any(), tt.req, tt.userID).Return(tt.mockBudget, tt.mockError)
			result, err := uc.CreateBudget(ctx, tt.req, tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBudget, result)
			}
		})
	}
}

func TestUseCase_UpdateBudget(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockBudgetService(ctrl)
	uc := &UseCase{budgetSvc: mockSvc}
	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())

	sum := 200.5
	tests := []struct {
		name           string
		req            models.UpdatedBudgetRequest
		userID         int
		budgetID       int
		mockBudget     models.Budget
		mockError      error
		expectedBudget models.Budget
		expectedError  error
	}{
		{
			name:           "success",
			req:            models.UpdatedBudgetRequest{Amount: &sum},
			userID:         1,
			budgetID:       1,
			mockBudget:     models.Budget{ID: 1, UserID: 1, Amount: 200},
			mockError:      nil,
			expectedBudget: models.Budget{ID: 1, UserID: 1, Amount: 200},
			expectedError:  nil,
		},
		{
			name:           "service error",
			req:            models.UpdatedBudgetRequest{Amount: &sum},
			userID:         1,
			budgetID:       1,
			mockBudget:     models.Budget{},
			mockError:      errors.New("db fail"),
			expectedBudget: models.Budget{},
			expectedError:  errors.New("budget.UpdateBudget"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc.EXPECT().UpdateBudget(gomock.Any(), tt.req, tt.userID, tt.budgetID).Return(tt.mockBudget, tt.mockError)
			result, err := uc.UpdateBudget(ctx, tt.req, tt.userID, tt.budgetID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBudget, result)
			}
		})
	}
}

func TestUseCase_DeleteBudget(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockBudgetService(ctrl)
	uc := &UseCase{budgetSvc: mockSvc}
	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())

	tests := []struct {
		name           string
		userID         int
		budgetID       int
		mockBudget     models.Budget
		mockError      error
		expectedBudget models.Budget
		expectedError  error
	}{
		{
			name:           "success",
			userID:         1,
			budgetID:       1,
			mockBudget:     models.Budget{ID: 1, UserID: 1, Amount: 1000},
			mockError:      nil,
			expectedBudget: models.Budget{ID: 1, UserID: 1, Amount: 1000},
			expectedError:  nil,
		},
		{
			name:           "service error",
			userID:         1,
			budgetID:       1,
			mockBudget:     models.Budget{},
			mockError:      errors.New("db fail"),
			expectedBudget: models.Budget{},
			expectedError:  errors.New("budget.DeleteBudget"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc.EXPECT().DeleteBudget(gomock.Any(), tt.userID, tt.budgetID).Return(tt.mockBudget, tt.mockError)
			result, err := uc.DeleteBudget(ctx, tt.userID, tt.budgetID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBudget, result)
			}
		})
	}
}
