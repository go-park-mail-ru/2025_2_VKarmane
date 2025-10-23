package budget

import (
	"context"
	"fmt"

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
	budgets, err := s.budgetRepo.GetBudgetsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	accounts, err := s.accountRepo.GetAccountsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	for i := range budgets {
		var actual float64
		for _, account := range accounts {
			ops, err := s.operationRepo.GetOperationsByAccount(ctx, account.ID)
			if err != nil {
				return nil, err
			}

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

func (s *Service) GetBudgetByID(ctx context.Context, userID, budgetID int) (models.Budget, error) {
	budgets, err := s.budgetRepo.GetBudgetsByUser(ctx, userID)
	if err != nil {
		return models.Budget{}, err
	}

	for _, budget := range budgets {
		if budget.ID == budgetID {
			return budget, nil
		}
	}

	return models.Budget{}, fmt.Errorf("budget not found: %s", models.ErrCodeBudgetNotFound)
}
