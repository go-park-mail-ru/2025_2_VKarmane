package balance

import "github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"

type BalanceUseCase interface {
	GetBalanceForUser(userID int) ([]models.Account, error)
	GetAccountByID(userID, accountID int) (models.Account, error)
}
