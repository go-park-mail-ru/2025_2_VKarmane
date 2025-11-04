package account

import (
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestAccountDBToModel(t *testing.T) {
	now := time.Now()
	accountDB := AccountDB{
		ID:         1,
		Balance:    1000.50,
		Type:       "debit",
		CurrencyID: 1,
		CreatedAt:  now,
	}

	account := AccountDBToModel(accountDB)

	assert.Equal(t, 1, account.ID)
	assert.Equal(t, 1000.50, account.Balance)
	assert.Equal(t, "debit", account.Type)
	assert.Equal(t, 1, account.CurrencyID)
	assert.Equal(t, now, account.CreatedAt)
}

func TestAccountModelToDB(t *testing.T) {
	now := time.Now()
	account := models.Account{
		ID:         2,
		Balance:    500.25,
		Type:       "credit",
		CurrencyID: 2,
		CreatedAt:  now,
	}

	accountDB := AccountModelToDB(account)

	assert.Equal(t, 2, accountDB.ID)
	assert.Equal(t, 500.25, accountDB.Balance)
	assert.Equal(t, "credit", accountDB.Type)
	assert.Equal(t, 2, accountDB.CurrencyID)
	assert.Equal(t, now, accountDB.CreatedAt)
}
