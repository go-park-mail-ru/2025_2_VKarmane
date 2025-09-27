package budget

import (
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestUseCase_GetBudgetsForUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         int
		mockBudgets    []models.Budget
		mockError      error
		expectedResult []models.Budget
		expectedError  error
	}{
		{
			name:   "successful get budgets for user",
			userID: 1,
			mockBudgets: []models.Budget{
				{
					ID:          1,
					UserID:      1,
					Amount:      5000.00,
					Actual:      2500.50,
					CurrencyID:  1,
					Description: "Monthly budget",
					CreatedAt:   time.Now(),
					PeriodStart: time.Now(),
					PeriodEnd:   time.Now().AddDate(0, 1, 0),
				},
				{
					ID:          2,
					UserID:      1,
					Amount:      2000.00,
					Actual:      800.25,
					CurrencyID:  1,
					Description: "Food budget",
					CreatedAt:   time.Now(),
					PeriodStart: time.Now(),
					PeriodEnd:   time.Now().AddDate(0, 1, 0),
				},
			},
			mockError: nil,
			expectedResult: []models.Budget{
				{
					ID:          1,
					UserID:      1,
					Amount:      5000.00,
					Actual:      2500.50,
					CurrencyID:  1,
					Description: "Monthly budget",
					CreatedAt:   time.Now(),
					PeriodStart: time.Now(),
					PeriodEnd:   time.Now().AddDate(0, 1, 0),
				},
				{
					ID:          2,
					UserID:      1,
					Amount:      2000.00,
					Actual:      800.25,
					CurrencyID:  1,
					Description: "Food budget",
					CreatedAt:   time.Now(),
					PeriodStart: time.Now(),
					PeriodEnd:   time.Now().AddDate(0, 1, 0),
				},
			},
			expectedError: nil,
		},
		{
			name:           "service error",
			userID:         1,
			mockBudgets:    nil,
			mockError:      errors.New("database connection failed"),
			expectedResult: nil,
			expectedError:  errors.New("budget.GetBudgetsForUser: database connection failed"),
		},
		{
			name:           "empty budgets list",
			userID:         1,
			mockBudgets:    []models.Budget{},
			mockError:      nil,
			expectedResult: []models.Budget{},
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBudgetService := &MockBudgetService{}
			uc := &UseCase{budgetSvc: mockBudgetService}

			mockBudgetService.EXPECT().
				GetBudgetsForUser(tt.userID).
				Return(tt.mockBudgets, tt.mockError).
				Once()

			result, err := uc.GetBudgetsForUser(tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedResult), len(result))
				for i, expectedBudget := range tt.expectedResult {
					assert.Equal(t, expectedBudget.ID, result[i].ID)
					assert.Equal(t, expectedBudget.UserID, result[i].UserID)
					assert.Equal(t, expectedBudget.Amount, result[i].Amount)
					assert.Equal(t, expectedBudget.Actual, result[i].Actual)
					assert.Equal(t, expectedBudget.CurrencyID, result[i].CurrencyID)
					assert.Equal(t, expectedBudget.Description, result[i].Description)
				}
			}

			mockBudgetService.AssertExpectations(t)
		})
	}
}

func TestUseCase_GetBudgetByID(t *testing.T) {
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
			name:     "successful get budget by id",
			userID:   1,
			budgetID: 1,
			mockBudgets: []models.Budget{
				{
					ID:          1,
					UserID:      1,
					Amount:      5000.00,
					Actual:      2500.50,
					CurrencyID:  1,
					Description: "Monthly budget",
					CreatedAt:   time.Now(),
					PeriodStart: time.Now(),
					PeriodEnd:   time.Now().AddDate(0, 1, 0),
				},
				{
					ID:          2,
					UserID:      1,
					Amount:      2000.00,
					Actual:      800.25,
					CurrencyID:  1,
					Description: "Food budget",
					CreatedAt:   time.Now(),
					PeriodStart: time.Now(),
					PeriodEnd:   time.Now().AddDate(0, 1, 0),
				},
			},
			mockError: nil,
			expectedBudget: models.Budget{
				ID:          1,
				UserID:      1,
				Amount:      5000.00,
				Actual:      2500.50,
				CurrencyID:  1,
				Description: "Monthly budget",
				CreatedAt:   time.Now(),
				PeriodStart: time.Now(),
				PeriodEnd:   time.Now().AddDate(0, 1, 0),
			},
			expectedError: nil,
		},
		{
			name:     "budget not found",
			userID:   1,
			budgetID: 999,
			mockBudgets: []models.Budget{
				{
					ID:          1,
					UserID:      1,
					Amount:      5000.00,
					Actual:      2500.50,
					CurrencyID:  1,
					Description: "Monthly budget",
					CreatedAt:   time.Now(),
					PeriodStart: time.Now(),
					PeriodEnd:   time.Now().AddDate(0, 1, 0),
				},
			},
			mockError:      nil,
			expectedBudget: models.Budget{},
			expectedError:  errors.New("budget.GetBudgetByID: BUDGET_NOT_FOUND"),
		},
		{
			name:           "service error",
			userID:         1,
			budgetID:       1,
			mockBudgets:    nil,
			mockError:      errors.New("database connection failed"),
			expectedBudget: models.Budget{},
			expectedError:  errors.New("budget.GetBudgetByID: database connection failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBudgetService := &MockBudgetService{}
			uc := &UseCase{budgetSvc: mockBudgetService}

			mockBudgetService.EXPECT().
				GetBudgetsForUser(tt.userID).
				Return(tt.mockBudgets, tt.mockError).
				Once()

			budget, err := uc.GetBudgetByID(tt.userID, tt.budgetID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBudget.ID, budget.ID)
				assert.Equal(t, tt.expectedBudget.UserID, budget.UserID)
				assert.Equal(t, tt.expectedBudget.Amount, budget.Amount)
				assert.Equal(t, tt.expectedBudget.Actual, budget.Actual)
				assert.Equal(t, tt.expectedBudget.CurrencyID, budget.CurrencyID)
				assert.Equal(t, tt.expectedBudget.Description, budget.Description)
			}

			mockBudgetService.AssertExpectations(t)
		})
	}
}
