package budget

import (
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestBudgetDBToModel(t *testing.T) {
	now := time.Now()
	budgetDB := BudgetDB{
		ID:          1,
		UserID:      10,
		Amount:      5000.00,
		CurrencyID:  1,
		Description: "Monthly budget",
		CreatedAt:   now,
		PeriodStart: now.AddDate(0, 0, -15),
		PeriodEnd:   now.AddDate(0, 0, 15),
	}

	budget := BudgetDBToModel(budgetDB)

	assert.Equal(t, 1, budget.ID)
	assert.Equal(t, 10, budget.UserID)
	assert.Equal(t, 5000.00, budget.Amount)
	assert.Equal(t, 1, budget.CurrencyID)
	assert.Equal(t, "Monthly budget", budget.Description)
	assert.Equal(t, now, budget.CreatedAt)
}

func TestBudgetModelToDB(t *testing.T) {
	now := time.Now()
	budget := models.Budget{
		ID:          2,
		UserID:      20,
		Amount:      2000.00,
		CurrencyID:  2,
		Description: "Weekly budget",
		CreatedAt:   now,
		PeriodStart: now.AddDate(0, 0, -7),
		PeriodEnd:   now.AddDate(0, 0, 7),
	}

	budgetDB := BudgetModelToDB(budget)

	assert.Equal(t, 2, budgetDB.ID)
	assert.Equal(t, 20, budgetDB.UserID)
	assert.Equal(t, 2000.00, budgetDB.Amount)
	assert.Equal(t, 2, budgetDB.CurrencyID)
	assert.Equal(t, "Weekly budget", budgetDB.Description)
	assert.Equal(t, now, budgetDB.CreatedAt)
}
