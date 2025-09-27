package budget

import "github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"

type BudgetService interface {
	GetBudgetsForUser(userID int) ([]models.Budget, error)
}
