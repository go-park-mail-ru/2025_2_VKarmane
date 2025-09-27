package balance

import (
	"fmt"

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

func (uc *UseCase) GetBalanceForUser(userID int) ([]models.Account, error) {
	accounts, err := uc.balanceSvc.GetBalanceForUser(userID)
	if err != nil {
		return nil, fmt.Errorf("balance.GetBalanceForUser: %w", err)
	}

	return accounts, nil
}

func (uc *UseCase) GetAccountByID(userID, accountID int) (models.Account, error) {
	accounts, err := uc.balanceSvc.GetBalanceForUser(userID)
	if err != nil {
		return models.Account{}, fmt.Errorf("balance.GetAccountByID: %w", err)
	}

	for _, account := range accounts {
		if account.ID == accountID {
			return account, nil
		}
	}

	return models.Account{}, fmt.Errorf("balance.GetAccountByID: %s", models.ErrCodeAccountNotFound)
}
