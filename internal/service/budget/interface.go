package budget

import "github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"

type BudgetService interface {
	GetBudgetsForUser(userID int) ([]models.Budget, error)
}

type BudgetRepository interface {
	GetBudgetsByUser(userID int) []models.Budget
}

type AccountRepository interface {
	GetAccountsByUser(userID int) []models.Account
}

type OperationRepository interface {
	GetOperationsByAccount(accountID int) []models.Operation
}
