package budget

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) GetBudgetsByUser(ctx context.Context, userID int) ([]models.Budget, error) {
	query := `
		SELECT _id, user_id, category_id, currency_id, amount, budget_description, 
		       created_at, updated_at, closed_at, period_start, period_end
		FROM budget
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get budgets by user: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var budgets []models.Budget
	for rows.Next() {
		var budget BudgetDB
		err := rows.Scan(
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
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan budget: %w", err)
		}
		budgets = append(budgets, BudgetDBToModel(budget))
	}

	return budgets, nil
}

func (r *PostgresRepository) CreateBudget(ctx context.Context, budget BudgetDB) (int, error) {
	query := `
		INSERT INTO budget (user_id, category_id, currency_id, amount, budget_description, 
		                   created_at, updated_at, period_start, period_end)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING _id
	`

	var id int
	err := r.db.QueryRowContext(ctx, query,
		budget.UserID,
		budget.CategoryID,
		budget.CurrencyID,
		budget.Amount,
		budget.Description,
		budget.CreatedAt,
		budget.UpdatedAt,
		budget.PeriodStart,
		budget.PeriodEnd,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create budget: %w", err)
	}

	return id, nil
}

func (r *PostgresRepository) UpdateBudget(ctx context.Context, budget BudgetDB) error {
	query := `
		UPDATE budget 
		SET amount = $1, budget_description = $2, updated_at = NOW()
		WHERE _id = $3
	`

	_, err := r.db.ExecContext(ctx, query, budget.Amount, budget.Description, budget.ID)
	if err != nil {
		return fmt.Errorf("failed to update budget: %w", err)
	}

	return nil
}

func (r *PostgresRepository) CloseBudget(ctx context.Context, budgetID int) error {
	query := `
		UPDATE budget 
		SET closed_at = NOW(), updated_at = NOW()
		WHERE _id = $1
	`

	_, err := r.db.ExecContext(ctx, query, budgetID)
	if err != nil {
		return fmt.Errorf("failed to close budget: %w", err)
	}

	return nil
}

