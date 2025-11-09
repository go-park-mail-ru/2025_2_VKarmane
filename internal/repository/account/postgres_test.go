package account

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresRepository_GetAccountsByUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	repo := NewPostgresRepository(db)

	userID := 1
	now := time.Now()
	rows := sqlmock.NewRows([]string{"_id", "balance", "account_type", "currency_id", "created_at", "updated_at"}).
		AddRow(1, 1000.50, "card", 1, now, now).
		AddRow(2, 500.25, "cash", 1, now, now)

	mock.ExpectPrepare(`SELECT.*FROM account`).
		ExpectQuery().
		WithArgs(userID).
		WillReturnRows(rows)

	accounts, err := repo.GetAccountsByUser(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, accounts, 2)
	assert.Equal(t, 1, accounts[0].ID)
	assert.Equal(t, 1000.50, accounts[0].Balance)
	assert.Equal(t, "card", accounts[0].Type)
	assert.Equal(t, 2, accounts[1].ID)
	assert.Equal(t, 500.25, accounts[1].Balance)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetAccountsByUser_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	repo := NewPostgresRepository(db)

	userID := 1
	rows := sqlmock.NewRows([]string{"_id", "balance", "account_type", "currency_id", "created_at", "updated_at"})

	mock.ExpectPrepare(`SELECT.*FROM account`).
		ExpectQuery().
		WithArgs(userID).
		WillReturnRows(rows)

	accounts, err := repo.GetAccountsByUser(context.Background(), userID)
	assert.NoError(t, err)
	assert.Empty(t, accounts)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetAccountsByUser_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	repo := NewPostgresRepository(db)

	userID := 1
	mock.ExpectPrepare(`SELECT.*FROM account`).
		ExpectQuery().
		WithArgs(userID).
		WillReturnError(sql.ErrConnDone)

	accounts, err := repo.GetAccountsByUser(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, accounts)
	assert.Contains(t, err.Error(), "failed to get accounts by user")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetAccountsByUser_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	repo := NewPostgresRepository(db)

	userID := 1
	rows := sqlmock.NewRows([]string{"_id", "balance", "account_type", "currency_id", "created_at", "updated_at"}).
		AddRow("invalid", "invalid", "invalid", "invalid", "invalid", "invalid")

	mock.ExpectPrepare(`SELECT.*FROM account`).
		ExpectQuery().
		WithArgs(userID).
		WillReturnRows(rows)

	accounts, err := repo.GetAccountsByUser(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, accounts)
	assert.Contains(t, err.Error(), "failed to scan account")
	assert.NoError(t, mock.ExpectationsWereMet())
}
