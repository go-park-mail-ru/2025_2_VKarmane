package operation

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresRepository_GetOperationsByAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	accountID := 1
	now := time.Now()
	accFromID := 1
	var accToID *int = nil
	categoryID := 5
	currencyID := 1
	rows := sqlmock.NewRows([]string{"o._id", "o.account_from_id", "o.account_to_id", "o.category_id", "o.currency_id", "o.operation_status", "o.operation_type", "o.operation_name", "o.operation_description", "o.receipt_url", "o.sum", "o.created_at", "o.operation_date", "category_name"}).
		AddRow(1, accFromID, accToID, categoryID, currencyID, models.OperationFinished, models.OperationExpense, "Test Op", "Test Desc", "", 100.50, now, now, "Food").
		AddRow(2, accFromID, accToID, categoryID, currencyID, models.OperationFinished, models.OperationIncome, "Test Op 2", "", "", 200.75, now, now, "Food")

	mock.ExpectQuery(`SELECT.*FROM operation`).
		WithArgs(accountID).
		WillReturnRows(rows)

	operations, err := repo.GetOperationsByAccount(context.Background(), accountID)
	assert.NoError(t, err)
	assert.Len(t, operations, 2)
	assert.Equal(t, 1, operations[0].ID)
	assert.Equal(t, "Test Op", operations[0].Name)
	assert.Equal(t, 100.50, operations[0].Sum)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetOperationsByAccount_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	accountID := 1
	rows := sqlmock.NewRows([]string{"o._id", "o.account_from_id", "o.account_to_id", "o.category_id", "o.currency_id", "o.operation_status", "o.operation_type", "o.operation_name", "o.operation_description", "o.receipt_url", "o.sum", "o.created_at", "o.operation_date", "category_name"})

	mock.ExpectQuery(`SELECT.*FROM operation`).
		WithArgs(accountID).
		WillReturnRows(rows)

	operations, err := repo.GetOperationsByAccount(context.Background(), accountID)
	assert.NoError(t, err)
	assert.Empty(t, operations)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetOperationsByAccount_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	accountID := 1
	mock.ExpectQuery(`SELECT.*FROM operation`).
		WithArgs(accountID).
		WillReturnError(sql.ErrConnDone)

	operations, err := repo.GetOperationsByAccount(context.Background(), accountID)
	assert.Error(t, err)
	assert.Nil(t, operations)
	assert.Contains(t, err.Error(), "failed to get operations by account")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetOperationByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	accountID := 1
	operationID := 5
	now := time.Now()
	accFromID := 1
	var accToID *int = nil
	categoryID := 5
	currencyID := 1
	rows := sqlmock.NewRows([]string{"o._id", "o.account_from_id", "o.account_to_id", "o.category_id", "o.currency_id", "o.operation_status", "o.operation_type", "o.operation_name", "o.operation_description", "o.receipt_url", "o.sum", "o.created_at", "o.operation_date", "category_name"}).
		AddRow(operationID, accFromID, accToID, categoryID, currencyID, models.OperationFinished, models.OperationExpense, "Test Op", "Test Desc", "", 100.50, now, now, "Food")

	mock.ExpectQuery(`SELECT.*FROM operation`).
		WithArgs(operationID, accountID).
		WillReturnRows(rows)

	operation, err := repo.GetOperationByID(context.Background(), accountID, operationID)
	assert.NoError(t, err)
	assert.Equal(t, operationID, operation.ID)
	assert.Equal(t, "Test Op", operation.Name)
	assert.Equal(t, 100.50, operation.Sum)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetOperationByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	accountID := 1
	operationID := 999
	mock.ExpectQuery(`SELECT.*FROM operation`).
		WithArgs(operationID, accountID).
		WillReturnError(sql.ErrNoRows)

	operation, err := repo.GetOperationByID(context.Background(), accountID, operationID)
	assert.Error(t, err)
	assert.Zero(t, operation.ID)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

