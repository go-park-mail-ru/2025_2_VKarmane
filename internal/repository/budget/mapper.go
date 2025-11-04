package budget

import (
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

func BudgetDBToModel(budgetDB BudgetDB) models.Budget {
	var closedAt time.Time
	if budgetDB.ClosedAt != nil {
		closedAt = *budgetDB.ClosedAt
	}

	return models.Budget{
		ID:          budgetDB.ID,
		UserID:      budgetDB.UserID,
		CategoryID:  budgetDB.CategoryID,
		Amount:      budgetDB.Amount,
		CurrencyID:  budgetDB.CurrencyID,
		Description: budgetDB.Description,
		CreatedAt:   budgetDB.CreatedAt,
		UpdatedAt:   budgetDB.UpdatedAt,
		ClosedAt:    closedAt,
		PeriodStart: budgetDB.PeriodStart,
		PeriodEnd:   budgetDB.PeriodEnd,
	}
}

func BudgetModelToDB(budget models.Budget) BudgetDB {
	var closedAt *time.Time
	if !budget.ClosedAt.IsZero() {
		closedAt = &budget.ClosedAt
	}

	return BudgetDB{
		ID:          budget.ID,
		UserID:      budget.UserID,
		CategoryID:  budget.CategoryID,
		Amount:      budget.Amount,
		CurrencyID:  budget.CurrencyID,
		Description: budget.Description,
		CreatedAt:   budget.CreatedAt,
		UpdatedAt:   budget.UpdatedAt,
		ClosedAt:    closedAt,
		PeriodStart: budget.PeriodStart,
		PeriodEnd:   budget.PeriodEnd,
	}
}
