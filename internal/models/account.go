package models

import "time"

type AccountType string

const (
	PrivateAccount = "private"
	SharedAccount  = "shared"
)

type Account struct {
	ID         int
	Balance    float64
	Type       string
	CurrencyID int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type UserAccount struct {
	ID        int
	UserID    int
	AccountID int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateAccountRequest struct {
	Name       string      `json:"name"`
	Balance    float64     `json:"balance" validate:"min=0"`
	Type       AccountType `json:"type"`
	CurrencyID int         `json:"currency_id"`
}

type UpdateAccountRequest struct {
	Name    *string  `json:"name,omitempty"`
	Balance *float64 `json:"balance,omitempty" validate:"min=0"`
}
