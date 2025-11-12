package budget

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

func TestPostgresRepository_GetBudgetsByUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)
	userID := 1
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"_id", "user_id", "category_id", "currency_id", "amount", "budget_description",
		"created_at", "updated_at", "closed_at", "period_start", "period_end",
	}).AddRow(1, userID, 0, 1, 5000.0, "Monthly budget", now, now, nil, now.AddDate(0, 0, -15), now.AddDate(0, 0, 15)).
		AddRow(2, userID, 0, 1, 2000.0, "Food budget", now, now, nil, now.AddDate(0, 0, -10), now.AddDate(0, 0, 20))

	mock.ExpectQuery(`SELECT.*FROM budget`).WithArgs(userID).WillReturnRows(rows)

	budgets, err := repo.GetBudgetsByUser(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, budgets, 2)
	assert.Equal(t, 1, budgets[0].ID)
	assert.Equal(t, 5000.0, budgets[0].Amount)
	assert.Equal(t, "Monthly budget", budgets[0].Description)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetBudgetsByUser_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)
	userID := 1

	rows := sqlmock.NewRows([]string{
		"_id", "user_id", "category_id", "currency_id", "amount", "budget_description",
		"created_at", "updated_at", "closed_at", "period_start", "period_end",
	})
	mock.ExpectQuery(`SELECT.*FROM budget`).WithArgs(userID).WillReturnRows(rows)

	budgets, err := repo.GetBudgetsByUser(context.Background(), userID)
	assert.NoError(t, err)
	assert.Empty(t, budgets)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetBudgetsByUser_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)
	userID := 1

	mock.ExpectQuery(`SELECT.*FROM budget`).WithArgs(userID).WillReturnError(sql.ErrConnDone)

	budgets, err := repo.GetBudgetsByUser(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, budgets)
	assert.Contains(t, err.Error(), "failed to get budgets by user")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_CreateBudget(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)
	now := time.Now()
	budget := models.Budget{
		UserID:      1,
		CategoryID:  0,
		CurrencyID:  1,
		Amount:      5000,
		Description: "Monthly budget",
		PeriodStart: now,
		PeriodEnd:   now.AddDate(0, 0, 30),
	}

	mock.ExpectQuery(`INSERT INTO budget`).
		WithArgs(budget.UserID, budget.CategoryID, budget.CurrencyID, budget.Amount,
			budget.Description, budget.PeriodStart, budget.PeriodEnd).
		WillReturnRows(sqlmock.NewRows([]string{"_id", "created_at", "updated_at"}).AddRow(1, now, now))

	created, err := repo.CreateBudget(context.Background(), budget)
	assert.NoError(t, err)
	assert.Equal(t, 1, created.ID)
	assert.NoError(t, mock.ExpectationsWereMet())

	pqErr := &pq.Error{Code: "23505"} // UniqueViolation
	mock.ExpectQuery(`INSERT INTO budget`).WillReturnError(pqErr)
	_, err = repo.CreateBudget(context.Background(), budget)
	assert.True(t, errors.Is(err, ErrActiveBudgetExists))

	pqErr = &pq.Error{Code: "23502"} // NotNullViolation
	mock.ExpectQuery(`INSERT INTO budget`).WillReturnError(pqErr)
	_, err = repo.CreateBudget(context.Background(), budget)
	assert.True(t, errors.Is(err, ErrNotNullViolation))

	pqErr = &pq.Error{Code: "23514"} // CheckViolation
	mock.ExpectQuery(`INSERT INTO budget`).WillReturnError(pqErr)
	_, err = repo.CreateBudget(context.Background(), budget)
	assert.True(t, errors.Is(err, ErrCheckViolation))
}

func TestPostgresRepository_UpdateBudget(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)
	userID := 1
	budgetID := 1
	now := time.Now()
	req := models.UpdatedBudgetRequest{
		Amount:      float64Ptr(6000),
		Description: strPtr("Updated"),
	}

	rows := sqlmock.NewRows([]string{
		"_id", "user_id", "category_id", "currency_id", "amount", "budget_description",
		"created_at", "updated_at", "period_start", "period_end",
	}).AddRow(budgetID, userID, 0, 1, 6000, "Updated", now, now, now, now.AddDate(0, 0, 30))

	mock.ExpectQuery(`UPDATE budget`).
		WithArgs(req.Amount, req.Description, nil, nil, budgetID, userID).
		WillReturnRows(rows)

	updated, err := repo.UpdateBudget(context.Background(), req, userID, budgetID)
	assert.NoError(t, err)
	assert.Equal(t, 6000.0, updated.Amount)
	assert.Equal(t, "Updated", updated.Description)
	assert.NoError(t, mock.ExpectationsWereMet())

	mock.ExpectQuery(`UPDATE budget`).WithArgs(req.Amount, req.Description, nil, nil, 999, userID).WillReturnError(sql.ErrNoRows)
	_, err = repo.UpdateBudget(context.Background(), req, userID, 999)
	assert.True(t, errors.Is(err, ErrBudgetNotFound))
}

func TestPostgresRepository_DeleteBudget(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)
	budgetID := 1
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"_id", "user_id", "category_id", "currency_id", "amount", "budget_description",
		"created_at", "updated_at", "period_start", "period_end",
	}).AddRow(budgetID, 1, 0, 1, 5000, "Monthly budget", now, now, now, now.AddDate(0, 0, 30))

	mock.ExpectQuery(`UPDATE budget`).
		WithArgs(budgetID).
		WillReturnRows(rows)

	deleted, err := repo.DeleteBudget(context.Background(), budgetID)
	assert.NoError(t, err)
	assert.Equal(t, 5000.0, deleted.Amount)
	assert.NoError(t, mock.ExpectationsWereMet())

	mock.ExpectQuery(`UPDATE budget`).WithArgs(999).WillReturnError(sql.ErrNoRows)
	_, err = repo.DeleteBudget(context.Background(), 999)
	assert.True(t, errors.Is(err, ErrBudgetNotFound))
}

func float64Ptr(f float64) *float64 { return &f }
func strPtr(s string) *string       { return &s }
