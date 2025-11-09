package operation

import (
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type OperationAPI struct {
	ID          int       `json:"operation_id"`
	AccountID   int       `json:"account_id"`
	CategoryID  int       `json:"category_id"`
	RecevierID  int       `json:"receiver_id"`
	Sum         float64   `json:"sum"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	ReceiptURL  string    `json:"receipt_url"`
	CurrencyID  int       `json:"currency_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type OperationsAPI struct {
	UserID     int            `json:"user_id"`
	Operations []OperationAPI `json:"operations"`
}

func OperationToAPI(op models.Operation) OperationAPI {
	return OperationAPI{
		ID:          op.ID,
		AccountID:   op.AccountID,
		CategoryID:  op.CategoryID,
		Sum:         op.Sum,
		Name:        op.Name,
		Type:        string(op.Type),
		Status:      string(op.Status),
		Description: op.Description,
		ReceiptURL:  op.ReceiptURL,
		CurrencyID:  op.CurrencyID,
		CreatedAt:   op.CreatedAt,
	}
}

func OperationsToAPI(ops []models.Operation) OperationsAPI {
	operations := make([]OperationAPI, len(ops))
	for i, op := range ops {
		operations[i] = OperationToAPI(op)
	}

	return OperationsAPI{
		Operations: operations,
	}
}
