package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	finmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/models"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewDBConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) GetAccountsByUser(ctx context.Context, userID int) ([]finmodels.Account, error) {
	query := `
		SELECT a._id, a.balance, a.account_type, a.currency_id, a.created_at, a.updated_at
		FROM account a
		JOIN sharings s ON a._id = s.account_id
		WHERE s.user_id = $1
		ORDER BY a.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts by user: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var accounts []finmodels.Account
	for rows.Next() {
		var account finmodels.Account
		err := rows.Scan(
			&account.ID,
			&account.Balance,
			&account.Type,
			&account.CurrencyID,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan account: %w", err)
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (r *PostgresRepository) GetAccountByID(ctx context.Context, userID, accountID int) (finmodels.Account, error) {
	query := `
		SELECT a._id, a.balance, a.account_type, a.currency_id, a.created_at, a.updated_at
		FROM account a
		JOIN sharings s ON a._id = s.account_id
		WHERE s.user_id = $1 AND a._id = $2
	`

	var account finmodels.Account
	err := r.db.QueryRowContext(ctx, query, userID, accountID).Scan(
		&account.ID,
		&account.Balance,
		&account.Type,
		&account.CurrencyID,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return finmodels.Account{}, ErrAccountNotFound
		}
		return finmodels.Account{}, fmt.Errorf("failed to get account by ID: %w", err)
	}

	return account, nil
}

func (r *PostgresRepository) CreateAccount(ctx context.Context, account finmodels.Account, userID int) (finmodels.Account, error) {
	query := `
		INSERT INTO account (balance, account_type, currency_id, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING _id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		account.Balance,
		account.Type,
		account.CurrencyID,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)

	if err != nil {
		return finmodels.Account{}, fmt.Errorf("failed to create account: %w", err)
	}

	if err := r.CreateUserAccount(ctx, userID, account.ID); err != nil {
		return finmodels.Account{}, fmt.Errorf("failed to link user to account: %w", err)
	}

	return account, nil
}

func (r *PostgresRepository) CreateUserAccount(ctx context.Context, userID, accountID int) error {
	query := `
		INSERT INTO sharings (account_id, user_id, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
	`

	_, err := r.db.ExecContext(ctx, query, accountID, userID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case ForeignKeyViolation:
				return ErrAccountNotFound
			default:
				return fmt.Errorf("failed to create user account: %w", err)
			}
		}
		return fmt.Errorf("failed to create user account: %w", err)
	}
	return nil
}

func (r *PostgresRepository) UpdateAccount(ctx context.Context, req finmodels.UpdateAccountRequest) (finmodels.Account, error) {
	query := `
		UPDATE account a
		SET 
			balance = COALESCE($1, a.balance),
			updated_at = NOW()
		FROM sharings s
		WHERE a._id = s.account_id AND s.user_id = $2 AND a._id = $3
		RETURNING a._id, a.balance, a.account_type, a.currency_id, a.created_at, a.updated_at
	`

	var acc finmodels.Account
	err := r.db.QueryRowContext(ctx, query,
		req.Balance,
		req.UserID,
		req.AccountID,
	).Scan(
		&acc.ID,
		&acc.Balance,
		&acc.Type,
		&acc.CurrencyID,
		&acc.CreatedAt,
		&acc.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return finmodels.Account{}, ErrAccountNotFound
		}
		return finmodels.Account{}, fmt.Errorf("failed to update account: %w", err)
	}

	return acc, nil
}

func (r *PostgresRepository) UpdateAccountBalance(ctx context.Context, accountID int, newBalance float64) error {
	query := `
		UPDATE account 
		SET balance = $1, updated_at = NOW()
		WHERE _id = $2
	`

	_, err := r.db.ExecContext(ctx, query, newBalance, accountID)
	if err != nil {
		return fmt.Errorf("failed to update account balance: %w", err)
	}

	return nil
}

func (r *PostgresRepository) DeleteAccount(ctx context.Context, userID, accID int) (finmodels.Account, error) {
	query := `
		DELETE FROM account
		WHERE _id = $1
		RETURNING _id, balance, account_type, currency_id, created_at, updated_at
	`

	var acc finmodels.Account
	err := r.db.QueryRowContext(ctx, query, accID).Scan(
		&acc.ID,
		&acc.Balance,
		&acc.Type,
		&acc.CurrencyID,
		&acc.CreatedAt,
		&acc.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return finmodels.Account{}, ErrAccountNotFound
		}
		return finmodels.Account{}, fmt.Errorf("failed to delete account: %w", err)
	}

	return acc, nil
}

var ErrAccountNotFound = errors.New("account not found")
