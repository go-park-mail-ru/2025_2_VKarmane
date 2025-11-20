package budget

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	bdgerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/errors"
	bdgmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestPostgresRepository_GetBudgetsByUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	ctx := context.Background()
	userID := 1

	rows := sqlmock.NewRows([]string{
		"_id", "user_id", "category_id", "currency_id", "amount", "budget_description",
		"created_at", "updated_at", "closed_at", "period_start", "period_end",
	}).AddRow(1, userID, 2, 1, 100.0, "desc", time.Now(), time.Now(), nil, time.Now(), time.Now())

	mock.ExpectQuery("SELECT _id, user_id, category_id, currency_id, amount, budget_description,").
		WithArgs(userID).
		WillReturnRows(rows)

	budgets, err := repo.GetBudgetsByUser(ctx, userID)
	require.NoError(t, err)
	require.Len(t, budgets, 1)
	require.Equal(t, 1, budgets[0].ID)
}

func TestPostgresRepository_CreateBudget(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)
	ctx := context.Background()

	budget := bdgmodels.Budget{
		UserID:      1,
		CategoryID:  2,
		CurrencyID:  1,
		Amount:      100.0,
		Description: "desc",
		PeriodStart: time.Now(),
		PeriodEnd:   time.Now().Add(24 * time.Hour),
	}

	mock.ExpectQuery("INSERT INTO budget").
		WithArgs(budget.UserID, budget.CategoryID, budget.CurrencyID, budget.Amount, budget.Description, budget.PeriodStart, budget.PeriodEnd).
		WillReturnRows(sqlmock.NewRows([]string{"_id", "created_at", "updated_at"}).AddRow(1, time.Now(), time.Now()))

	created, err := repo.CreateBudget(ctx, budget)
	require.NoError(t, err)
	require.Equal(t, 1, created.ID)
}

func TestPostgresRepository_UpdateBudget(t *testing.T) {
	fixedTime := time.Date(2030, 5, 20, 10, 0, 0, 0, time.UTC)
	fixed := clock.FixedClock{FixedTime: fixedTime}
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)
	ctx := context.Background()

	sum := 200.0
	descr := "updated"
	req := bdgmodels.UpdatedBudgetRequest{
		BudgetID:    1,
		UserID:      1,
		Amount:      &sum,
		Description: &descr,
		PeriodStart: &fixed.FixedTime,
		PeriodEnd:   &fixed.FixedTime,
	}

	amountArg := interface{}(nil)
	if req.Amount != nil {
		amountArg = *req.Amount
	}
	descArg := interface{}(nil)
	if req.Description != nil {
		descArg = *req.Description
	}

	mock.ExpectQuery("UPDATE budget").
		WithArgs(amountArg, descArg, req.PeriodStart, req.PeriodEnd, req.BudgetID, req.UserID).
		WillReturnRows(sqlmock.NewRows([]string{
			"_id", "user_id", "category_id", "currency_id", "amount", "budget_description",
			"created_at", "updated_at", "period_start", "period_end",
		}).AddRow(
			req.BudgetID,
			req.UserID,
			2,
			1,
			*req.Amount,
			*req.Description,
			fixed.FixedTime,
			fixed.FixedTime,
			*req.PeriodStart,
			*req.PeriodEnd,
		),
		)

	updated, err := repo.UpdateBudget(ctx, req)
	require.NoError(t, err)
	require.Equal(t, *req.Amount, updated.Amount)
	require.Equal(t, *req.Description, updated.Description)
}

func TestPostgresRepository_DeleteBudget(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)
	ctx := context.Background()
	budgetID := 1

	mock.ExpectQuery("UPDATE budget").
		WithArgs(budgetID).
		WillReturnRows(sqlmock.NewRows([]string{
			"_id", "user_id", "category_id", "currency_id", "amount", "budget_description",
			"created_at", "updated_at", "period_start", "period_end",
		}).AddRow(budgetID, 1, 2, 1, 100.0, "desc", time.Now(), time.Now(), time.Now(), time.Now()))

	deleted, err := repo.DeleteBudget(ctx, budgetID)
	require.NoError(t, err)
	require.Equal(t, budgetID, deleted.ID)
}

func TestPostgresRepository_GetBudgetsByUser_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)
	ctx := context.Background()
	userID := 1

	mock.ExpectQuery("SELECT _id, user_id, category_id, currency_id, amount, budget_description,").
		WithArgs(userID).
		WillReturnError(errors.New("db error"))

	_, err = repo.GetBudgetsByUser(ctx, userID)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get budgets by user")
}

func TestPostgresRepository_CreateBudget_UniqueViolation(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)
	ctx := context.Background()

	budget := bdgmodels.Budget{
		UserID:      1,
		CategoryID:  2,
		CurrencyID:  1,
		Amount:      100,
		Description: "desc",
		PeriodStart: time.Now(),
		PeriodEnd:   time.Now().Add(24 * time.Hour),
	}

	mockErr := &pq.Error{Code: UniqueViolation}
	mock.ExpectQuery("INSERT INTO budget").WithArgs(
		budget.UserID, budget.CategoryID, budget.CurrencyID, budget.Amount,
		budget.Description, budget.PeriodStart, budget.PeriodEnd,
	).WillReturnError(mockErr)

	_, err = repo.CreateBudget(ctx, budget)
	require.ErrorIs(t, err, bdgerrors.ErrBudgetExists)
}

func TestPostgresRepository_CreateBudget_NotNullViolation(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)
	ctx := context.Background()

	budget := bdgmodels.Budget{
		UserID:      1,
		CategoryID:  2,
		CurrencyID:  1,
		Amount:      100,
		Description: "desc",
		PeriodStart: time.Now(),
		PeriodEnd:   time.Now().Add(24 * time.Hour),
	}

	mockErr := &pq.Error{Code: NotNullViolation}
	mock.ExpectQuery("INSERT INTO budget").WithArgs(
		budget.UserID, budget.CategoryID, budget.CurrencyID, budget.Amount,
		budget.Description, budget.PeriodStart, budget.PeriodEnd,
	).WillReturnError(mockErr)

	_, err = repo.CreateBudget(ctx, budget)
	require.ErrorIs(t, err, bdgerrors.ErrInavlidData)
}

func TestPostgresRepository_UpdateBudget_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)
	ctx := context.Background()

	req := bdgmodels.UpdatedBudgetRequest{
		BudgetID: 1,
		UserID:   1,
	}

	mock.ExpectQuery("UPDATE budget").
		WithArgs(req.Amount, req.Description, req.PeriodStart, req.PeriodEnd, req.BudgetID, req.UserID).
		WillReturnError(sql.ErrNoRows)

	_, err = repo.UpdateBudget(ctx, req)
	require.ErrorIs(t, err, bdgerrors.ErrBudgetNotFound)
}

func TestPostgresRepository_DeleteBudget_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)
	ctx := context.Background()
	budgetID := 1

	mock.ExpectQuery("UPDATE budget").
		WithArgs(budgetID).
		WillReturnError(sql.ErrNoRows)

	_, err = repo.DeleteBudget(ctx, budgetID)
	require.ErrorIs(t, err, bdgerrors.ErrBudgetNotFound)
}
