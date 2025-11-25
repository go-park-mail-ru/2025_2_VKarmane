package budget

import "time"

type BudgetDB struct {
	ID          int        `db:"budget_id"`
	UserID      int        `db:"user_id"`
	CategoryID  int        `db:"category_id"`
	Amount      float64    `db:"amount"`
	CurrencyID  int        `db:"currency_id"`
	Description string     `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	ClosedAt    *time.Time `db:"closed_at"`
	PeriodStart time.Time  `db:"period_start"`
	PeriodEnd   time.Time  `db:"period_end"`
}
