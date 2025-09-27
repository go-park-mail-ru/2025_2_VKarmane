package balance

import (
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestUseCase_GetBalanceForUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         int
		mockAccounts   []models.Account
		mockError      error
		expectedResult []models.Account
		expectedError  error
	}{
		{
			name:   "successful get balance for user",
			userID: 1,
			mockAccounts: []models.Account{
				{
					ID:         1,
					Balance:    1000.50,
					Type:       "debit",
					CurrencyID: 1,
					CreatedAt:  time.Now(),
				},
				{
					ID:         2,
					Balance:    500.25,
					Type:       "credit",
					CurrencyID: 1,
					CreatedAt:  time.Now(),
				},
			},
			mockError: nil,
			expectedResult: []models.Account{
				{
					ID:         1,
					Balance:    1000.50,
					Type:       "debit",
					CurrencyID: 1,
					CreatedAt:  time.Now(),
				},
				{
					ID:         2,
					Balance:    500.25,
					Type:       "credit",
					CurrencyID: 1,
					CreatedAt:  time.Now(),
				},
			},
			expectedError: nil,
		},
		{
			name:           "service error",
			userID:         1,
			mockAccounts:   nil,
			mockError:      errors.New("database connection failed"),
			expectedResult: nil,
			expectedError:  errors.New("balance.GetBalanceForUser: database connection failed"),
		},
		{
			name:           "empty accounts list",
			userID:         1,
			mockAccounts:   []models.Account{},
			mockError:      nil,
			expectedResult: []models.Account{},
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBalanceService := &MockBalanceService{}
			uc := &UseCase{balanceSvc: mockBalanceService}

			mockBalanceService.EXPECT().
				GetBalanceForUser(tt.userID).
				Return(tt.mockAccounts, tt.mockError).
				Once()

			result, err := uc.GetBalanceForUser(tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedResult), len(result))
				for i, expectedAccount := range tt.expectedResult {
					assert.Equal(t, expectedAccount.ID, result[i].ID)
					assert.Equal(t, expectedAccount.Balance, result[i].Balance)
					assert.Equal(t, expectedAccount.Type, result[i].Type)
					assert.Equal(t, expectedAccount.CurrencyID, result[i].CurrencyID)
				}
			}

			mockBalanceService.AssertExpectations(t)
		})
	}
}

func TestUseCase_GetAccountByID(t *testing.T) {
	tests := []struct {
		name            string
		userID          int
		accountID       int
		mockAccounts    []models.Account
		mockError       error
		expectedAccount models.Account
		expectedError   error
	}{
		{
			name:      "successful get account by id",
			userID:    1,
			accountID: 1,
			mockAccounts: []models.Account{
				{
					ID:         1,
					Balance:    1000.50,
					Type:       "debit",
					CurrencyID: 1,
					CreatedAt:  time.Now(),
				},
				{
					ID:         2,
					Balance:    500.25,
					Type:       "credit",
					CurrencyID: 1,
					CreatedAt:  time.Now(),
				},
			},
			mockError: nil,
			expectedAccount: models.Account{
				ID:         1,
				Balance:    1000.50,
				Type:       "debit",
				CurrencyID: 1,
				CreatedAt:  time.Now(),
			},
			expectedError: nil,
		},
		{
			name:      "account not found",
			userID:    1,
			accountID: 999,
			mockAccounts: []models.Account{
				{
					ID:         1,
					Balance:    1000.50,
					Type:       "debit",
					CurrencyID: 1,
					CreatedAt:  time.Now(),
				},
			},
			mockError:       nil,
			expectedAccount: models.Account{},
			expectedError:   errors.New("balance.GetAccountByID: ACCOUNT_NOT_FOUND"),
		},
		{
			name:            "service error",
			userID:          1,
			accountID:       1,
			mockAccounts:    nil,
			mockError:       errors.New("database connection failed"),
			expectedAccount: models.Account{},
			expectedError:   errors.New("balance.GetAccountByID: database connection failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBalanceService := &MockBalanceService{}
			uc := &UseCase{balanceSvc: mockBalanceService}

			mockBalanceService.EXPECT().
				GetBalanceForUser(tt.userID).
				Return(tt.mockAccounts, tt.mockError).
				Once()

			account, err := uc.GetAccountByID(tt.userID, tt.accountID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAccount.ID, account.ID)
				assert.Equal(t, tt.expectedAccount.Balance, account.Balance)
				assert.Equal(t, tt.expectedAccount.Type, account.Type)
				assert.Equal(t, tt.expectedAccount.CurrencyID, account.CurrencyID)
			}

			mockBalanceService.AssertExpectations(t)
		})
	}
}
