package budget

import "github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"

type BudgetUseCase interface {
	GetBudgetsForUser(userID int) ([]models.Budget, error)
	GetBudgetByID(userID, budgetID int) (models.Budget, error)
}
