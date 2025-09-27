package budget

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type Service struct {
	budgetRepo    BudgetRepository
	accountRepo   AccountRepository
	operationRepo OperationRepository
}

func NewService(budgetRepo BudgetRepository, accountRepo AccountRepository, operationRepo OperationRepository) *Service {
	return &Service{
		budgetRepo:    budgetRepo,
		accountRepo:   accountRepo,
		operationRepo: operationRepo,
	}
}

func (s *Service) GetBudgetsForUser(ctx context.Context, userID int) ([]models.Budget, error) {
	budgets := s.budgetRepo.GetBudgetsByUser(ctx, userID)
	accounts := s.accountRepo.GetAccountsByUser(ctx, userID)

	for i := range budgets {
		var actual float64
		for _, account := range accounts {
			ops := s.operationRepo.GetOperationsByAccount(ctx, account.ID)
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
