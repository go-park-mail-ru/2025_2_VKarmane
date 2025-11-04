package account

import (
	"context"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	"github.com/stretchr/testify/assert"
)

func TestGetAccountsByUser_Table(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	repo := NewRepository(
		[]AccountDB{{ID: 1, Balance: 10, Type: "card"}, {ID: 2, Balance: 5, Type: "cash"}},
		[]UserAccountDB{{ID: 1, UserID: 1, AccountID: 1}},
		fixedClock,
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

func TestGetAccountsByUser_Multiple(t *testing.T) {
	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	repo := NewRepository(
		[]AccountDB{
			{ID: 1, Balance: 100},
			{ID: 2, Balance: 200},
			{ID: 3, Balance: 300},
		},
		[]UserAccountDB{
			{ID: 1, UserID: 1, AccountID: 1},
			{ID: 2, UserID: 1, AccountID: 2},
			{ID: 3, UserID: 2, AccountID: 3},
		},
		fixedClock,
	)

	result, err := repo.GetAccountsByUser(context.Background(), 1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestGetAccountsByUser_Empty(t *testing.T) {
	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	repo := NewRepository(
		[]AccountDB{},
		[]UserAccountDB{},
		fixedClock,
	)

	result, err := repo.GetAccountsByUser(context.Background(), 1)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestGetAccountsByUser_NoAccountsForUser(t *testing.T) {
	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	repo := NewRepository(
		[]AccountDB{
			{ID: 1, Balance: 100, Type: "card"},
			{ID: 2, Balance: 200, Type: "cash"},
		},
		[]UserAccountDB{
			{ID: 1, UserID: 1, AccountID: 1},
		},
		fixedClock,
	)

	result, err := repo.GetAccountsByUser(context.Background(), 2)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestGetAccountsByUser_AllFields(t *testing.T) {
	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	repo := NewRepository(
		[]AccountDB{
			{ID: 1, Balance: 150.50, Type: "card", CurrencyID: 1, CreatedAt: time.Now()},
		},
		[]UserAccountDB{
			{ID: 1, UserID: 1, AccountID: 1},
		},
		fixedClock,
	)

	result, err := repo.GetAccountsByUser(context.Background(), 1)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 1, result[0].ID)
	assert.Equal(t, 150.50, result[0].Balance)
	assert.Equal(t, "card", result[0].Type)
	assert.Equal(t, 1, result[0].CurrencyID)
}
