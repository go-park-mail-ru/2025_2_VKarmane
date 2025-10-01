package budget

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBudgetsByUser_Table(t *testing.T) {
	repo := NewRepository([]BudgetDB{{ID: 1, UserID: 1, Amount: 100}, {ID: 2, UserID: 2, Amount: 200}})

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
			got := repo.GetBudgetsByUser(context.Background(), tt.userID)
			assert.Len(t, got, tt.expect)
			if tt.expect == 1 {
				assert.Equal(t, 1, got[0].ID)
			}
		})
	}
}
