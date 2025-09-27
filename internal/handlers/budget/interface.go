package budget

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type BudgetUseCase interface {
	GetBudgetsForUser(ctx context.Context, userID int) ([]models.Budget, error)
	GetBudgetByID(ctx context.Context, userID, budgetID int) (models.Budget, error)
}
