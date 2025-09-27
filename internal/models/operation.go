package models

import "time"

type Operation struct {
	ID          int
	AccountID   int
	CategoryID  int
	Type        string
	Status      string
	Description string
	ReceiptURL  string
	Name        string
	Sum         float64
	CurrencyID  int
	CreatedAt   time.Time
}
