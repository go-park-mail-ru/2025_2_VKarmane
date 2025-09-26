package account

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

func AccountDBToModel(accountDB AccountDB) models.Account {
	return models.Account{
		ID:         accountDB.ID,
		Balance:    accountDB.Balance,
		Type:       accountDB.Type,
		CurrencyID: accountDB.CurrencyID,
		CreatedAt:  accountDB.CreatedAt,
		UpdatedAt:  accountDB.UpdatedAt,
	}
}

func AccountModelToDB(account models.Account) AccountDB {
	return AccountDB{
		ID:         account.ID,
		Balance:    account.Balance,
		Type:       account.Type,
		CurrencyID: account.CurrencyID,
		CreatedAt:  account.CreatedAt,
		UpdatedAt:  account.UpdatedAt,
	}
}

func UserAccountDBToModel(userAccountDB UserAccountDB) models.UserAccount {
	return models.UserAccount{
		ID:        userAccountDB.ID,
		UserID:    userAccountDB.UserID,
		AccountID: userAccountDB.AccountID,
		CreatedAt: userAccountDB.CreatedAt,
		UpdatedAt: userAccountDB.UpdatedAt,
	}
}

func UserAccountModelToDB(userAccount models.UserAccount) UserAccountDB {
	return UserAccountDB{
		ID:        userAccount.ID,
		UserID:    userAccount.UserID,
		AccountID: userAccount.AccountID,
		CreatedAt: userAccount.CreatedAt,
		UpdatedAt: userAccount.UpdatedAt,
	}
}
