package budget

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

func BudgetToAPI(budget models.Budget) BudgetAPI {
	return BudgetAPI{
		ID:          budget.ID,
		UserID:      budget.UserID,
		Amount:      budget.Amount,
		Actual:      budget.Actual,
		CurrencyID:  budget.CurrencyID,
		Description: budget.Description,
		PeriodStart: budget.PeriodStart,
		PeriodEnd:   budget.PeriodEnd,
	}
}

func BudgetsToAPI(userID int, budgets []models.Budget) BudgetsAPI {
	budgetDTOs := make([]BudgetAPI, 0, len(budgets))
	for _, budget := range budgets {
		budgetDTOs = append(budgetDTOs, BudgetToAPI(budget))
	}

	return BudgetsAPI{
		UserID:  userID,
		Budgets: budgetDTOs,
	}
}
