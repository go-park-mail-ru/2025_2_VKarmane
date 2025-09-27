package balance

import "github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"

type BalanceService interface {
	GetBalanceForUser(userID int) ([]models.Account, error)
}

type AccountRepository interface {
	GetAccountsByUser(userID int) []models.Account
}
