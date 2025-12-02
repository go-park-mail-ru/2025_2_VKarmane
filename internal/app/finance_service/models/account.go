package models

import "time"

type Account struct {
	ID          int
	Balance     float64
	Type        string
	Name        string
	Description *string
	CurrencyID  int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateAccountRequest struct {
	UserID      int
	Balance     float64
	Name        string
	Description string
	Type        string
	CurrencyID  int
}

type UpdateAccountRequest struct {
	UserID      int
	AccountID   int
	Description *string
	Name        *string
	Balance     *float64
}
