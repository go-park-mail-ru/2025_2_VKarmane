package models

import "time"

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

