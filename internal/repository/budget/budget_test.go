package budget

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

func TestGetBudgetsByUser_Table(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	repo := NewRepository([]BudgetDB{{ID: 1, UserID: 1, Amount: 100}, {ID: 2, UserID: 2, Amount: 200}}, fixedClock)

	tests := []struct {
		name   string
		userID int
		expect int
	}{
		{"has-one", 1, 1},
		{"has-none", 3, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := repo.GetBudgetsByUser(context.Background(), tt.userID)
			assert.Len(t, got, tt.expect)
			if tt.expect == 1 {
				assert.Equal(t, 1, got[0].ID)
			}
		})
	}
}

func TestGetBudgetsByUser_Multiple(t *testing.T) {
	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	repo := NewRepository([]BudgetDB{
		{ID: 1, UserID: 1, Amount: 100},
		{ID: 2, UserID: 1, Amount: 200},
		{ID: 3, UserID: 2, Amount: 300},
	}, fixedClock)

	result, err := repo.GetBudgetsByUser(context.Background(), 1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 100.0, result[0].Amount)
	assert.Equal(t, 200.0, result[1].Amount)
}

func TestGetBudgetsByUser_Empty(t *testing.T) {
	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	repo := NewRepository([]BudgetDB{}, fixedClock)

	result, err := repo.GetBudgetsByUser(context.Background(), 1)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestGetBudgetsByUser_AllFields(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	now := fixedClock.FixedTime
	repo := NewRepository([]BudgetDB{
		{
			ID:          1,
			UserID:      1,
			Amount:      5000.0,
			CurrencyID:  1,
			Description: "Monthly budget",
			CreatedAt:   now,
			PeriodStart: now.AddDate(0, 0, -15),
			PeriodEnd:   now.AddDate(0, 0, 15),
		},
	}, fixedClock)

	result, err := repo.GetBudgetsByUser(context.Background(), 1)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 1, result[0].ID)
	assert.Equal(t, 1, result[0].UserID)
	assert.Equal(t, 5000.0, result[0].Amount)
	assert.Equal(t, 1, result[0].CurrencyID)
	assert.Equal(t, "Monthly budget", result[0].Description)
}

func TestGetBudgetsByUser_EmptyResult(t *testing.T) {
	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	repo := NewRepository([]BudgetDB{
		{ID: 1, UserID: 1, Amount: 100},
		{ID: 2, UserID: 2, Amount: 200},
	}, fixedClock)

	result, err := repo.GetBudgetsByUser(context.Background(), 99)
	assert.NoError(t, err)
	assert.Empty(t, result)
}
