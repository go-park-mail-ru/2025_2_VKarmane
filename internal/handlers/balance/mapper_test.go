package balance

import (
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestAccountToAPI(t *testing.T) {
	account := models.Account{
		ID:         1,
		Balance:    1000.50,
		Type:       "debit",
		CurrencyID: 1,
	}

	apiAccount := AccountToAPI(account)

	assert.Equal(t, 1, apiAccount.ID)
	assert.Equal(t, 1000.50, apiAccount.Balance)
	assert.Equal(t, "debit", apiAccount.Type)
	assert.Equal(t, 1, apiAccount.CurrencyID)
}

func TestAccountsToBalanceAPI(t *testing.T) {
	accounts := []models.Account{
		{ID: 1, Balance: 1000.00},
		{ID: 2, Balance: 500.00},
	}

	balanceAPI := AccountsToBalanceAPI(1, accounts)

	assert.Equal(t, 1, balanceAPI.UserID)
	assert.Len(t, balanceAPI.Accounts, 2)
	assert.Equal(t, 1, balanceAPI.Accounts[0].ID)
	assert.Equal(t, 2, balanceAPI.Accounts[1].ID)
}
