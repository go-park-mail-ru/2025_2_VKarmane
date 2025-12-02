package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	serviceerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/errors"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/models"
	"github.com/stretchr/testify/require"
)

func setupDB(t *testing.T) (*PostgresRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewPostgresRepository(db)
	return repo, mock, func() { _ = db.Close() }
}

func TestGetAccountByID_Success(t *testing.T) {
	repo, mock, close := setupDB(t)
	defer close()

	accountID := 1
	userID := 2
	created := time.Now()
	updated := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`
    SELECT a._id, a.balance, a.account_type, a.currency_id, a.created_at, a.updated_at,
           a.account_name, a.account_description
    FROM account a
    JOIN sharings s ON a._id = s.account_id
    WHERE s.user_id = $1 AND a._id = $2
`)).
		WithArgs(userID, accountID).
		WillReturnRows(sqlmock.NewRows([]string{
			"_id", "balance", "account_type", "currency_id",
			"created_at", "updated_at", "account_name", "account_description",
		}).AddRow(accountID, 100.0, "cash", 1, created, updated, "name", "desc"))

	acc, err := repo.GetAccountByID(context.Background(), userID, accountID)
	require.NoError(t, err)
	require.Equal(t, accountID, acc.ID)
	require.Equal(t, 100.0, acc.Balance)
}

func TestGetAccountByID_NotFound(t *testing.T) {
	repo, mock, close := setupDB(t)
	defer close()

	mock.ExpectQuery(regexp.QuoteMeta(`
    SELECT a._id, a.balance, a.account_type, a.currency_id, a.created_at, a.updated_at,
           a.account_name, a.account_description
    FROM account a
    JOIN sharings s ON a._id = s.account_id
    WHERE s.user_id = $1 AND a._id = $2
	`)).
		WithArgs(1, 2).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetAccountByID(context.Background(), 1, 2)
	require.ErrorIs(t, err, serviceerrors.ErrAccountNotFound)
}

func TestGetAccountsByUser(t *testing.T) {
	repo, mock, close := setupDB(t)
	defer close()

	userID := 1
	created := time.Now()
	updated := time.Now()

	rows := sqlmock.NewRows([]string{
		"_id", "balance", "account_type", "currency_id",
		"created_at", "updated_at", "account_name", "account_description",
	}).
		AddRow(1, 100.0, "cash", 1, created, updated, "Cash Wallet", "My cash").
		AddRow(2, 200.0, "card", 2, created, updated, "Bank Card", "My card")

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT a._id, a.balance, a.account_type, a.currency_id, 
		       a.created_at, a.updated_at, a.account_name, a.account_description
		FROM account a
		JOIN sharings s ON a._id = s.account_id
		WHERE s.user_id = $1
		ORDER BY a.created_at DESC
	`)).
		WithArgs(userID).
		WillReturnRows(rows)

	accs, err := repo.GetAccountsByUser(context.Background(), userID)
	require.NoError(t, err)
	require.Len(t, accs, 2)
	require.Equal(t, 100.0, accs[0].Balance)
	require.Equal(t, 200.0, accs[1].Balance)
}

func TestCreateAccount_Success(t *testing.T) {
	repo, mock, close := setupDB(t)
	defer close()

	userID := 1
	account := models.Account{
		Balance:    100,
		Type:       "cash",
		CurrencyID: 1,
		Name:       "My Account",
	}
	created := time.Now()
	updated := time.Now()

	// mock insert account
	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO account (account_name, account_description, balance, account_type, currency_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING _id, created_at, updated_at, account_name, account_description
	`)).
		WithArgs(
			account.Name,
			account.Description,
			account.Balance,
			account.Type,
			account.CurrencyID,
		).
		WillReturnRows(sqlmock.NewRows([]string{
			"_id", "created_at", "updated_at", "account_name", "account_description",
		}).AddRow(5, created, updated, account.Name, account.Description))

	// mock insert into sharings
	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO sharings (account_id, user_id, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
	`)).
		WithArgs(5, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	acc, err := repo.CreateAccount(context.Background(), account, userID)
	require.NoError(t, err)
	require.Equal(t, 5, acc.ID)
	require.Equal(t, 100.0, acc.Balance)
	require.Equal(t, account.Name, acc.Name)
	require.Equal(t, account.Description, acc.Description)
}

func TestUpdateAccount_Success(t *testing.T) {
	repo, mock, close := setupDB(t)
	defer close()
	balance := 500.
	name := "beba"
	desrc := "descr"

	req := models.UpdateAccountRequest{UserID: 1, AccountID: 2, Balance: &balance, Name: &name, Description: &desrc}
	created := time.Now()
	updated := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`
    UPDATE account a
    SET 
        balance = COALESCE($1, a.balance),
        account_name = COALESCE($2, a.account_name),
        account_description = COALESCE($3, a.account_description),
        updated_at = NOW()
    FROM sharings s
    WHERE a._id = s.account_id AND s.user_id = $4 AND a._id = $5
    RETURNING a._id, a.account_name, a.account_description, a.balance, a.account_type, a.currency_id, a.created_at, a.updated_at
`)).
		WithArgs(req.Balance, req.Name, req.Description, req.UserID, req.AccountID).
		WillReturnRows(sqlmock.NewRows([]string{
			"_id", "account_name", "account_description", "balance", "account_type", "currency_id", "created_at", "updated_at",
		}).AddRow(req.AccountID, "", "", req.Balance, "cash", 1, created, updated))

	acc, err := repo.UpdateAccount(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, 500.0, acc.Balance)
}

func TestUpdateAccountBalance(t *testing.T) {
	repo, mock, close := setupDB(t)
	defer close()

	accountID := 2
	newBalance := 300.0

	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE account 
		SET balance = $1, updated_at = NOW()
		WHERE _id = $2
	`)).WithArgs(newBalance, accountID).WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.UpdateAccountBalance(context.Background(), accountID, newBalance)
	require.NoError(t, err)
}

func TestDeleteAccount(t *testing.T) {
	repo, mock, close := setupDB(t)
	defer close()

	accID := 3
	created := time.Now()
	updated := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`
		DELETE FROM account
		WHERE _id = $1
		RETURNING _id, balance, account_type, currency_id, created_at, updated_at
	`)).WithArgs(accID).WillReturnRows(sqlmock.NewRows([]string{"_id", "balance", "account_type", "currency_id", "created_at", "updated_at"}).
		AddRow(accID, 400.0, "card", 1, created, updated))

	acc, err := repo.DeleteAccount(context.Background(), 1, accID)
	require.NoError(t, err)
	require.Equal(t, 400.0, acc.Balance)
}
