package budget

import "time"

type BudgetAPI struct {
	ID          int       `json:"budget_id"`
	UserID      int       `json:"user_id"`
	Amount      float64   `json:"amount"`
	Actual      float64   `json:"actual"`
	CurrencyID  int       `json:"currency_id"`
	Description string    `json:"description"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
}

type BudgetsAPI struct {
	UserID  int         `json:"user_id"`
	Budgets []BudgetAPI `json:"budgets"`
}
