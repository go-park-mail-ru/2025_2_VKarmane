package operation

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

func OperationDBToModel(operationDB OperationDB) models.Operation {
	return models.Operation{
		ID:          operationDB.ID,
		AccountID:   operationDB.AccountID,
		CategoryID:  operationDB.CategoryID,
		Type:       models.OperationType(operationDB.Type),
		Status:      models.OperationStatus(operationDB.Status),
		Description: operationDB.Description,
		ReceiptURL:  operationDB.ReceiptURL,
		Name:        operationDB.Name,
		Sum:         operationDB.Sum,
		CurrencyID:  operationDB.CurrencyID,
		CreatedAt:   operationDB.CreatedAt,
	}
}

func OperationModelToDB(operation models.Operation) OperationDB {
	return OperationDB{
		ID:          operation.ID,
		AccountID:   operation.AccountID,
		CategoryID:  operation.CategoryID,
		Type:        string(operation.Type),
		Status:      string(operation.Status),
		Description: operation.Description,
		ReceiptURL:  operation.ReceiptURL,
		Name:        operation.Name,
		Sum:         operation.Sum,
		CurrencyID:  operation.CurrencyID,
		CreatedAt:   operation.CreatedAt,
	}
}
