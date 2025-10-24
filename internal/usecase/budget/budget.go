package budget

import (
	"context"
	"fmt"

	pkgerrors "github.com/pkg/errors"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

type UseCase struct {
	budgetSvc BudgetService
	clock     clock.Clock
}

func NewUseCase(budgetService BudgetService) *UseCase {
	realClock := clock.RealClock{}
	return &UseCase{
		budgetSvc: budgetService,
		clock:     realClock,
	}
}

func (uc *UseCase) GetBudgetsForUser(ctx context.Context, userID int) ([]models.Budget, error) {
	log := logger.FromContext(ctx)
	budgetsData, err := uc.budgetSvc.GetBudgetsForUser(ctx, userID)
	if err != nil {
		if log != nil {
			log.Error("Failed to get budgets for user", "error", err, "user_id", userID)
		}

		return nil, pkgerrors.Wrap(err, "budget.GetBudgetsForUser")
	}

	return budgetsData, nil
}

func (uc *UseCase) GetBudgetByID(ctx context.Context, userID, budgetID int) (models.Budget, error) {
	log := logger.FromContext(ctx)
	budgetsData, err := uc.budgetSvc.GetBudgetsForUser(ctx, userID)
	if err != nil {
		if log != nil {
			log.Error("Failed to get budgets for user", "error", err, "user_id", userID)
		}

		return models.Budget{}, pkgerrors.Wrap(err, "budget.GetBudgetByID")
	}

	for _, budget := range budgetsData {
		if budget.ID == budgetID {
			return budget, nil
		}
	}

	if log != nil {
		log.Warn("Budget not found", "user_id", userID, "budget_id", budgetID)
	}

	return models.Budget{}, fmt.Errorf("budget.GetBudgetByID: %s", models.ErrCodeBudgetNotFound)
}
