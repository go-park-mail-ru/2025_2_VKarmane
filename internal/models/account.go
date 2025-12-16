package models

import "time"

type AccountType string

const (
	PrivateAccount = "private"
	SharedAccount  = "shared"
)

type Account struct {
	ID          int
	Balance     float64
	Name        string
	Description string
	Type        string
	CurrencyID  int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserAccount struct {
	ID        int
	UserID    int
	AccountID int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateAccountRequest struct {
	Name        string      `json:"name"`
	Description *string     `json:"description,omitempty"`
	Balance     float64     `json:"balance" validate:"min=0"`
	Type        AccountType `json:"type"`
	CurrencyID  int         `json:"currency_id"`
}

type UpdateAccountRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Balance     *float64 `json:"balance,omitempty" validate:"min=0"`
}

type AddUserToAccountRequest struct {
	UserLogin string `json:"user_login"`
	AccountID int    `json:"account_id"`
}
