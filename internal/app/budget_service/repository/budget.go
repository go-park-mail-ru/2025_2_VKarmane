package budget

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	bdgmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/models"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewDBConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Настройка connection pool
	maxOpenConns := getEnvAsInt("DB_MAX_OPEN_CONNS", 25)
	db.SetMaxOpenConns(maxOpenConns)

	maxIdleConns := getEnvAsInt("DB_MAX_IDLE_CONNS", 5)
	db.SetMaxIdleConns(maxIdleConns)

	connMaxLifetime := getEnvAsInt("DB_CONN_MAX_LIFETIME", 30) // минуты
	if connMaxLifetime > 0 {
		db.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Minute)
	}

	connMaxIdleTime := getEnvAsInt("DB_CONN_MAX_IDLE_TIME", 10) // минуты
	if connMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(time.Duration(connMaxIdleTime) * time.Minute)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}
	return defaultValue
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetBudgetsByUser(ctx context.Context, userID int) ([]bdgmodels.Budget, error) {
	query := `
		SELECT _id, user_id, category_id, currency_id, amount, budget_description, 
		       created_at, updated_at, closed_at, period_start, period_end
		FROM budget
		WHERE user_id = $1 AND closed_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get budgets by user: %w", err)
	}
	defer rows.Close()

	var budgets []bdgmodels.Budget

	for rows.Next() {
		var b BudgetDB
		if err := rows.Scan(
			&b.ID,
			&b.UserID,
			&b.CategoryID,
			&b.CurrencyID,
			&b.Amount,
			&b.Description,
			&b.CreatedAt,
			&b.UpdatedAt,
			&b.ClosedAt,
			&b.PeriodStart,
			&b.PeriodEnd,
		); err != nil {
			return nil, fmt.Errorf("failed to scan budget: %w", err)
		}

		pb := bdgmodels.Budget{
			ID:          b.ID,
			UserID:      b.UserID,
			CategoryID:  b.CategoryID,
			CurrencyID:  b.CurrencyID,
			Amount:      b.Amount,
			Description: b.Description,
			CreatedAt:   b.CreatedAt,
			UpdatedAt:   b.UpdatedAt,
			PeriodStart: b.PeriodStart,
			PeriodEnd:   b.PeriodEnd,
		}
		budgets = append(budgets, pb)
	}

	return budgets, nil
}

func (r *PostgresRepository) CreateBudget(ctx context.Context, budget bdgmodels.Budget) (bdgmodels.Budget, error) {
	query := `
		INSERT INTO budget (
			user_id, category_id, currency_id, amount, budget_description, 
			created_at, updated_at, period_start, period_end
		)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW(), $6, $7)
		RETURNING _id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		budget.UserID,
		budget.CategoryID,
		budget.CurrencyID,
		budget.Amount,
		budget.Description,
		budget.PeriodStart,
		budget.PeriodEnd,
	).Scan(&budget.ID, &budget.CreatedAt, &budget.UpdatedAt)

	if err != nil {
		return bdgmodels.Budget{}, MapPgError(err)
	}

	return budget, nil
}

func (r *PostgresRepository) UpdateBudget(ctx context.Context, req bdgmodels.UpdatedBudgetRequest) (bdgmodels.Budget, error) {
	query := `
		UPDATE budget
		SET
			amount = COALESCE($1, amount),
			budget_description = COALESCE($2, budget_description),
			period_start = COALESCE($3, period_start),
			period_end = COALESCE($4, period_end),
			updated_at = NOW()
		WHERE _id = $5 AND user_id = $6 AND closed_at IS NULL
		RETURNING _id, user_id, category_id, currency_id, amount, budget_description,
				  created_at, updated_at, period_start, period_end
	`

	var b bdgmodels.Budget
	err := r.db.QueryRowContext(ctx, query,
		req.Amount,
		req.Description,
		req.PeriodStart,
		req.PeriodEnd,
		req.BudgetID,
		req.UserID,
	).Scan(
		&b.ID,
		&b.UserID,
		&b.CategoryID,
		&b.CurrencyID,
		&b.Amount,
		&b.Description,
		&b.CreatedAt,
		&b.UpdatedAt,
		&b.PeriodStart,
		&b.PeriodEnd,
	)

	if err != nil {
		return bdgmodels.Budget{}, MapPgError(err)
	}

	return b, nil
}

func (r *PostgresRepository) DeleteBudget(ctx context.Context, budgetID int) (bdgmodels.Budget, error) {
	query := `
		UPDATE budget 
		SET closed_at = NOW(), updated_at = NOW()
		WHERE _id = $1
		RETURNING _id, user_id, category_id, currency_id, amount, budget_description,
				  created_at, updated_at, period_start, period_end
	`

	var b bdgmodels.Budget
	err := r.db.QueryRowContext(ctx, query, budgetID).Scan(
		&b.ID,
		&b.UserID,
		&b.CategoryID,
		&b.CurrencyID,
		&b.Amount,
		&b.Description,
		&b.CreatedAt,
		&b.UpdatedAt,
		&b.PeriodStart,
		&b.PeriodEnd,
	)

	if err != nil {
		return bdgmodels.Budget{}, MapPgError(err)
	}

	return b, nil
}
