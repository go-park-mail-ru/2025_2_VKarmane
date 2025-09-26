package budget

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

func BudgetDBToModel(budgetDB BudgetDB) models.Budget {
	return models.Budget{
		ID:          budgetDB.ID,
		UserID:      budgetDB.UserID,
		Amount:      budgetDB.Amount,
		CurrencyID:  budgetDB.CurrencyID,
		Description: budgetDB.Description,
		CreatedAt:   budgetDB.CreatedAt,
		UpdatedAt:   budgetDB.UpdatedAt,
		ClosedAt:    budgetDB.ClosedAt,
		PeriodStart: budgetDB.PeriodStart,
		PeriodEnd:   budgetDB.PeriodEnd,
	}
}

func BudgetModelToDB(budget models.Budget) BudgetDB {
	return BudgetDB{
		ID:          budget.ID,
		UserID:      budget.UserID,
		Amount:      budget.Amount,
		CurrencyID:  budget.CurrencyID,
		Description: budget.Description,
		CreatedAt:   budget.CreatedAt,
		UpdatedAt:   budget.UpdatedAt,
		ClosedAt:    budget.ClosedAt,
		PeriodStart: budget.PeriodStart,
		PeriodEnd:   budget.PeriodEnd,
	}
}
