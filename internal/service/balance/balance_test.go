package balance

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

func TestService_GetBalanceForUser(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	tests := []struct {
		name           string
		userID         int
		mockAccounts   []models.Account
		expectedResult []models.Account
	}{
		{
			name:   "successful get balance for user with multiple accounts",
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
		},
		{
			name:           "successful get balance for user with no accounts",
			userID:         1,
			mockAccounts:   []models.Account{},
			expectedResult: []models.Account{},
		},
		{
			name:   "successful get balance for user with single account",
			userID: 1,
			mockAccounts: []models.Account{
				{
					ID:         1,
					Balance:    2500.75,
					Type:       "debit",
					CurrencyID: 1,
					CreatedAt:  time.Now(),
				},
			},
			expectedResult: []models.Account{
				{
					ID:         1,
					Balance:    2500.75,
					Type:       "debit",
					CurrencyID: 1,
					CreatedAt:  time.Now(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAccountRepo := mocks.NewMockAccountRepository(ctrl)
			mockAccountRepo.EXPECT().GetAccountsByUser(gomock.Any(), tt.userID).Return(tt.mockAccounts, nil)

			service := NewService(mockAccountRepo, fixedClock)

			result, err := service.GetBalanceForUser(context.Background(), tt.userID)

			assert.NoError(t, err)
			assert.Equal(t, len(tt.expectedResult), len(result))

			for i, expectedAccount := range tt.expectedResult {
				assert.Equal(t, expectedAccount.ID, result[i].ID)
				assert.Equal(t, expectedAccount.Balance, result[i].Balance)
				assert.Equal(t, expectedAccount.Type, result[i].Type)
				assert.Equal(t, expectedAccount.CurrencyID, result[i].CurrencyID)
			}
		})
	}
}

func TestService_GetBalanceForUser_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	repo := mocks.NewMockAccountRepository(ctrl)
	repo.EXPECT().GetAccountsByUser(gomock.Any(), 99).Return([]models.Account{}, nil)
	svc := NewService(repo, fixedClock)
	res, err := svc.GetBalanceForUser(context.Background(), 99)
	assert.NoError(t, err)
	assert.Empty(t, res)
}
