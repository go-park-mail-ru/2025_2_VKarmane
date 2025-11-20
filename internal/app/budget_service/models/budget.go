package models

import (
	"time"
)

type Budget struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	CategoryID  int       `json:"category_id"`
	Amount      float64   `json:"sum"`
	Actual      float64   // Фактические расходы (вычисляемое поле)
	CurrencyID  int       `json:"currency_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ClosedAt    time.Time `json:"closed_at"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
}

type CreateBudgetRequest struct {
	CategoryID  int       `json:"category_id"`
	Amount      float64   `json:"sum" validate:"min=0"`
	Description string    `json:"description,omitempty" validate:"max=80"`
	CreatedAt   time.Time `json:"created_at"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
}

type UpdatedBudgetRequest struct {
	UserID      int
	BudgetID    int
	Amount      *float64   `json:"sum,omitempty" validate:"min=0"`
	Description *string    `json:"description,omitempty" validate:"max=80"`
	PeriodStart *time.Time `json:"period_start,omitempty"`
	PeriodEnd   *time.Time `json:"period_end,omitempty"`
}
