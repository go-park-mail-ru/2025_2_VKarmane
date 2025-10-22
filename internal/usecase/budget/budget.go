package budget

import (
	"context"
	"fmt"

	pkgerrors "github.com/pkg/errors"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/account"
	budgetRepo "github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/budget"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/operation"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/budget"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

type UseCase struct {
	budgetSvc BudgetService
	clock     clock.Clock
}

func NewUseCase(store *repository.Store, clck clock.Clock) *UseCase {
	accountRepo := account.NewRepository(store.Accounts, store.UserAccounts, clck)
	budgetRepo := budgetRepo.NewRepository(store.Budget, clck)
	operationRepo := operation.NewRepository(store.Operations, clck)
	budgetService := budget.NewService(budgetRepo, accountRepo, operationRepo, clck)

	return &UseCase{
		budgetSvc: budgetService,
		clock:     clck,
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
