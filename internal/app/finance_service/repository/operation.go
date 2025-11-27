package repository

import (
	"context"
	"time"

	finmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/models"
)

type OperationDB struct {
	ID            int
	AccountFromID *int
	AccountToID   *int
	CategoryID    *int
	CategoryName  string
	CurrencyID    *int
	Status        finmodels.OperationStatus
	Type          finmodels.OperationType
	Name          string
	Description   string
	ReceiptURL    string
	Sum           float64
	CreatedAt     time.Time
	Date          time.Time
}

func (r *PostgresRepository) GetOperationsByAccount(ctx context.Context, accountID int) ([]finmodels.OperationInList, error) {
	query := `
		SELECT o._id, o.account_from_id, o.account_to_id, o.category_id, o.currency_id, 
		       o.operation_status, o.operation_type, o.operation_name, o.operation_description, 
		       o.receipt_url, o.sum, o.created_at, o.operation_date,
		       COALESCE(c.category_name, 'Без категории') as category_name,
			   COALESCE(c.logo_hashed_id, '') AS logo_hashed_idc
		FROM operation o
		LEFT JOIN category c ON o.category_id = c._id
		WHERE (o.account_from_id = $1 OR o.account_to_id = $1) AND o.operation_status != 'reverted'
		ORDER BY o.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		return nil, MapPgOperationError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var operations []finmodels.OperationInList
	for rows.Next() {
		var opDB OperationDB
		var categoryLogoHashID string
		err := rows.Scan(
			&opDB.ID,
			&opDB.AccountFromID,
			&opDB.AccountToID,
			&opDB.CategoryID,
			&opDB.CurrencyID,
			&opDB.Status,
			&opDB.Type,
			&opDB.Name,
			&opDB.Description,
			&opDB.ReceiptURL,
			&opDB.Sum,
			&opDB.CreatedAt,
			&opDB.Date,
			&opDB.CategoryName,
			&categoryLogoHashID,
		)
		if err != nil {
			return nil, MapPgOperationError(err)
		}

		opInList := operationDBToModelInList(opDB, categoryLogoHashID)
		operations = append(operations, opInList)
	}

	return operations, nil
}

func (r *PostgresRepository) GetOperationByID(ctx context.Context, accID int, opID int) (finmodels.Operation, error) {
	query := `
		SELECT o._id, o.account_from_id, o.account_to_id, o.category_id, o.currency_id, 
		       o.operation_status, o.operation_type, o.operation_name, o.operation_description, 
		       o.receipt_url, o.sum, o.created_at, o.operation_date,
		       COALESCE(c.category_name, 'Без категории') as category_name
		FROM operation o
		LEFT JOIN category c ON o.category_id = c._id
		WHERE o._id = $1 AND (o.account_from_id = $2 OR o.account_to_id = $2) AND o.operation_status != 'reverted'
	`

	var operation OperationDB
	err := r.db.QueryRowContext(ctx, query, opID, accID).Scan(
		&operation.ID,
		&operation.AccountFromID,
		&operation.AccountToID,
		&operation.CategoryID,
		&operation.CurrencyID,
		&operation.Status,
		&operation.Type,
		&operation.Name,
		&operation.Description,
		&operation.ReceiptURL,
		&operation.Sum,
		&operation.CreatedAt,
		&operation.Date,
		&operation.CategoryName,
	)

	if err != nil {
		return finmodels.Operation{}, MapPgOperationError(err)
	}

	return operationDBToModel(operation), nil
}

func (r *PostgresRepository) CreateOperation(ctx context.Context, op finmodels.Operation) (finmodels.Operation, error) {
	query := `
		INSERT INTO operation (account_from_id, account_to_id, category_id, currency_id, 
		                      operation_status, operation_type, operation_name, operation_description, 
		                      receipt_url, sum, created_at, operation_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING _id
	`

	var categoryID interface{}
	if op.CategoryID > 0 {
		categoryID = op.CategoryID
	} else {
		categoryID = nil
	}

	var currencyID interface{}
	if op.CurrencyID > 0 {
		currencyID = op.CurrencyID
	} else {
		currencyID = nil
	}

	var id int
	err := r.db.QueryRowContext(ctx, query,
		op.AccountID,
		nil,
		categoryID,
		currencyID,
		op.Status,
		op.Type,
		op.Name,
		op.Description,
		op.ReceiptURL,
		op.Sum,
		time.Now(),
		op.Date,
	).Scan(&id)

	if err != nil {
		return finmodels.Operation{}, MapPgOperationError(err)
	}

	op.ID = id

	if op.CategoryID > 0 {
		var categoryName string
		err = r.db.QueryRowContext(ctx, "SELECT category_name FROM category WHERE _id = $1", op.CategoryID).Scan(&categoryName)
		if err != nil {
			op.CategoryName = "Без категории"
		} else {
			op.CategoryName = categoryName
		}
	} else {
		op.CategoryName = "Без категории"
	}

	return op, nil
}

func (r *PostgresRepository) UpdateOperation(ctx context.Context, req finmodels.UpdateOperationRequest, accID int, opID int) (finmodels.Operation, error) {
	query := `
		WITH updated_operation AS (
			UPDATE operation 
			SET 
			    operation_name = COALESCE($1, operation_name),
			    operation_description = COALESCE($2, operation_description),
			    sum = COALESCE($3, sum)
			WHERE _id = $4 AND (account_from_id = $5 OR account_to_id = $5) AND operation_status != 'reverted'
			RETURNING _id, account_from_id, account_to_id, category_id, currency_id, 
			          operation_status, operation_type, operation_name, operation_description, 
			          receipt_url, sum, created_at, operation_date
		)
		SELECT o._id, o.account_from_id, o.account_to_id, o.category_id, o.currency_id, 
		       o.operation_status, o.operation_type, o.operation_name, o.operation_description, 
		       o.receipt_url, o.sum, o.created_at, o.operation_date,
		       COALESCE(c.category_name, 'Без категории') as category_name
		FROM updated_operation o
		LEFT JOIN category c ON o.category_id = c._id
	`

	var name *string
	if req.Name != nil {
		name = req.Name
	}

	var description *string
	if req.Description != nil {
		description = req.Description
	}

	var sum *float64
	if req.Sum != nil {
		sum = req.Sum
	}

	var operation OperationDB
	err := r.db.QueryRowContext(ctx, query,
		name,
		description,
		sum,
		opID,
		accID,
	).Scan(
		&operation.ID,
		&operation.AccountFromID,
		&operation.AccountToID,
		&operation.CategoryID,
		&operation.CurrencyID,
		&operation.Status,
		&operation.Type,
		&operation.Name,
		&operation.Description,
		&operation.ReceiptURL,
		&operation.Sum,
		&operation.CreatedAt,
		&operation.Date,
		&operation.CategoryName,
	)

	if err != nil {
		return finmodels.Operation{}, MapPgOperationError(err)
	}

	return operationDBToModel(operation), nil
}

func (r *PostgresRepository) DeleteOperation(ctx context.Context, accID int, opID int) (finmodels.Operation, error) {
	query := `
		WITH updated_operation AS (
			UPDATE operation 
			SET operation_status = 'reverted'
			WHERE _id = $1 AND (account_from_id = $2 OR account_to_id = $2) AND operation_status != 'reverted'
			RETURNING _id, account_from_id, account_to_id, category_id, currency_id, 
			          operation_status, operation_type, operation_name, operation_description, 
			          receipt_url, sum, created_at, operation_date
		)
		SELECT o._id, o.account_from_id, o.account_to_id, o.category_id, o.currency_id, 
		       o.operation_status, o.operation_type, o.operation_name, o.operation_description, 
		       o.receipt_url, o.sum, o.created_at, o.operation_date,
		       COALESCE(c.category_name, 'Без категории') as category_name
		FROM updated_operation o
		LEFT JOIN category c ON o.category_id = c._id
	`

	var operation OperationDB
	err := r.db.QueryRowContext(ctx, query, opID, accID).Scan(
		&operation.ID,
		&operation.AccountFromID,
		&operation.AccountToID,
		&operation.CategoryID,
		&operation.CurrencyID,
		&operation.Status,
		&operation.Type,
		&operation.Name,
		&operation.Description,
		&operation.ReceiptURL,
		&operation.Sum,
		&operation.CreatedAt,
		&operation.Date,
		&operation.CategoryName,
	)

	if err != nil {
		return finmodels.Operation{}, MapPgOperationError(err)
	}

	return operationDBToModel(operation), nil
}

func operationDBToModel(operationDB OperationDB) finmodels.Operation {
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

	return finmodels.Operation{
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

func operationDBToModelInList(opDB OperationDB, categoryLogoHash string) finmodels.OperationInList {
	var accountID int
	if opDB.AccountFromID != nil {
		accountID = *opDB.AccountFromID
	} else if opDB.AccountToID != nil {
		accountID = *opDB.AccountToID
	}

	var categoryID int
	if opDB.CategoryID != nil {
		categoryID = *opDB.CategoryID
	}

	var currencyID int
	if opDB.CurrencyID != nil {
		currencyID = *opDB.CurrencyID
	}

	return finmodels.OperationInList{
		ID:                   opDB.ID,
		AccountID:            accountID,
		CategoryID:           categoryID,
		CategoryName:         opDB.CategoryName,
		Type:                 opDB.Type,
		Description:          opDB.Description,
		Name:                 opDB.Name,
		CategoryLogoHashedID: categoryLogoHash,
		Sum:                  opDB.Sum,
		CurrencyID:           currencyID,
		CreatedAt:            opDB.CreatedAt,
		Date:                 opDB.Date,
	}
}
