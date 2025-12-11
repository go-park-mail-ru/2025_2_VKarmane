package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	serviceerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/errors"
	finmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/models"

	"github.com/stretchr/testify/require"
)

func setupOperationDB(t *testing.T) (*PostgresRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	return NewPostgresRepository(db), mock, func() { _ = db.Close() }
}

func TestGetOperationsByAccount_Success(t *testing.T) {
	repo, mock, close := setupOperationDB(t)
	defer close()

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"_id", "account_from_id", "account_to_id", "category_id", "currency_id",
		"operation_status", "operation_type", "operation_name", "operation_description",
		"receipt_url", "sum", "created_at", "operation_date",
		"category_name", "logo_hashed_idc",
		"account_type",
	}).AddRow(
		10, 1, nil, 3, 1, "done", "expense", "Оплата", "Описание",
		"url", 100.5, now, now, "Еда", "hash123", "cash",
	)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT o._id, o.account_from_id, o.account_to_id, o.category_id, o.currency_id, 
		       o.operation_status, o.operation_type, o.operation_name, o.operation_description,
		       o.receipt_url, o.sum, o.created_at, o.operation_date, 
			   COALESCE(c.category_name, 'Без категории') as category_name,
			   COALESCE(c.logo_hashed_id, '') AS logo_hashed_idc,
			   a.account_type
		FROM operation o
		LEFT JOIN category c ON o.category_id = c._id
		JOIN account a ON a._id = o.account_from_id
		WHERE (o.account_from_id = $1 OR o.account_to_id = $1) AND o.operation_status != 'reverted'
		ORDER BY o.created_at DESC
	`)).
		WithArgs(1).
		WillReturnRows(rows)

	ops, err := repo.GetOperationsByAccount(context.Background(), 1)
	require.NoError(t, err)
	require.Len(t, ops, 1)
	require.Equal(t, 10, ops[0].ID)
	require.Equal(t, 100.5, ops[0].Sum)
	require.Equal(t, "Еда", ops[0].CategoryName)
	require.Equal(t, "hash123", ops[0].CategoryLogoHashedID)
	require.Equal(t, "cash", ops[0].AccountType)
}

func TestGetOperationsByAccount_DBError(t *testing.T) {
	repo, mock, close := setupOperationDB(t)
	defer close()

	mock.ExpectQuery("SELECT").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetOperationsByAccount(context.Background(), 1)
	require.ErrorIs(t, err, serviceerrors.ErrOperationNotFound)
}

func TestGetOperationByID_Success(t *testing.T) {
	repo, mock, close := setupOperationDB(t)
	defer close()

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"_id", "account_from_id", "account_to_id", "category_id", "currency_id",
		"operation_status", "operation_type", "operation_name", "operation_description",
		"receipt_url", "sum", "created_at", "operation_date", "category_name",
	}).
		AddRow(
			10, 1, nil, 3, 1, "done", "income", "Зачисление", "Описание",
			"url", 200.5, now, now, "Зарплата",
		)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT o._id, o.account_from_id, o.account_to_id, o.category_id, o.currency_id,
		       o.operation_status, o.operation_type, o.operation_name, o.operation_description,
		       o.receipt_url, o.sum, o.created_at, o.operation_date,
		       COALESCE(c.category_name, 'Без категории') as category_name
		FROM operation o
		LEFT JOIN category c ON o.category_id = c._id
		WHERE o._id = $1 AND (o.account_from_id = $2 OR o.account_to_id = $2) AND o.operation_status != 'reverted'
	`)).
		WithArgs(10, 1).
		WillReturnRows(rows)

	op, err := repo.GetOperationByID(context.Background(), 1, 10)
	require.NoError(t, err)
	require.Equal(t, 200.5, op.Sum)
	require.Equal(t, "Зарплата", op.CategoryName)
}

func TestGetOperationByID_NotFound(t *testing.T) {
	repo, mock, close := setupOperationDB(t)
	defer close()

	mock.ExpectQuery("SELECT").
		WithArgs(10, 1).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetOperationByID(context.Background(), 1, 10)
	require.ErrorIs(t, err, serviceerrors.ErrOperationNotFound)
}

func TestCreateOperation_Success(t *testing.T) {
	repo, mock, close := setupOperationDB(t)
	defer close()

	now := time.Now()

	op := finmodels.Operation{
		AccountID:   1,
		CategoryID:  2,
		CurrencyID:  1,
		Status:      "done",
		Type:        "expense",
		Name:        "Покупка",
		Description: "Описание",
		Sum:         99.9,
		Date:        now,
	}

	// Начало транзакции
	mock.ExpectBegin()

	// Обновление баланса счета — теперь ExpectExec
	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE account
		SET balance = balance - $1, updated_at = NOW()
		WHERE _id = $2
		RETURNING balance
	`)).
		WithArgs(op.Sum, op.AccountID).
		WillReturnResult(sqlmock.NewResult(1, 1)) // 1 row affected

	// Вставка operation
	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO operation (
		    account_from_id, account_to_id, category_id, currency_id, 
		    operation_status, operation_type, operation_name, operation_description, 
		    receipt_url, sum, created_at, operation_date
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		RETURNING _id
	`)).
		WithArgs(
			op.AccountID,
			nil,
			op.CategoryID,
			op.CurrencyID,
			op.Status,
			op.Type,
			op.Name,
			op.Description,
			op.ReceiptURL,
			op.Sum,
			sqlmock.AnyArg(),
			op.Date,
		).
		WillReturnRows(sqlmock.NewRows([]string{"_id"}).AddRow(100))

	// Коммит транзакции
	mock.ExpectCommit()

	// Запрос имени категории
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT category_name FROM category WHERE _id = $1
	`)).
		WithArgs(op.CategoryID).
		WillReturnRows(sqlmock.NewRows([]string{"category_name"}).AddRow("Категория"))

	// Запрос account_type
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT account_type from account where _id = $1
	`)).
		WithArgs(op.AccountID).
		WillReturnRows(sqlmock.NewRows([]string{"account_type"}).AddRow("cash"))

	res, err := repo.CreateOperation(context.Background(), op)
	require.NoError(t, err)
	require.Equal(t, 100, res.ID)
	require.Equal(t, "Категория", res.CategoryName)
	require.Equal(t, "cash", string(res.AccountType))
}

func TestUpdateOperation_Success(t *testing.T) {
	repo, mock, close := setupOperationDB(t)
	defer close()

	name := "Новое имя"
	sum := 123.4
	req := finmodels.UpdateOperationRequest{
		Name: &name,
		Sum:  &sum,
	}

	now := time.Now()

	mock.ExpectBegin()

	rows := sqlmock.NewRows([]string{
    "_id", "account_from_id", "account_to_id", "category_id", "currency_id",
    "operation_status", "operation_type", "operation_name", "operation_description",
    "receipt_url", "sum", "created_at", "operation_date",
    "category_name", "account_type",
	}).AddRow(
		5, 1, nil, 3, 1, "done", "expense", name, "old desc",
		"url", sum, now, now, "Еда", "debit",
	)


	mock.ExpectQuery(regexp.QuoteMeta(`
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
			COALESCE(c.category_name, 'Без категории') as category_name, acc.account_type
		FROM updated_operation o
		LEFT JOIN category c ON o.category_id = c._id
		JOIN account acc ON acc._id = $5
	`)).
		WithArgs(&name, nil, &sum, 5, 1).
		WillReturnRows(rows)

	mock.ExpectExec(regexp.QuoteMeta(`
	UPDATE account SET balance = balance + $1 WHERE _id = $2
	`)).
		WithArgs(-1*sum, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	op, err := repo.UpdateOperation(context.Background(), req, 1, 5)
	require.NoError(t, err)
	require.Equal(t, sum, op.Sum)
	require.Equal(t, name, op.Name)
	require.Equal(t, "Еда", op.CategoryName)
}

func TestDeleteOperation_Success(t *testing.T) {
	repo, mock, close := setupOperationDB(t)
	defer close()

	now := time.Now()

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT sum, operation_type FROM operation WHERE _id = $1`)).
		WithArgs(7).
		WillReturnRows(sqlmock.NewRows([]string{"sum", "type"}).AddRow(50.0, "expense"))

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE account SET balance = balance - $1 WHERE _id = $2`)).
		WithArgs(-50.0, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rows := sqlmock.NewRows([]string{
		"_id", "account_from_id", "account_to_id", "category_id", "currency_id",
		"operation_status", "operation_type", "operation_name", "operation_description",
		"receipt_url", "sum", "created_at", "operation_date", "category_name",
	}).AddRow(
		7, 1, nil, 3, 1, "reverted", "expense", "Имя", "Описание",
		"url", 50.0, now, now, "Категория",
	)

	mock.ExpectQuery(regexp.QuoteMeta(`
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
	`)).
		WithArgs(7, 1).
		WillReturnRows(rows)

	mock.ExpectCommit()

	op, err := repo.DeleteOperation(context.Background(), 1, 7)
	require.NoError(t, err)
	require.Equal(t, 7, op.ID)
	require.Equal(t, "reverted", string(op.Status))
}
