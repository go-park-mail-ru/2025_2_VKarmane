package balance

import (
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestAccountToAPI(t *testing.T) {
	m := models.Account{ID: 3, Balance: 10.5, Type: "cash", CurrencyID: 1}
	api := AccountToAPI(m)
	assert.Equal(t, 3, api.ID)
	assert.Equal(t, 10.5, api.Balance)
	assert.Equal(t, "cash", api.Type)
	assert.Equal(t, 1, api.CurrencyID)
}

func TestAccountsToBalanceAPI(t *testing.T) {
	accs := []models.Account{{ID: 1}, {ID: 2}}
	dto := AccountsToBalanceAPI(7, accs)
	assert.Equal(t, 7, dto.UserID)
	assert.Len(t, dto.Accounts, 2)
	assert.Equal(t, 2, dto.Accounts[1].ID)
}
