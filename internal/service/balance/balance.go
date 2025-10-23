package balance

import (
	"context"
	"fmt"

	pkgErrors "github.com/pkg/errors"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type Service struct {
	accountRepo AccountRepository
}

func NewService(accountRepo AccountRepository) *Service {
	return &Service{accountRepo: accountRepo}
}

func (s *Service) GetBalanceForUser(ctx context.Context, userID int) ([]models.Account, error) {
	accounts, err := s.accountRepo.GetAccountsByUser(ctx, userID)
	if err != nil {
		return []models.Account{}, pkgErrors.Wrap(err, "Failed to get balance for user")
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
