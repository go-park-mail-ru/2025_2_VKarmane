package budget

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type BudgetService interface {
	GetBudgetsForUser(ctx context.Context, userID int) ([]models.Budget, error)
}
