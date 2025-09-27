package balance

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type Service struct {
	accountRepo AccountRepository
}

func NewService(accountRepo AccountRepository) *Service {
	return &Service{accountRepo: accountRepo}
}

func (s *Service) GetBalanceForUser(userID int) ([]models.Account, error) {
	accounts := s.accountRepo.GetAccountsByUser(userID)
	return accounts, nil
}
