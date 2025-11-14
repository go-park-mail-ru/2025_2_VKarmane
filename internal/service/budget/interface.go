package budget

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type BudgetService interface {
	GetBudgetsForUser(ctx context.Context, userID int) ([]models.Budget, error)
	GetBudgetByID(ctx context.Context, userID, budgetID int) (models.Budget, error)
}

type BudgetRepository interface {
	GetBudgetsByUser(ctx context.Context, userID int) ([]models.Budget, error)
}

type AccountRepository interface {
	GetAccountsByUser(ctx context.Context, userID int) ([]models.Account, error)
}

type OperationRepository interface {
	GetOperationsByAccount(ctx context.Context, accountID int) ([]models.OperationInList, error)
}
