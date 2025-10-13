package operation

import "time"

type OperationAPI struct {
	ID int `json:"transaction_id"`
	AccountID int `json:"account_id"`
	CategoryID int `json:"category_id"`
	Sum float64 `json:"sum"`
	Name string `json:"name"`
	Type string `json:"type"`
	Status string `json:"status"`
	Description string `json:"description"`
	ReceiptURL string `json:"receipt"`
	Date time.Time `json:"date"`
}

type OperationsAPI struct {
	UserID int `json:"user_id"`
	Operations []OperationAPI `json:"operations"`
}