package budget

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type BudgetUseCase interface {
	GetBudgetsForUser(ctx context.Context, userID int) ([]models.Budget, error)
	GetBudgetByID(ctx context.Context, userID, budgetID int) (models.Budget, error)
	CreateBudget(ctx context.Context, req models.CreateBudgetRequest, userID int) (models.Budget, error)
	UpdateBudget(ctx context.Context, req models.UpdatedBudgetRequest, userID, budgetID int) (models.Budget, error)
	DeleteBudget(ctx context.Context, userID, budgetID int) (models.Budget, error)
}
