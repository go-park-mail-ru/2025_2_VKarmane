package operation

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

func OperationDBToModel(operationDB OperationDB) models.Operation {
	var accountID int
	if operationDB.AccountFromID != nil {
		accountID = *operationDB.AccountFromID
	} else if operationDB.AccountToID != nil {
		accountID = *operationDB.AccountToID
	}

	var categoryID int
	if operationDB.CategoryID != nil {
		categoryID = *operationDB.CategoryID
	}

	var currencyID int
	if operationDB.CurrencyID != nil {
		currencyID = *operationDB.CurrencyID
	}

	return models.Operation{
		ID:           operationDB.ID,
		AccountID:    accountID,
		CategoryID:   categoryID,
		CategoryName: operationDB.CategoryName,
		Type:         operationDB.Type,
		Status:       operationDB.Status,
		Description:  operationDB.Description,
		ReceiptURL:   operationDB.ReceiptURL,
		Name:         operationDB.Name,
		Sum:          operationDB.Sum,
		CurrencyID:   currencyID,
		CreatedAt:    operationDB.CreatedAt,
		Date:         operationDB.Date,
	}
}

func OperationModelToDB(operation models.Operation) OperationDB {
	return OperationDB{
		ID:            operation.ID,
		AccountFromID: &operation.AccountID,
		AccountToID:   nil,
		CategoryID:    &operation.CategoryID,
		CategoryName:  operation.CategoryName,
		Type:          operation.Type,
		Status:        operation.Status,
		Description:   operation.Description,
		ReceiptURL:    operation.ReceiptURL,
		Name:          operation.Name,
		Sum:           operation.Sum,
		CurrencyID:    &operation.CurrencyID,
		CreatedAt:     operation.CreatedAt,
		Date:          operation.Date,
	}
}
