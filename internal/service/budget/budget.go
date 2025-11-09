package budget

import (
	"context"
	"fmt"

	pkgerrors "github.com/pkg/errors"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

type Service struct {
	repo interface {
		GetBudgetsByUser(ctx context.Context, userID int) ([]models.Budget, error)
		GetAccountsByUser(ctx context.Context, userID int) ([]models.Account, error)
		GetOperationsByAccount(ctx context.Context, accountID int) ([]models.Operation, error)
	}
	clock clock.Clock
}

func NewService(repo interface {
	GetBudgetsByUser(ctx context.Context, userID int) ([]models.Budget, error)
	GetAccountsByUser(ctx context.Context, userID int) ([]models.Account, error)
	GetOperationsByAccount(ctx context.Context, accountID int) ([]models.Operation, error)
}, clck clock.Clock) *Service {
	return &Service{
		repo:  repo,
		clock: clck,
	}
}

func (s *Service) GetBudgetsForUser(ctx context.Context, userID int) ([]models.Budget, error) {
	budgets, err := s.repo.GetBudgetsByUser(ctx, userID)
	if err != nil {
		return []models.Budget{}, pkgerrors.Wrap(err, "Failed to get budgets for user")
	}

	if budgets == nil {
		budgets = []models.Budget{}
	}
	accounts, err := s.repo.GetAccountsByUser(ctx, userID)
	if err != nil {
		return []models.Budget{}, pkgerrors.Wrap(err, "Failed to get accounts for user")
	}

	if accounts == nil {
		accounts = []models.Account{}
	}

	for i := range budgets {
		var actual float64
		for _, account := range accounts {
			ops, err := s.repo.GetOperationsByAccount(ctx, account.ID)
			if err != nil {
				return []models.Budget{}, pkgerrors.Wrap(err, "Failed to get budgets for user")
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
	budgets, err := s.repo.GetBudgetsByUser(ctx, userID)
	if err != nil {
		return models.Budget{}, pkgerrors.Wrap(err, "budget.GetBudgetByID: failed to get budgets")
	}

	for _, budget := range budgets {
		if budget.ID == budgetID {
			return budget, nil
		}
	}

	return models.Budget{}, fmt.Errorf("budget not found: %s", models.ErrCodeBudgetNotFound)
}
