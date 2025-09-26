package budget

import "time"

// BudgetAPI - DTO для API ответов с бюджетом
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

// BudgetsAPI - DTO для API ответов со списком бюджетов
type BudgetsAPI struct {
	UserID  int         `json:"user_id"`
	Budgets []BudgetAPI `json:"budgets"`
}
