package models

import "time"

type User struct {
	ID        int       `json:"user_id"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Email     string    `json:"email"`
	Login     string    `json:"login"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Currency struct {
	ID        int       `json:"currency_id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Account struct {
	ID         int       `json:"account_id"`
	Balance    float64   `json:"balance"`
	Type       string    `json:"type"`
	CurrencyID int       `json:"currency_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserAccount struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	AccountID int       `json:"account_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Category struct {
	ID          int       `json:"category_id"`
	UserID      int       `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Operation struct {
	ID          int       `json:"operation_id"`
	AccountID   int       `json:"account_id"`
	CategoryID  int       `json:"category_id"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	ReceiptURL  string    `json:"receipt_url"`
	Name        string    `json:"name"`
	Sum         float64   `json:"sum"`
	CurrencyID  int       `json:"currency_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type Budget struct {
	ID          int       `json:"budget_id"`
	UserID      int       `json:"user_id"`
	Amount      float64   `json:"amount"`
	CurrencyID  int       `json:"currency_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ClosedAt    time.Time `json:"closed_at,omitempty"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
}
