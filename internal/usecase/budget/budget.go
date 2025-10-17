package budget

import (
	"context"
	"fmt"

	pkgErrors "github.com/pkg/errors"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/account"
	budgetRepo "github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/budget"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/operation"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/budget"
)

type UseCase struct {
	budgetSvc BudgetService
}

func NewUseCase(store *repository.Store) *UseCase {
	accountRepo := account.NewRepository(store.Accounts, store.UserAccounts)
	budgetRepo := budgetRepo.NewRepository(store.Budget)
	operationRepo := operation.NewRepository(store.Operations)
	budgetService := budget.NewService(budgetRepo, accountRepo, operationRepo)

	return &UseCase{
		budgetSvc: budgetService,
	}
}

func (uc *UseCase) GetBudgetsForUser(ctx context.Context, userID int) ([]models.Budget, error) {
	log := logger.FromContext(ctx)
	budgetsData, err := uc.budgetSvc.GetBudgetsForUser(ctx, userID)
	if err != nil {
		if log != nil {
			log.Error("Failed to get budgets for user", "error", err, "user_id", userID)
		}

		return nil, pkgErrors.Wrap(err, "budget.GetBudgetsForUser")
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

		return models.Budget{}, pkgErrors.Wrap(err, "budget.GetBudgetByID")
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
