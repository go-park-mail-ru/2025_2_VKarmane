package balance

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type BalanceUseCase interface {
	GetBalanceForUser(ctx context.Context, userID int) ([]models.Account, error)
	GetAccountByID(ctx context.Context, userID, accountID int) (models.Account, error)
}
