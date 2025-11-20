package balance

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockBalanceService := mocks.NewMockBalanceService(ctrl)
			uc := &UseCase{balanceSvc: mockBalanceService}

			mockBalanceService.EXPECT().GetBalanceForUser(gomock.Any(), tt.userID).Return(tt.mockAccounts, tt.mockError)

			ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
			result, err := uc.GetBalanceForUser(ctx, tt.userID)

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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockBalanceService := mocks.NewMockBalanceService(ctrl)
			uc := &UseCase{balanceSvc: mockBalanceService}

			mockBalanceService.EXPECT().GetBalanceForUser(gomock.Any(), tt.userID).Return(tt.mockAccounts, tt.mockError)

			ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
			account, err := uc.GetAccountByID(ctx, tt.userID, tt.accountID)

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
		})
	}
}


func TestUseCase_CreateAccount(t *testing.T) {
	tests := []struct {
		name            string
		userID          int
		req             models.CreateAccountRequest
		mockAccount     models.Account
		mockError       error
		expectedAccount models.Account
		expectedError   error
	}{
		{
			name:   "successful create account",
			userID: 1,
			req: models.CreateAccountRequest{
				Balance:    100.0,
				Type:       models.PrivateAccount,
				CurrencyID: 1,
			},
			mockAccount: models.Account{
				ID:         10,
				Balance:    100.0,
				Type:       string(models.PrivateAccount),
				CurrencyID: 1,
			},
			mockError:       nil,
			expectedAccount: models.Account{ID: 10, Balance: 100.0, Type: string(models.PrivateAccount), CurrencyID: 1},
			expectedError:   nil,
		},
		{
			name:   "service error",
			userID: 1,
			req: models.CreateAccountRequest{
				Balance:    100.0,
				Type:       models.SharedAccount,
				CurrencyID: 2,
			},
			mockAccount:     models.Account{},
			mockError:       errors.New("failed to insert"),
			expectedAccount: models.Account{},
			expectedError:   errors.New("balance.CreateAccount: failed to insert"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockBalanceService := mocks.NewMockBalanceService(ctrl)
			uc := &UseCase{balanceSvc: mockBalanceService}

			mockBalanceService.EXPECT().CreateAccount(gomock.Any(), tt.req, tt.userID).Return(tt.mockAccount, tt.mockError)

			ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
			account, err := uc.CreateAccount(ctx, tt.req, tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAccount, account)
			}
		})
	}
}

func TestUseCase_UpdateAccount(t *testing.T) {
	tests := []struct {
		name            string
		userID          int
		accID           int
		req             models.UpdateAccountRequest
		mockAccount     models.Account
		mockError       error
		expectedAccount models.Account
		expectedError   error
	}{
		{
			name:   "successful update account",
			userID: 1,
			accID:  2,
			req: models.UpdateAccountRequest{
				Balance: 200.0,
			},
			mockAccount: models.Account{
				ID:         2,
				Balance:    200.0,
				Type:       string(models.PrivateAccount),
				CurrencyID: 1,
			},
			mockError:       nil,
			expectedAccount: models.Account{ID: 2, Balance: 200.0, Type: string(models.PrivateAccount), CurrencyID: 1},
			expectedError:   nil,
		},
		{
			name:   "service error",
			userID: 1,
			accID:  2,
			req: models.UpdateAccountRequest{
				Balance: 300.0,
			},
			mockAccount:     models.Account{},
			mockError:       errors.New("db update failed"),
			expectedAccount: models.Account{},
			expectedError:   errors.New("balance.UpdateAccount: db update failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockBalanceService := mocks.NewMockBalanceService(ctrl)
			uc := &UseCase{balanceSvc: mockBalanceService}

			mockBalanceService.EXPECT().UpdateAccount(gomock.Any(), tt.req, tt.userID, tt.accID).Return(tt.mockAccount, tt.mockError)

			ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
			account, err := uc.UpdateAccount(ctx, tt.req, tt.userID, tt.accID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAccount, account)
			}
		})
	}
}

func TestUseCase_DeleteAccount(t *testing.T) {
	tests := []struct {
		name            string
		userID          int
		accID           int
		mockAccount     models.Account
		mockError       error
		expectedAccount models.Account
		expectedError   error
	}{
		{
			name:   "successful delete account",
			userID: 1,
			accID:  10,
			mockAccount: models.Account{
				ID:         10,
				Balance:    0,
				Type:       string(models.PrivateAccount),
				CurrencyID: 1,
			},
			mockError:       nil,
			expectedAccount: models.Account{ID: 10, Balance: 0, Type: string(models.PrivateAccount), CurrencyID: 1},
			expectedError:   nil,
		},
		{
			name:            "service error",
			userID:          1,
			accID:           99,
			mockAccount:     models.Account{},
			mockError:       errors.New("delete failed"),
			expectedAccount: models.Account{},
			expectedError:   errors.New("balance.DeleteAccount: delete failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockBalanceService := mocks.NewMockBalanceService(ctrl)
			uc := &UseCase{balanceSvc: mockBalanceService}

			mockBalanceService.EXPECT().DeleteAccount(gomock.Any(), tt.userID, tt.accID).Return(tt.mockAccount, tt.mockError)

			ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
			account, err := uc.DeleteAccount(ctx, tt.userID, tt.accID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAccount, account)
			}
		})
	}
}
