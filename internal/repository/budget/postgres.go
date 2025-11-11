package budget

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	postgreserrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/errors"
	
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetBudgetsByUser(ctx context.Context, userID int) ([]models.Budget, error) {
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

	var budgets []models.Budget
	for rows.Next() {
		var budget BudgetDB
		if err := rows.Scan(
			&budget.ID,
			&budget.UserID,
			&budget.CategoryID,
			&budget.CurrencyID,
			&budget.Amount,
			&budget.Description,
			&budget.CreatedAt,
			&budget.UpdatedAt,
			&budget.ClosedAt,
			&budget.PeriodStart,
			&budget.PeriodEnd,
		); err != nil {
			return nil, fmt.Errorf("failed to scan budget: %w", err)
		}
		budgets = append(budgets, BudgetDBToModel(budget))
	}

	return budgets, nil
}

func (r *PostgresRepository) CreateBudget(ctx context.Context, budget models.Budget) (models.Budget, error) {
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
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case postgreserrors.UniqueViolation:
				return models.Budget{}, ErrUniqueViolation
			case postgreserrors.ForeignKeyViolation:
				return models.Budget{}, ErrForeignKeyViolation
			case postgreserrors.NotNullViolation:
				return models.Budget{}, ErrNotNullViolation
			case postgreserrors.CheckViolation:
				return models.Budget{}, ErrCheckViolation
			}
		}
		return models.Budget{}, fmt.Errorf("failed to create budget: %w", err)
	}

	return budget, nil
}

func (r *PostgresRepository) UpdateBudget(ctx context.Context, req models.UpdatedBudgetRequest, userID, budgetID int) (models.Budget, error) {
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

	var b models.Budget
	err := r.db.QueryRowContext(ctx, query,
		req.Amount,
		req.Description,
		req.PeriodStart,
		req.PeriodEnd,
		budgetID,
		userID,
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
		if errors.Is(err, sql.ErrNoRows) {
			return models.Budget{}, ErrBudgetNotFound
		}
		return models.Budget{}, fmt.Errorf("failed to update budget: %w", err)
	}

	return b, nil
}

func (r *PostgresRepository) DeleteBudget(ctx context.Context, budgetID int) (models.Budget, error) {
	query := `
		UPDATE budget 
		SET closed_at = NOW(), updated_at = NOW()
		WHERE _id = $1
		RETURNING _id, user_id, category_id, currency_id, amount, budget_description,
				  created_at, updated_at, period_start, period_end
	`

	var b models.Budget
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
		if errors.Is(err, sql.ErrNoRows) {
			return models.Budget{}, ErrBudgetNotFound
		}
		return models.Budget{}, fmt.Errorf("failed to delete budget: %w", err)
	}

	return b, nil
}
