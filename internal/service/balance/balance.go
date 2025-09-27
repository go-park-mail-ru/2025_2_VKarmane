package balance

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type Service struct {
	accountRepo AccountRepository
}

func NewService(accountRepo AccountRepository) *Service {
	return &Service{accountRepo: accountRepo}
}

func (s *Service) GetBalanceForUser(ctx context.Context, userID int) ([]models.Account, error) {
	accounts := s.accountRepo.GetAccountsByUser(ctx, userID)

	return accounts, nil
}
