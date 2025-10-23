package balance

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/balance"
)

type UseCase struct {
	balanceSvc BalanceService
}

func NewUseCase(store repository.Repository) *UseCase {
	accountRepoAdapter := balance.NewPostgresAccountRepositoryAdapter(store)
	balanceService := balance.NewService(accountRepoAdapter)

	return &UseCase{
		balanceSvc: balanceService,
	}
}

func (uc *UseCase) GetBalanceForUser(ctx context.Context, userID int) ([]models.Account, error) {
	accounts, err := uc.balanceSvc.GetBalanceForUser(ctx, userID)
	if err != nil {
		if log := logger.FromContext(ctx); log != nil {
			log.Error("Failed to get balance for user", "error", err, "user_id", userID)
		}

		return nil, fmt.Errorf("balance.GetBalanceForUser: %w", err)
	}

	return accounts, nil
}

func (uc *UseCase) GetAccountByID(ctx context.Context, userID, accountID int) (models.Account, error) {
	accounts, err := uc.balanceSvc.GetBalanceForUser(ctx, userID)
	if err != nil {
		if log := logger.FromContext(ctx); log != nil {
			log.Error("Failed to get balance for user", "error", err, "user_id", userID)
		}

		return models.Account{}, fmt.Errorf("balance.GetAccountByID: %w", err)
	}

	for _, account := range accounts {
		if account.ID == accountID {
			return account, nil
		}
	}

	if log := logger.FromContext(ctx); log != nil {
		log.Warn("Account not found", "user_id", userID, "account_id", accountID)
	}

	return models.Account{}, fmt.Errorf("balance.GetAccountByID: %s", models.ErrCodeAccountNotFound)
}
