package models

import (
	"fmt"
	"time"
)

type OperationType string

const (
	OperationIncome  OperationType = "income"
	OperationExpense OperationType = "expense"
)

type OperationStatus string

const (
	OperationFinished OperationStatus = "finished"
	OperationReverted OperationStatus = "reverted"
)

type Operation struct {
	ID           int
	AccountID    int
	CategoryID   int
	CategoryName string
	Type         OperationType
	Status       OperationStatus
	Description  string
	ReceiptURL   string
	Name         string
	Sum          float64
	CurrencyID   int
	CreatedAt    time.Time
	Date         time.Time // Дата операции (может отличаться от created_at)
}

type OperationInList struct {
	ID           int
	AccountID    int
	CategoryID   int
	CategoryName string
	Type         OperationType
	Description  string
	Name         string
	CategoryLogoHashedID string
	CategoryLogo string
	Sum          float64
	CurrencyID   int
	CreatedAt    time.Time
	Date         time.Time // Дата операции (может отличаться от created_at)
}


type UpdateOperationRequest struct {
	CategoryID  *int       `json:"category_id,omitempty"`
	Name        *string    `json:"name,omitempty" validate:"omitempty,max=50"`
	Description *string    `json:"description,omitempty" validate:"omitempty,max=60"`
	Sum         *float64   `json:"sum,omitempty" validate:"omitempty,min=0"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
}

func (r *UpdateOperationRequest) Validate() error {
	if r.Sum != nil && *r.Sum <= 0 {
		return fmt.Errorf("invalid sum: %s", ErrCodeInvalidAmount)
	}
	return nil
}

type DeleteOperationRequest struct {
	Status string `json:"status"`
}

type CreateOperationRequest struct {
	AccountID   int           `json:"account_id"`
	CategoryID  *int          `json:"category_id,omitempty"`
	Type        OperationType `json:"type"`
	Name        string        `json:"name" validate:"required,max=50"`
	Description string        `json:"description,omitempty" validate:"max=60"`
	Sum         float64       `json:"sum" validate:"min=0"`
	Date        *time.Time    `json:"date,omitempty"` // Дата операции
}

func (r *CreateOperationRequest) Validate() error {
	if r.Sum < 0 {
		return fmt.Errorf("invalid sum: %s", ErrCodeInvalidAmount)
	}
	if r.Name == "" {
		return fmt.Errorf("name is required: %s", ErrCodeMissingFields)
	}
	return nil
}

type OperationResponse struct {
	ID           int       `json:"id"`
	AccountID    int       `json:"account_id"`
	CategoryID   int       `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Type         string    `json:"type"`
	Status       string    `json:"status"`
	Description  string    `json:"description"`
	ReceiptURL   string    `json:"receipt_url"`
	Name         string    `json:"name"`
	Sum          float64   `json:"sum"`
	CurrencyID   int       `json:"currency_id"`
	CreatedAt    time.Time `json:"created_at"`
	Date         time.Time `json:"date"` // Дата операции
}


type OperationInListResponse struct {
	ID           int       `json:"id"`
	AccountID    int       `json:"account_id"`
	CategoryID   int       `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Type         string    `json:"type"`
	Description  string    `json:"description"`
	CategoryLogo string		`json:"category_logo"`	
	CategoryHashedID string `json:"-"`
	Sum          float64   `json:"sum"`
	CurrencyID   int       `json:"currency_id"`
	CreatedAt    time.Time `json:"created_at"`
	Date         time.Time `json:"date"` // Дата операции
}
