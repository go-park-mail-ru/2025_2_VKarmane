package balance

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type BalanceService interface {
	GetBalanceForUser(ctx context.Context, userID int) ([]models.Account, error)
}

type AccountRepository interface {
	GetAccountsByUser(ctx context.Context, userID int) ([]models.Account, error)
}
