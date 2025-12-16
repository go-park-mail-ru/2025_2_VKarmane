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

type AccountType string

const (
	AccountShared  AccountType = "shared"
	AccountPrivate AccountType = "private"
)

type Operation struct {
	ID           int
	AccountID    int
	CategoryID   int
	CategoryName string
	Type         OperationType
	Status       OperationStatus
	AccountType  AccountType
	Description  string
	ReceiptURL   string
	Name         string
	Sum          float64
	CurrencyID   int
	CreatedAt    time.Time
	Date         time.Time
}

type OperationInList struct {
	ID                   int
	AccountID            int
	CategoryID           int
	CategoryName         string
	Type                 OperationType
	Description          string
	Name                 string
	CategoryLogoHashedID string
	CategoryLogo         string
	Sum                  float64
	CurrencyID           int
	AccountType          string
	CreatedAt            time.Time
	Date                 time.Time
}

type CreateOperationRequest struct {
	UserID      int
	AccountID   int
	CategoryID  *int
	Type        OperationType
	Name        string
	Description string
	Sum         float64
	Date        *time.Time
}

type UpdateOperationRequest struct {
	UserID      int
	AccountID   int
	OperationID int
	CategoryID  *int
	Name        *string
	Description *string
	Sum         *float64
	CreatedAt   *time.Time
}
