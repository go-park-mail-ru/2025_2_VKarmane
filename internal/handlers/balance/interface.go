package balance

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type BalanceUseCase interface {
	GetBalanceForUser(ctx context.Context, userID int) ([]models.Account, error)
	GetAccountByID(ctx context.Context, userID, accountID int) (models.Account, error)
	CreateAccount(ctx context.Context, req models.CreateAccountRequest, userID int) (models.Account, error)
	UpdateAccount(ctx context.Context, req models.UpdateAccountRequest, userID, accID int) (models.Account, error)
	DeleteAccount(ctx context.Context, userID, accID int) (models.Account, error)
}
