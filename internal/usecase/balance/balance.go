package balance

import (
	"context"
	"fmt"

	pkgErrors "github.com/pkg/errors"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/account"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/balance"
)

type UseCase struct {
	balanceSvc BalanceService
}

func NewUseCase(store *repository.Store) *UseCase {
	accountRepo := account.NewRepository(store.Accounts, store.UserAccounts)
	balanceService := balance.NewService(accountRepo)

	return &UseCase{
		balanceSvc: balanceService,
	}
}

func (uc *UseCase) GetBalanceForUser(ctx context.Context, userID int) ([]models.Account, error) {
	log := logger.FromContext(ctx)
	accounts, err := uc.balanceSvc.GetBalanceForUser(ctx, userID)
	if err != nil {
		if log != nil {
			log.Error("Failed to get balance for user", "error", err, "user_id", userID)
		}

		return nil, pkgErrors.Wrap(err, "balance.GetBalanceForUser")
	}

	return accounts, nil
}

func (uc *UseCase) GetAccountByID(ctx context.Context, userID, accountID int) (models.Account, error) {
	log := logger.FromContext(ctx)
	accounts, err := uc.balanceSvc.GetBalanceForUser(ctx, userID)
	if err != nil {
		if log != nil {
			log.Error("Failed to get balance for user", "error", err, "user_id", userID)
		}

		return models.Account{}, pkgErrors.Wrap(err, "balance.GetAccountByID")
	}

	for _, account := range accounts {
		if account.ID == accountID {
			return account, nil
		}
	}

	if log != nil {
		log.Warn("Account not found", "user_id", userID, "account_id", accountID)
	}

	return models.Account{}, fmt.Errorf("balance.GetAccountByID: %s", models.ErrCodeAccountNotFound)
}
