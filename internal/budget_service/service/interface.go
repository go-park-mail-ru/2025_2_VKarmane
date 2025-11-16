package budget

import (
	"context"

	bdgmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/models"
)

type BudgetRepository interface {
	GetBudgetsByUser(ctx context.Context, userID int) ([]bdgmodels.Budget, error)
	// GetAccountsByUser(ctx context.Context, userID int) ([]models.Account, error)
	// GetOperationsByAccount(ctx context.Context, accountID int) ([]models.OperationInList, error) //пока нет микросервисов
	CreateBudget(ctx context.Context, budget bdgmodels.Budget) (bdgmodels.Budget, error)
	UpdateBudget(ctx context.Context, req bdgmodels.UpdatedBudgetRequest) (bdgmodels.Budget, error)
	DeleteBudget(ctx context.Context, budgetID int) (bdgmodels.Budget, error)
}