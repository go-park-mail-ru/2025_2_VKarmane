package balance

import (
	"context"
	"fmt"

	pkgerrors "github.com/pkg/errors"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

type UseCase struct {
	balanceSvc BalanceService
	clock      clock.Clock
}

func NewUseCase(balanceService BalanceService) *UseCase {
	realClock := clock.RealClock{}
	return &UseCase{
		balanceSvc: balanceService,
		clock:      realClock,
	}
}

func (uc *UseCase) GetBalanceForUser(ctx context.Context, userID int) ([]models.Account, error) {
	log := logger.FromContext(ctx)
	accounts, err := uc.balanceSvc.GetBalanceForUser(ctx, userID)
	if err != nil {
		log.Error("Failed to get balance for user", "error", err, "user_id", userID)
		return nil, pkgerrors.Wrap(err, "balance.GetBalanceForUser")
	}

	return accounts, nil
}

func (uc *UseCase) GetAccountByID(ctx context.Context, userID, accountID int) (models.Account, error) {
	log := logger.FromContext(ctx)
	accounts, err := uc.balanceSvc.GetBalanceForUser(ctx, userID)
	if err != nil {
		log.Error("Failed to get balance for user", "error", err, "user_id", userID)
		return models.Account{}, pkgerrors.Wrap(err, "balance.GetAccountByID")
	}

	for _, account := range accounts {
		if account.ID == accountID {
			return account, nil
		}
	}

	log.Warn("Account not found", "user_id", userID, "account_id", accountID)
	return models.Account{}, fmt.Errorf("balance.GetAccountByID: %s", models.ErrCodeAccountNotFound)
}
