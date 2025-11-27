package models

import "time"

type TransactionSearch struct {
	ID                   int       `json:"id"`
	AccountID            int       `json:"account_id"`
	CategoryID           int       `json:"category_id"`
	CategoryName         string    `json:"category_name"`
	Type                 string    `json:"type"`
	Description          string    `json:"description"`
	Name                 string    `json:"name"`
	CategoryLogoHashedID string    `json:"category_logo_hashed_id"`
	CategoryLogo         string    `json:"category_logo"`
	Sum                  float64   `json:"sum"`
	CurrencyID           int       `json:"curerncy"`
	CreatedAt            time.Time `json:"created_at"`
	Date                 time.Time `json:"date"`
}
