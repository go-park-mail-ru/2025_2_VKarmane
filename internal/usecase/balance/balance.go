package balance

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/balance"
)

type UseCase struct {
	balanceSvc *balance.Service
}

func NewUseCase(store *repository.Store) *UseCase {
	return &UseCase{
		balanceSvc: balance.NewService(store),
	}
}

func (uc *UseCase) GetBalanceForUser(userID int) ([]models.Account, error) {
	accounts, err := uc.balanceSvc.GetBalanceForUser(userID)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (uc *UseCase) GetAccountByID(userID, accountID int) (models.Account, error) {
	accounts, err := uc.balanceSvc.GetBalanceForUser(userID)
	if err != nil {
		return models.Account{}, err
	}

	for _, account := range accounts {
		if account.ID == accountID {
			return account, nil
		}
	}

	return models.Account{}, models.ErrAccountNotFound
}
