package operation

import (
	"context"
	"database/sql"

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
		SELECT _id, account_from_id, account_to_id, category_id, currency_id, 
		       operation_status, operation_type, operation_name, operation_description, 
		       receipt_url, sum, created_at
		FROM operation
		WHERE account_from_id = $1 OR account_to_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
		)
		if err != nil {
			return nil, err
		}
		operations = append(operations, OperationDBToModel(operation))
	}

	return operations, nil
}

func (r *PostgresRepository) CreateOperation(ctx context.Context, operation OperationDB) (int, error) {
	query := `
		INSERT INTO operation (account_from_id, account_to_id, category_id, currency_id, 
		                      operation_status, operation_type, operation_name, operation_description, 
		                      receipt_url, sum, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING _id
	`

	var id int
	err := r.db.QueryRowContext(ctx, query,
		operation.AccountFromID,
		operation.AccountToID,
		operation.CategoryID,
		operation.CurrencyID,
		operation.Status,
		operation.Type,
		operation.Name,
		operation.Description,
		operation.ReceiptURL,
		operation.Sum,
		operation.CreatedAt,
	).Scan(&id)

	return id, err
}

func (r *PostgresRepository) GetOperationsByUser(ctx context.Context, userID int) ([]models.Operation, error) {
	query := `
		SELECT o._id, o.account_from_id, o.account_to_id, o.category_id, o.currency_id, 
		       o.operation_status, o.operation_type, o.operation_name, o.operation_description, 
		       o.receipt_url, o.sum, o.created_at
		FROM operation o
		JOIN sharings s ON (o.account_from_id = s.account_id OR o.account_to_id = s.account_id)
		WHERE s.user_id = $1
		ORDER BY o.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
		)
		if err != nil {
			return nil, err
		}
		operations = append(operations, OperationDBToModel(operation))
	}

	return operations, nil
}

