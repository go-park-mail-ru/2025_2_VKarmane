package account

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/lib/pq"
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

	mock.ExpectQuery(`SELECT.*FROM account`).
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

	mock.ExpectQuery(`SELECT.*FROM account`).
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
	mock.ExpectQuery(`SELECT.*FROM account`).
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

	mock.ExpectQuery(`SELECT.*FROM account`).
		WithArgs(userID).
		WillReturnRows(rows)

	accounts, err := repo.GetAccountsByUser(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, accounts)
	assert.Contains(t, err.Error(), "failed to scan account")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_CreateAccount(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		account       models.Account
		userID        int
		mockSetup     func(sqlmock.Sqlmock)
		expectError   bool
		errorContains string
	}{
		{
			name: "success",
			account: models.Account{
				Balance:    100.5,
				Type:       "card",
				CurrencyID: 1,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO account`).
					WithArgs(100.5, "card", 1, now, now).
					WillReturnRows(sqlmock.NewRows([]string{"_id", "created_at", "updated_at"}).
						AddRow(10, now, now))

				mock.ExpectExec(`INSERT INTO sharings`).
					WithArgs(10, 1, now, now).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectError: false,
		},
		{
			name: "error on account insert",
			account: models.Account{
				Balance:    200,
				Type:       "cash",
				CurrencyID: 2,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO account`).
					WillReturnError(sql.ErrConnDone)
			},
			expectError:   true,
			errorContains: "failed to create account",
		},
		{
			name: "error on linking user",
			account: models.Account{
				Balance:    300,
				Type:       "card",
				CurrencyID: 3,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO account`).
					WillReturnRows(sqlmock.NewRows([]string{"_id", "created_at", "updated_at"}).
						AddRow(15, now, now))
				mock.ExpectExec(`INSERT INTO sharings`).
					WillReturnError(errors.New("insert failed"))
			},
			expectError:   true,
			errorContains: "failed to link user to account",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := NewPostgresRepository(db)
			tt.mockSetup(mock)

			_, err = repo.CreateAccount(context.Background(), tt.account, tt.userID)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresRepository_CreateUserAccount(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		userAccount   models.UserAccount
		mockSetup     func(sqlmock.Sqlmock)
		expectError   bool
		errorContains string
	}{
		{
			name: "success",
			userAccount: models.UserAccount{
				AccountID: 1, UserID: 2, CreatedAt: now, UpdatedAt: now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO sharings`).
					WithArgs(1, 2, now, now).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectError: false,
		},
		{
			name: "foreign key violation",
			userAccount: models.UserAccount{
				AccountID: 999, UserID: 2, CreatedAt: now, UpdatedAt: now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				pqErr := &pq.Error{Code: "23503"} // foreign key violation
				mock.ExpectExec(`INSERT INTO sharings`).
					WillReturnError(pqErr)
			},
			expectError:   true,
			errorContains: "account not found",
		},
		{
			name: "generic error",
			userAccount: models.UserAccount{
				AccountID: 1, UserID: 2, CreatedAt: now, UpdatedAt: now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO sharings`).
					WillReturnError(errors.New("db insert error"))
			},
			expectError:   true,
			errorContains: "failed to create user account",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			repo := NewPostgresRepository(db)
			tt.mockSetup(mock)

			err := repo.CreateUserAccount(context.Background(), tt.userAccount)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresRepository_UpdateAccountBalance(t *testing.T) {
	tests := []struct {
		name          string
		accountID     int
		newBalance    float64
		mockSetup     func(sqlmock.Sqlmock)
		expectError   bool
		errorContains string
	}{
		{
			name:       "success",
			accountID:  1,
			newBalance: 123.45,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE account`).
					WithArgs(123.45, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectError: false,
		},
		{
			name:       "exec error",
			accountID:  2,
			newBalance: 200,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE account`).
					WillReturnError(sql.ErrConnDone)
			},
			expectError:   true,
			errorContains: "failed to update account balance",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			repo := NewPostgresRepository(db)
			tt.mockSetup(mock)

			err := repo.UpdateAccountBalance(context.Background(), tt.accountID, tt.newBalance)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresRepository_UpdateAccount(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name          string
		req           models.UpdateAccountRequest
		userID        int
		accID         int
		mockSetup     func(sqlmock.Sqlmock)
		expectError   bool
		errorContains string
	}{
		{
			name:   "success",
			req:    models.UpdateAccountRequest{Balance: 300},
			userID: 1, accID: 2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"_id", "balance", "account_type", "currency_id", "created_at", "updated_at"}).
					AddRow(2, 300.0, "cash", 1, now, now)
				mock.ExpectQuery(`UPDATE account`).
					WillReturnRows(rows)
			},
			expectError: false,
		},
		{
			name:   "no rows found",
			req:    models.UpdateAccountRequest{Balance: 400},
			userID: 1, accID: 3,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`UPDATE account`).
					WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			errorContains: "account not found",
		},
		{
			name:   "query error",
			req:    models.UpdateAccountRequest{Balance: 500},
			userID: 1, accID: 3,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`UPDATE account`).
					WillReturnError(errors.New("db failed"))
			},
			expectError:   true,
			errorContains: "failed to update account",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			repo := NewPostgresRepository(db)
			tt.mockSetup(mock)

			_, err := repo.UpdateAccount(context.Background(), tt.req, tt.userID, tt.accID)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresRepository_DeleteAccount(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name          string
		userID        int
		accID         int
		mockSetup     func(sqlmock.Sqlmock)
		expectError   bool
		errorContains string
	}{
		{
			name:   "success",
			userID: 1, accID: 2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"_id", "balance", "account_type", "currency_id", "created_at", "updated_at"}).
					AddRow(2, 100, "card", 1, now, now)
				mock.ExpectQuery(`DELETE FROM account`).
					WithArgs(2).
					WillReturnRows(rows)
			},
			expectError: false,
		},
		{
			name:   "no rows found",
			userID: 1, accID: 3,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`DELETE FROM account`).
					WithArgs(3).
					WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			errorContains: "account not found",
		},
		{
			name:   "query error",
			userID: 1, accID: 4,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`DELETE FROM account`).
					WithArgs(4).
					WillReturnError(errors.New("delete failed"))
			},
			expectError:   true,
			errorContains: "failed to delete account",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			repo := NewPostgresRepository(db)
			tt.mockSetup(mock)

			_, err := repo.DeleteAccount(context.Background(), tt.userID, tt.accID)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
