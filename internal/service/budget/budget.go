package budget

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
)

type Service struct {
	store *repository.Store
}

func NewService(store *repository.Store) *Service {
	return &Service{store: store}
}

func (s *Service) GetBudgetsForUser(userID int) ([]models.Budget, error) {
	budgets := s.store.BudgetRepo.GetBudgetsByUser(userID)
	accounts := s.store.AccountRepo.GetAccountsByUser(userID)

	for i := range budgets {
		var actual float64
		for _, account := range accounts {
			ops := s.store.OperationRepo.GetOperationsByAccount(account.ID)
			for _, op := range ops {
				if op.CurrencyID != budgets[i].CurrencyID {
					continue
				}
				if op.CreatedAt.Before(budgets[i].PeriodStart) || op.CreatedAt.After(budgets[i].PeriodEnd) {
					continue
				}

				if op.Type == "expense" {
					actual += op.Sum
				}
			}
		}
		budgets[i].Actual = actual
	}

	return budgets, nil
}
