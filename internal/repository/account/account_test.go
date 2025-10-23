package account

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAccountsByUser_Table(t *testing.T) {
	repo := NewRepository(
		[]AccountDB{{ID: 1, Balance: 10, Type: "card"}, {ID: 2, Balance: 5, Type: "cash"}},
		[]UserAccountDB{{ID: 1, UserID: 1, AccountID: 1}},
	)

	tests := []struct {
		name   string
		userID int
		expect int
	}{
		{"has-one", 1, 1},
		{"has-none", 2, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := repo.GetAccountsByUser(context.Background(), tt.userID)
			assert.Len(t, got, tt.expect)
			if tt.expect == 1 {
				assert.Equal(t, 1, got[0].ID)
			}
		})
	}
}
