package models

import "time"

const (
	WRITE  string = "create"
	DELETE string = "delete"
	UPDATE string = "update"
)

const (
	TRANSACTIONS string = "transactions"
	CATEGORIES   string = "categories"
)

type TransactionSearch struct {
	ID                   int       `json:"id"`
	AccountID            int       `json:"account_id"`
	CategoryID           int       `json:"category_id"`
	CategoryName         string    `json:"category_name"`
	Type                 string    `json:"type"`
	Description          string    `json:"description"`
	Status               string    `json:"status"`
	Name                 string    `json:"name"`
	CategoryLogoHashedID string    `json:"category_logo_hashed_id"`
	CategoryLogo         string    `json:"category_logo"`
	Sum                  float64   `json:"sum"`
	AccountType          string    `json:"account_type"`
	CurrencyID           int       `json:"curerncy"`
	CreatedAt            time.Time `json:"created_at"`
	Date                 time.Time `json:"date"`
	Action               string    `json:"action"`
}

type UpdateCategoryInOperationSearch struct {
	CategoryID           int    `json:"category_id"`
	CategoryName         string `json:"category_name"`
	CategoryLogoHashedID string `json:"category_logo_hashed_id"`
	CategoryLogo         string `json:"category_logo"`
	Action               string `json:"action"`
}
