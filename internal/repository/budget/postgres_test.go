package budget

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresRepository_GetBudgetsByUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	repo := NewPostgresRepository(db)

	userID := 1
	now := time.Now()
	closedAt := sql.NullTime{Time: time.Time{}, Valid: false}
	rows := sqlmock.NewRows([]string{"_id", "user_id", "category_id", "currency_id", "amount", "budget_description", "created_at", "updated_at", "closed_at", "period_start", "period_end"}).
		AddRow(1, userID, 0, 1, 5000.0, "Monthly budget", now, now, closedAt, now.AddDate(0, 0, -15), now.AddDate(0, 0, 15)).
		AddRow(2, userID, 0, 1, 2000.0, "Food budget", now, now, closedAt, now.AddDate(0, 0, -10), now.AddDate(0, 0, 20))

	mock.ExpectQuery(`SELECT.*FROM budget`).
		WithArgs(userID).
		WillReturnRows(rows)

	budgets, err := repo.GetBudgetsByUser(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, budgets, 2)
	assert.Equal(t, 1, budgets[0].ID)
	assert.Equal(t, 5000.0, budgets[0].Amount)
	assert.Equal(t, "Monthly budget", budgets[0].Description)
	assert.Equal(t, 2, budgets[1].ID)
	assert.Equal(t, 2000.0, budgets[1].Amount)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetBudgetsByUser_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	repo := NewPostgresRepository(db)

	userID := 1
	rows := sqlmock.NewRows([]string{"_id", "user_id", "category_id", "currency_id", "amount", "budget_description", "created_at", "updated_at", "closed_at", "period_start", "period_end"})

	mock.ExpectQuery(`SELECT.*FROM budget`).
		WithArgs(userID).
		WillReturnRows(rows)

	budgets, err := repo.GetBudgetsByUser(context.Background(), userID)
	assert.NoError(t, err)
	assert.Empty(t, budgets)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetBudgetsByUser_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	repo := NewPostgresRepository(db)

	userID := 1
	mock.ExpectQuery(`SELECT.*FROM budget`).
		WithArgs(userID).
		WillReturnError(sql.ErrConnDone)

	budgets, err := repo.GetBudgetsByUser(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, budgets)
	assert.Contains(t, err.Error(), "failed to get budgets by user")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetBudgetsByUser_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	repo := NewPostgresRepository(db)

	userID := 1
	rows := sqlmock.NewRows([]string{"_id", "user_id", "category_id", "currency_id", "amount", "budget_description", "created_at", "updated_at", "closed_at", "period_start", "period_end"}).
		AddRow("invalid", "invalid", "invalid", "invalid", "invalid", "invalid", "invalid", "invalid", "invalid", "invalid", "invalid")

	mock.ExpectQuery(`SELECT.*FROM budget`).
		WithArgs(userID).
		WillReturnRows(rows)

	budgets, err := repo.GetBudgetsByUser(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, budgets)
	assert.Contains(t, err.Error(), "failed to scan budget")
	assert.NoError(t, mock.ExpectationsWereMet())
}
