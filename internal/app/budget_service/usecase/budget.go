package budget

import (
	"context"

	pkgerrors "github.com/pkg/errors"

	bdgmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/models"
	bdgpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

type UseCase struct {
	budgetSvc BudgetService
	clock     clock.Clock
}

func NewBudgetUseCase(budgetService BudgetService) *UseCase {
	realClock := clock.RealClock{}
	return &UseCase{
		budgetSvc: budgetService,
		clock:     realClock,
	}
}

func (uc *UseCase) GetBudgets(ctx context.Context, userID int) (*bdgpb.ListBudgetsResponse, error) {
	log := logger.FromContext(ctx)
	budgetsData, err := uc.budgetSvc.GetBudgets(ctx, userID)
	if err != nil {
		log.Error("Failed to get budgets for user", "error", err, "user_id", userID)
		return nil, pkgerrors.Wrap(err, "budget.GetBudgetsForUser")
	}
	return budgetsData, nil
}

func (uc *UseCase) GetBudget(ctx context.Context, budgetID, userID int) (*bdgpb.Budget, error) {
	log := logger.FromContext(ctx)
	budgetsData, err := uc.budgetSvc.GetBudgetByID(ctx, budgetID, userID)
	if err != nil {
		log.Error("Failed to get budgets for user", "error", err, "user_id", userID)
		return nil, pkgerrors.Wrap(err, "budget.GetBudgetByID")
	}
	return budgetsData, nil
}

func (uc *UseCase) CreateBudget(ctx context.Context, req bdgmodels.CreateBudgetRequest, userID int) (*bdgpb.Budget, error) {
	log := logger.FromContext(ctx)
	budgetData, err := uc.budgetSvc.CreateBudget(ctx, req, userID)
	if err != nil {
		log.Error("Failed to create budget for user", "error", err, "user_id", userID)
		return nil, pkgerrors.Wrap(err, "budget.CreateBudget")
	}
	return budgetData, nil
}

func (uc *UseCase) UpdateBudget(ctx context.Context, req bdgmodels.UpdatedBudgetRequest) (*bdgpb.Budget, error) {
	log := logger.FromContext(ctx)
	budgetData, err := uc.budgetSvc.UpdateBudget(ctx, req)
	if err != nil {
		log.Error("Failed to update budget for user", "error", err, "user_id", req.UserID)
		return nil, pkgerrors.Wrap(err, "budget.UpdateBudget")
	}
	return budgetData, nil
}

func (uc *UseCase) DeleteBudget(ctx context.Context, budgetID, userID int) (*bdgpb.Budget, error) {
	log := logger.FromContext(ctx)
	budgetData, err := uc.budgetSvc.DeleteBudget(ctx, budgetID, userID)
	if err != nil {
		log.Error("Failed to delete budget for user", "error", err, "user_id", userID)
		return nil, pkgerrors.Wrap(err, "budget.DeleteBudget")
	}
	return budgetData, nil
}
