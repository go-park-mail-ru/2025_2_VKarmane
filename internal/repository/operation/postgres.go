package operation

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) GetOperationsByAccount(ctx context.Context, accountID int) ([]models.Operation, error) {
	query := `
		SELECT o._id, o.account_from_id, o.account_to_id, o.category_id, o.currency_id, 
		       o.operation_status, o.operation_type, o.operation_name, o.operation_description, 
		       o.receipt_url, o.sum, o.created_at, o.operation_date,
		       COALESCE(c.category_name, 'Без категории') as category_name
		FROM operation o
		LEFT JOIN category c ON o.category_id = c._id
		WHERE (o.account_from_id = $1 OR o.account_to_id = $1) AND o.operation_status != 'reverted'
		ORDER BY o.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var operations []models.Operation
	for rows.Next() {
		var operation OperationDB
		err := rows.Scan(
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
			return nil, err
		}
		operations = append(operations, OperationDBToModel(operation))
	}

	return operations, nil
}

func (r *PostgresRepository) GetOperationByID(ctx context.Context, accID int, opID int) (models.Operation, error) {
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
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Operation{}, sql.ErrNoRows
		}
		return models.Operation{}, err
	}

	return OperationDBToModel(operation), nil
}

func (r *PostgresRepository) CreateOperation(ctx context.Context, op models.Operation) (models.Operation, error) {
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
		return models.Operation{}, err
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

func (r *PostgresRepository) UpdateOperation(ctx context.Context, req models.UpdateOperationRequest, accID int, opID int) (models.Operation, error) {
	existingOp, err := r.GetOperationByID(ctx, accID, opID)
	if err != nil {
		return models.Operation{}, err
	}

	if req.CategoryID != nil {
		existingOp.CategoryID = *req.CategoryID
	}
	if req.Name != nil {
		existingOp.Name = *req.Name
	}
	if req.Description != nil {
		existingOp.Description = *req.Description
	}
	if req.Sum != nil {
		existingOp.Sum = *req.Sum
	}

	query := `
		UPDATE operation 
		SET category_id = $1, operation_name = $2, operation_description = $3, sum = $4
		WHERE _id = $5 AND (account_from_id = $6 OR account_to_id = $6) AND operation_status != 'reverted'
		RETURNING _id, account_from_id, account_to_id, category_id, currency_id, 
		          operation_status, operation_type, operation_name, operation_description, 
		          receipt_url, sum, created_at, operation_date
	`

	var categoryID interface{}
	if existingOp.CategoryID > 0 {
		categoryID = existingOp.CategoryID
	} else {
		categoryID = nil
	}

	var operation OperationDB
	err = r.db.QueryRowContext(ctx, query,
		categoryID,
		existingOp.Name,
		existingOp.Description,
		existingOp.Sum,
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
	)

	if err != nil {
		return models.Operation{}, err
	}

	if operation.CategoryID != nil && *operation.CategoryID > 0 {
		var categoryName string
		err = r.db.QueryRowContext(ctx, "SELECT category_name FROM category WHERE _id = $1", *operation.CategoryID).Scan(&categoryName)
		if err != nil {
			operation.CategoryName = "Без категории"
		} else {
			operation.CategoryName = categoryName
		}
	} else {
		operation.CategoryName = "Без категории"
	}

	return OperationDBToModel(operation), nil
}

func (r *PostgresRepository) DeleteOperation(ctx context.Context, accID int, opID int) (models.Operation, error) {
	operation, err := r.GetOperationByID(ctx, accID, opID)
	if err != nil {
		return models.Operation{}, err
	}

	query := `
		UPDATE operation 
		SET operation_status = 'reverted'
		WHERE _id = $1 AND (account_from_id = $2 OR account_to_id = $2) AND operation_status != 'reverted'
	`

	_, err = r.db.ExecContext(ctx, query, opID, accID)
	if err != nil {
		return models.Operation{}, err
	}

	operation.Status = models.OperationReverted
	return operation, nil
}

func (r *PostgresRepository) GetOperationsByUser(ctx context.Context, userID int) ([]models.Operation, error) {
	query := `
		SELECT o._id, o.account_from_id, o.account_to_id, o.category_id, o.currency_id, 
		       o.operation_status, o.operation_type, o.operation_name, o.operation_description, 
		       o.receipt_url, o.sum, o.created_at, o.operation_date,
		       COALESCE(c.category_name, 'Без категории') as category_name
		FROM operation o
		LEFT JOIN category c ON o.category_id = c._id
		JOIN account a ON (o.account_from_id = a._id OR o.account_to_id = a._id)
		WHERE a.user_id = $1 AND o.operation_status != 'reverted'
		ORDER BY o.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var operations []models.Operation
	for rows.Next() {
		var operation OperationDB
		err := rows.Scan(
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
			return nil, err
		}
		operations = append(operations, OperationDBToModel(operation))
	}

	return operations, nil
}
