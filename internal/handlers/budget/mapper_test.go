package budget

import (
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestBudgetToAPI(t *testing.T) {
	now := time.Now()
	budget := models.Budget{
		ID:          1,
		UserID:      10,
		Amount:      5000.00,
		Actual:      2500.50,
		CurrencyID:  1,
		Description: "Monthly budget",
		CreatedAt:   now,
		PeriodStart: now.AddDate(0, 0, -15),
		PeriodEnd:   now.AddDate(0, 0, 15),
	}

	apiBudget := BudgetToAPI(budget)

	assert.Equal(t, 1, apiBudget.ID)
	assert.Equal(t, 10, apiBudget.UserID)
	assert.Equal(t, 5000.00, apiBudget.Amount)
	assert.Equal(t, 2500.50, apiBudget.Actual)
	assert.Equal(t, 1, apiBudget.CurrencyID)
	assert.Equal(t, "Monthly budget", apiBudget.Description)
}

func TestBudgetsToAPI(t *testing.T) {
	now := time.Now()
	budgets := []models.Budget{
		{
			ID:          1,
			Amount:      1000.00,
			Actual:      500.00,
			CurrencyID:  1,
			Description: "Budget 1",
			CreatedAt:   now,
		},
		{
			ID:          2,
			Amount:      2000.00,
			Actual:      1000.00,
			CurrencyID:  1,
			Description: "Budget 2",
			CreatedAt:   now,
		},
	}

	apiBudgets := BudgetsToAPI(1, budgets)

	assert.Equal(t, 1, apiBudgets.UserID)
	assert.Len(t, apiBudgets.Budgets, 2)
	assert.Equal(t, 1, apiBudgets.Budgets[0].ID)
	assert.Equal(t, 2, apiBudgets.Budgets[1].ID)
	assert.Equal(t, 1000.00, apiBudgets.Budgets[0].Amount)
	assert.Equal(t, 2000.00, apiBudgets.Budgets[1].Amount)
}

func TestBudgetsToAPI_Empty(t *testing.T) {
	apiBudgets := BudgetsToAPI(1, []models.Budget{})
	assert.Empty(t, apiBudgets.Budgets)
	assert.Equal(t, 1, apiBudgets.UserID)
}
