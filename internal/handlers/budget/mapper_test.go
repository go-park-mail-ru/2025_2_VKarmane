package budget

import (
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestBudgetToAPI(t *testing.T) {
	now := time.Now()
	b := models.Budget{ID: 5, UserID: 2, Amount: 100, Actual: 30, CurrencyID: 1, Description: "d", PeriodStart: now, PeriodEnd: now}
	api := BudgetToAPI(b)
	assert.Equal(t, 5, api.ID)
	assert.Equal(t, 2, api.UserID)
	assert.Equal(t, 100.0, api.Amount)
	assert.Equal(t, 30.0, api.Actual)
	assert.Equal(t, 1, api.CurrencyID)
	assert.Equal(t, "d", api.Description)
}

func TestBudgetsToAPI(t *testing.T) {
	bs := []models.Budget{{ID: 1}, {ID: 2}}
	out := BudgetsToAPI(9, bs)
	assert.Equal(t, 9, out.UserID)
	assert.Len(t, out.Budgets, 2)
	assert.Equal(t, 1, out.Budgets[0].ID)
}
