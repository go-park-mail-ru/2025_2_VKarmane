package budget

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/budget"
)

type UseCase struct {
	budgetSvc BudgetService
}

func NewUseCase(store repository.Repository) *UseCase {
	budgetRepoAdapter := budget.NewPostgresBudgetRepositoryAdapter(store)
	accountRepoAdapter := budget.NewPostgresAccountRepositoryAdapter(store)
	operationRepoAdapter := budget.NewPostgresOperationRepositoryAdapter(store)

	budgetService := budget.NewService(budgetRepoAdapter, accountRepoAdapter, operationRepoAdapter)

	return &UseCase{
		budgetSvc: budgetService,
	}
}

func (uc *UseCase) GetBudgetsForUser(ctx context.Context, userID int) ([]models.Budget, error) {
	budgetsData, err := uc.budgetSvc.GetBudgetsForUser(ctx, userID)
	if err != nil {
		if log := logger.FromContext(ctx); log != nil {
			log.Error("Failed to get budgets for user", "error", err, "user_id", userID)
		}

		return nil, fmt.Errorf("budget.GetBudgetsForUser: %w", err)
	}

	return budgetsData, nil
}

func (uc *UseCase) GetBudgetByID(ctx context.Context, userID, budgetID int) (models.Budget, error) {
	budgetsData, err := uc.budgetSvc.GetBudgetsForUser(ctx, userID)
	if err != nil {
		if log := logger.FromContext(ctx); log != nil {
			log.Error("Failed to get budgets for user", "error", err, "user_id", userID)
		}

		return models.Budget{}, fmt.Errorf("budget.GetBudgetByID: %w", err)
	}

	for _, budget := range budgetsData {
		if budget.ID == budgetID {
			return budget, nil
		}
	}

	if log := logger.FromContext(ctx); log != nil {
		log.Warn("Budget not found", "user_id", userID, "budget_id", budgetID)
	}

	return models.Budget{}, fmt.Errorf("budget.GetBudgetByID: %s", models.ErrCodeBudgetNotFound)
}
