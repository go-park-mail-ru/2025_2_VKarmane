package models

import "time"

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
	ID          int
	AccountID   int
	CategoryID  int
	ReceiverID  int
	Type        OperationType
	Status      OperationStatus
	Description string
	ReceiptURL  string
	Name        string
	Sum         float64
	CurrencyID  int
	CreatedAt   time.Time
}

type UpdateOperationRequest struct {
	CategoryID  *int       `json:"category_id,omitempty"`
	Name        *string    `json:"name,omitempty" validate:"omitempty,max=50"`
	Description *string    `json:"description,omitempty" validate:"omitempty,max=60"`
	Sum         *float64   `json:"sum,omitempty" validate:"omitempty,min=0"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
}

type DeleteOperationRequest struct {
	Status string `json:"status"`
}

type CreateOperationRequest struct {
	AccountID   int           `json:"account_id"`
	CategoryID  int           `json:"category_id"`
	ReceiverID  int           `json:"receiver_id,omitempty"`
	Type        OperationType `json:"type" validate:"required"`
	Name        string        `json:"name" validate:"required,max=50"`
	Description string        `json:"description,omitempty" validate:"max=60"`
	Sum         float64       `json:"sum" validate:"required,min=0"`
	CreatedAt   time.Time     `json:"created_at"`
}
