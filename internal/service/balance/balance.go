package balance

import (
	"context"

	pkgerrors "github.com/pkg/errors"

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
		return []models.Account{}, pkgerrors.Wrap(err, "Failed to get balance for user")
	}

	return accounts, nil
}
