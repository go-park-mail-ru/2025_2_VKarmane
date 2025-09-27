package models

import "time"

type Budget struct {
	ID          int
	UserID      int
	Amount      float64
	Actual      float64 // Фактические расходы (вычисляемое поле)
	CurrencyID  int
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ClosedAt    time.Time
	PeriodStart time.Time
	PeriodEnd   time.Time
}
