package budget

import (
	"context"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// combinedRepo is a helper struct to combine three mocked repositories into one
type combinedRepo struct {
	budgetRepo    *mocks.MockBudgetRepository
	accountRepo   *mocks.MockAccountRepository
	operationRepo *mocks.MockOperationRepository
}

func (c *combinedRepo) GetBudgetsByUser(ctx context.Context, userID int) ([]models.Budget, error) {
	return c.budgetRepo.GetBudgetsByUser(ctx, userID)
}

func (c *combinedRepo) GetAccountsByUser(ctx context.Context, userID int) ([]models.Account, error) {
	return c.accountRepo.GetAccountsByUser(ctx, userID)
}

func (c *combinedRepo) GetOperationsByAccount(ctx context.Context, accountID int) ([]models.Operation, error) {
	return c.operationRepo.GetOperationsByAccount(ctx, accountID)
}

func TestService_GetBudgetsForUser(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	tests := []struct {
		name           string
		userID         int
		mockBudgets    []models.Budget
		mockAccounts   []models.Account
		mockOperations []models.Operation
		expectedResult []models.Budget
	}{
		{
			name:   "successful get budgets with calculated actual amounts",
			userID: 1,
			mockBudgets: []models.Budget{
				{
					ID:          1,
					UserID:      1,
					Amount:      5000.00,
					CurrencyID:  1,
					Description: "Monthly budget",
					CreatedAt:   time.Now(),
					PeriodStart: time.Now().AddDate(0, 0, -15),
					PeriodEnd:   time.Now().AddDate(0, 0, 15),
				},
			},
			mockAccounts: []models.Account{
				{
					ID:         1,
					Balance:    1000.00,
					Type:       "debit",
					CurrencyID: 1,
					CreatedAt:  time.Now(),
				},
			},
			mockOperations: []models.Operation{
				{
					ID:         1,
					AccountID:  1,
					Sum:        100.50,
					Type:       "expense",
					CurrencyID: 1,
					CreatedAt:  time.Now().AddDate(0, 0, -10),
				},
				{
					ID:         2,
					AccountID:  1,
					Sum:        200.75,
					Type:       "expense",
					CurrencyID: 1,
					CreatedAt:  time.Now().AddDate(0, 0, -5),
				},
				{
					ID:         3,
					AccountID:  1,
					Sum:        50.00,
					Type:       "income",
					CurrencyID: 1,
					CreatedAt:  time.Now().AddDate(0, 0, -3),
				},
			},
			expectedResult: []models.Budget{
				{
					ID:          1,
					UserID:      1,
					Amount:      5000.00,
					Actual:      301.25, // 100.50 + 200.75 (only expenses)
					CurrencyID:  1,
					Description: "Monthly budget",
					CreatedAt:   time.Now(),
					PeriodStart: time.Now().AddDate(0, 0, -15),
					PeriodEnd:   time.Now().AddDate(0, 0, 15),
				},
			},
		},
		{
			name:   "successful get budgets with no operations",
			userID: 1,
			mockBudgets: []models.Budget{
				{
					ID:          1,
					UserID:      1,
					Amount:      2000.00,
					CurrencyID:  1,
					Description: "Food budget",
					CreatedAt:   time.Now(),
					PeriodStart: time.Now().AddDate(0, 0, -10),
					PeriodEnd:   time.Now().AddDate(0, 0, 20),
				},
			},
			mockAccounts: []models.Account{
				{
					ID:         1,
					Balance:    1000.00,
					Type:       "debit",
					CurrencyID: 1,
					CreatedAt:  time.Now(),
				},
			},
			mockOperations: []models.Operation{},
			expectedResult: []models.Budget{
				{
					ID:          1,
					UserID:      1,
					Amount:      2000.00,
					Actual:      0.00, // No operations
					CurrencyID:  1,
					Description: "Food budget",
					CreatedAt:   time.Now(),
					PeriodStart: time.Now().AddDate(0, 0, -10),
					PeriodEnd:   time.Now().AddDate(0, 0, 20),
				},
			},
		},
		{
			name:   "successful get budgets with operations outside period",
			userID: 1,
			mockBudgets: []models.Budget{
				{
					ID:          1,
					UserID:      1,
					Amount:      1000.00,
					CurrencyID:  1,
					Description: "Weekly budget",
					CreatedAt:   time.Now(),
					PeriodStart: time.Now().AddDate(0, 0, -3),
					PeriodEnd:   time.Now().AddDate(0, 0, 4),
				},
			},
			mockAccounts: []models.Account{
				{
					ID:         1,
					Balance:    1000.00,
					Type:       "debit",
					CurrencyID: 1,
					CreatedAt:  time.Now(),
				},
			},
			mockOperations: []models.Operation{
				{
					ID:         1,
					AccountID:  1,
					Sum:        100.00,
					Type:       "expense",
					CurrencyID: 1,
					CreatedAt:  time.Now().AddDate(0, 0, -10),
				},
				{
					ID:         2,
					AccountID:  1,
					Sum:        50.00,
					Type:       "expense",
					CurrencyID: 1,
					CreatedAt:  time.Now().AddDate(0, 0, 10),
				},
			},
			expectedResult: []models.Budget{
				{
					ID:          1,
					UserID:      1,
					Amount:      1000.00,
					Actual:      0.00,
					CurrencyID:  1,
					Description: "Weekly budget",
					CreatedAt:   time.Now(),
					PeriodStart: time.Now().AddDate(0, 0, -3),
					PeriodEnd:   time.Now().AddDate(0, 0, 4),
				},
			},
		},
		{
			name:   "successful get budgets with different currency operations",
			userID: 1,
			mockBudgets: []models.Budget{
				{
					ID:          1,
					UserID:      1,
					Amount:      1000.00,
					CurrencyID:  1,
					Description: "USD budget",
					CreatedAt:   time.Now(),
					PeriodStart: time.Now().AddDate(0, 0, -10),
					PeriodEnd:   time.Now().AddDate(0, 0, 20),
				},
			},
			mockAccounts: []models.Account{
				{
					ID:         1,
					Balance:    1000.00,
					Type:       "debit",
					CurrencyID: 1,
					CreatedAt:  time.Now(),
				},
			},
			mockOperations: []models.Operation{
				{
					ID:         1,
					AccountID:  1,
					Sum:        100.00,
					Type:       "expense",
					CurrencyID: 1,
					CreatedAt:  time.Now(),
				},
				{
					ID:         2,
					AccountID:  1,
					Sum:        200.00,
					Type:       "expense",
					CurrencyID: 2,
					CreatedAt:  time.Now(),
				},
			},
			expectedResult: []models.Budget{
				{
					ID:          1,
					UserID:      1,
					Amount:      1000.00,
					Actual:      100.00,
					CurrencyID:  1,
					Description: "USD budget",
					CreatedAt:   time.Now(),
					PeriodStart: time.Now().AddDate(0, 0, -10),
					PeriodEnd:   time.Now().AddDate(0, 0, 20),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockBudgetRepo := mocks.NewMockBudgetRepository(ctrl)
			mockAccountRepo := mocks.NewMockAccountRepository(ctrl)
			mockOperationRepo := mocks.NewMockOperationRepository(ctrl)

			mockBudgetRepo.EXPECT().GetBudgetsByUser(gomock.Any(), tt.userID).Return(tt.mockBudgets, nil)
			mockAccountRepo.EXPECT().GetAccountsByUser(gomock.Any(), tt.userID).Return(tt.mockAccounts, nil)
			mockOperationRepo.EXPECT().GetOperationsByAccount(gomock.Any(), 1).Return(tt.mockOperations, nil)

			combinedRepo := &combinedRepo{
				budgetRepo:    mockBudgetRepo,
				accountRepo:   mockAccountRepo,
				operationRepo: mockOperationRepo,
			}

			service := NewService(combinedRepo, fixedClock)

			result, err := service.GetBudgetsForUser(context.Background(), tt.userID)

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
		})
	}
}

func TestService_GetBudgetsForUser_MultipleAccountsAggregation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	mockBudgetRepo := mocks.NewMockBudgetRepository(ctrl)
	mockAccountRepo := mocks.NewMockAccountRepository(ctrl)
	mockOperationRepo := mocks.NewMockOperationRepository(ctrl)

	now := time.Now()
	budgets := []models.Budget{{ID: 1, UserID: 1, Amount: 1000, CurrencyID: 1, PeriodStart: now.Add(-24 * time.Hour), PeriodEnd: now.Add(24 * time.Hour)}}
	accounts := []models.Account{{ID: 1, CurrencyID: 1}, {ID: 2, CurrencyID: 1}}
	ops1 := []models.Operation{{ID: 1, AccountID: 1, Type: "expense", Sum: 100, CurrencyID: 1, CreatedAt: now}}
	ops2 := []models.Operation{{ID: 2, AccountID: 2, Type: "expense", Sum: 50, CurrencyID: 1, CreatedAt: now}}

	mockBudgetRepo.EXPECT().GetBudgetsByUser(gomock.Any(), 1).Return(budgets, nil)
	mockAccountRepo.EXPECT().GetAccountsByUser(gomock.Any(), 1).Return(accounts, nil)
	mockOperationRepo.EXPECT().GetOperationsByAccount(gomock.Any(), 1).Return(ops1, nil)
	mockOperationRepo.EXPECT().GetOperationsByAccount(gomock.Any(), 2).Return(ops2, nil)

	combinedRepo := &combinedRepo{
		budgetRepo:    mockBudgetRepo,
		accountRepo:   mockAccountRepo,
		operationRepo: mockOperationRepo,
	}

	svc := NewService(combinedRepo, fixedClock)
	res, err := svc.GetBudgetsForUser(context.Background(), 1)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, 150.0, res[0].Actual)
}

func TestService_GetBudgetsForUser_EmptyBudgets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	mockBudgetRepo := mocks.NewMockBudgetRepository(ctrl)
	mockAccountRepo := mocks.NewMockAccountRepository(ctrl)
	mockOperationRepo := mocks.NewMockOperationRepository(ctrl)

	mockBudgetRepo.EXPECT().GetBudgetsByUser(gomock.Any(), 1).Return([]models.Budget{}, nil)
	mockAccountRepo.EXPECT().GetAccountsByUser(gomock.Any(), 1).Return([]models.Account{}, nil)

	combinedRepo := &combinedRepo{
		budgetRepo:    mockBudgetRepo,
		accountRepo:   mockAccountRepo,
		operationRepo: mockOperationRepo,
	}

	svc := NewService(combinedRepo, fixedClock)
	res, err := svc.GetBudgetsForUser(context.Background(), 1)
	assert.NoError(t, err)
	assert.Empty(t, res)
}

func TestService_GetBudgetsForUser_WithIncome(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	mockBudgetRepo := mocks.NewMockBudgetRepository(ctrl)
	mockAccountRepo := mocks.NewMockAccountRepository(ctrl)
	mockOperationRepo := mocks.NewMockOperationRepository(ctrl)

	now := time.Now()
	budgets := []models.Budget{{ID: 1, UserID: 1, Amount: 1000, CurrencyID: 1, PeriodStart: now.Add(-24 * time.Hour), PeriodEnd: now.Add(24 * time.Hour)}}
	accounts := []models.Account{{ID: 1, CurrencyID: 1}}
	ops := []models.Operation{
		{ID: 1, AccountID: 1, Type: "expense", Sum: 100, CurrencyID: 1, CreatedAt: now},
		{ID: 2, AccountID: 1, Type: "income", Sum: 200, CurrencyID: 1, CreatedAt: now},
	}

	mockBudgetRepo.EXPECT().GetBudgetsByUser(gomock.Any(), 1).Return(budgets, nil)
	mockAccountRepo.EXPECT().GetAccountsByUser(gomock.Any(), 1).Return(accounts, nil)
	mockOperationRepo.EXPECT().GetOperationsByAccount(gomock.Any(), 1).Return(ops, nil)

	combinedRepo := &combinedRepo{
		budgetRepo:    mockBudgetRepo,
		accountRepo:   mockAccountRepo,
		operationRepo: mockOperationRepo,
	}

	svc := NewService(combinedRepo, fixedClock)
	res, err := svc.GetBudgetsForUser(context.Background(), 1)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, 100.0, res[0].Actual) // Only expense counted
}
