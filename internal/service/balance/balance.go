package balance

import (
	"context"
	"fmt"

	pkgerrors "github.com/pkg/errors"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

type Service struct {
	accountRepo AccountRepository
	clock       clock.Clock
}

func NewService(accountRepo AccountRepository, clck clock.Clock) *Service {
	return &Service{accountRepo: accountRepo, clock: clck}
}

func (s *Service) GetBalanceForUser(ctx context.Context, userID int) ([]models.Account, error) {
	accounts, err := s.accountRepo.GetAccountsByUser(ctx, userID)
	if err != nil {
		return []models.Account{}, pkgerrors.Wrap(err, "Failed to get balance for user")
	}

	return accounts, nil
}

func (s *Service) GetAccountByID(ctx context.Context, userID, accountID int) (models.Account, error) {
	accounts, err := s.accountRepo.GetAccountsByUser(ctx, userID)
	if err != nil {
		return models.Account{}, err
	}

	for _, account := range accounts {
		if account.ID == accountID {
			return account, nil
		}
	}

	return models.Account{}, fmt.Errorf("account not found: %s", models.ErrCodeAccountNotFound)
}
