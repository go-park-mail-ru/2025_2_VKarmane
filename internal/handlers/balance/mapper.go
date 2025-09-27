package balance

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

func AccountToAPI(account models.Account) AccountAPI {
	return AccountAPI{
		ID:         account.ID,
		Balance:    account.Balance,
		Type:       account.Type,
		CurrencyID: account.CurrencyID,
	}
}

func AccountsToBalanceAPI(userID int, accounts []models.Account) BalanceAPI {
	accountDTOs := make([]AccountAPI, 0, len(accounts))
	for _, account := range accounts {
		accountDTOs = append(accountDTOs, AccountToAPI(account))
	}

	return BalanceAPI{
		UserID:   userID,
		Accounts: accountDTOs,
	}
}
