package operation

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

func OperationDBToModel(operationDB OperationDB) models.Operation {
	var accountID, categoryID, currencyID int
	if operationDB.AccountFromID != nil {
		accountID = *operationDB.AccountFromID
	}
	if operationDB.CategoryID != nil {
		categoryID = *operationDB.CategoryID
	}
	if operationDB.CurrencyID != nil {
		currencyID = *operationDB.CurrencyID
	}

	return models.Operation{
		ID:          operationDB.ID,
		AccountID:   accountID,
		CategoryID:  categoryID,
		Type:        operationDB.Type,
		Status:      operationDB.Status,
		Description: operationDB.Description,
		ReceiptURL:  operationDB.ReceiptURL,
		Name:        operationDB.Name,
		Sum:         operationDB.Sum,
		CurrencyID:  currencyID,
		CreatedAt:   operationDB.CreatedAt,
	}
}

func OperationModelToDB(operation models.Operation) OperationDB {
	var accountFromID, categoryID, currencyID *int
	if operation.AccountID != 0 {
		accountFromID = &operation.AccountID
	}
	if operation.CategoryID != 0 {
		categoryID = &operation.CategoryID
	}
	if operation.CurrencyID != 0 {
		currencyID = &operation.CurrencyID
	}

	return OperationDB{
		ID:            operation.ID,
		AccountFromID: accountFromID,
		CategoryID:    categoryID,
		Type:          operation.Type,
		Status:        operation.Status,
		Description:   operation.Description,
		ReceiptURL:    operation.ReceiptURL,
		Name:          operation.Name,
		Sum:           operation.Sum,
		CurrencyID:    currencyID,
		CreatedAt:     operation.CreatedAt,
	}
}
