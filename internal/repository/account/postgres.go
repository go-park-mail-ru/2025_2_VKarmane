package account

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	postgreserrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/errors"
	"github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) GetAccountsByUser(ctx context.Context, userID int) ([]models.Account, error) {
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

	var accounts []models.Account
	for rows.Next() {
		var account AccountDB
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
		accounts = append(accounts, AccountDBToModel(account))
	}

	return accounts, nil
}

// func (r *PostgresRepository) CreateAccount(ctx context.Context, account AccountDB) (int, error) {
// 	query := `
// 		INSERT INTO account (balance, account_type, currency_id, created_at, updated_at)
// 		VALUES ($1, $2, $3, $4, $5)
// 		RETURNING _id
// 	`

// 	var id int
// 	err := r.db.QueryRowContext(ctx, query,
// 		account.Balance,
// 		account.Type,
// 		account.CurrencyID,
// 		account.CreatedAt,
// 		account.UpdatedAt,
// 	).Scan(&id)

// 	if err != nil {
// 		return 0, fmt.Errorf("failed to create account: %w", err)
// 	}

// 	return id, nil
// }

// func (r *PostgresRepository) CreateUserAccount(ctx context.Context, userAccount UserAccountDB) error {
// 	query := `
// 		INSERT INTO sharings (account_id, user_id, created_at, updated_at)
// 		VALUES ($1, $2, $3, $4)
// 	`

// 	_, err := r.db.ExecContext(ctx, query,
// 		userAccount.AccountID,
// 		userAccount.UserID,
// 		userAccount.CreatedAt,
// 		userAccount.UpdatedAt,
// 	)

// 	if err != nil {
// 		return fmt.Errorf("failed to create user account: %w", err)
// 	}

// 	return nil
// }

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

func (r *PostgresRepository) CreateAccount(ctx context.Context, account models.Account, userID int) (models.Account, error) {
	query := `
		INSERT INTO account (balance, account_type, currency_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING _id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		account.Balance,
		account.Type,
		account.CurrencyID,
		account.CreatedAt,
		account.UpdatedAt,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)

	if err != nil {
		return models.Account{}, fmt.Errorf("failed to create account: %w", err)
	}

	if err := r.CreateUserAccount(ctx, models.UserAccount{
		UserID:    userID,
		AccountID: account.ID,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}); err != nil {
		return models.Account{}, fmt.Errorf("failed to link user to account: %w", err)
	}

	return account, nil
}

func (r *PostgresRepository) CreateUserAccount(ctx context.Context, ua models.UserAccount) error {
	query := `
		INSERT INTO sharings (account_id, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query,
		ua.AccountID,
		ua.UserID,
		ua.CreatedAt,
		ua.UpdatedAt,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case postgreserrors.ForeignKeyViolation:
				return ErrAccountNotFound
			default:
				return fmt.Errorf("failed to create user account: %w", err)
			}
		}
		return fmt.Errorf("failed to create user account: %w", err)
	}
	return nil
}

func (r *PostgresRepository) UpdateAccount(ctx context.Context, req models.UpdateAccountRequest, userID, accID int) (models.Account, error) {
	query := `
		UPDATE account a
		SET 
			balance = COALESCE($1, a.balance),
			updated_at = NOW()
		FROM sharings s
		WHERE a._id = s.account_id AND s.user_id = $2 AND a._id = $3
		RETURNING a._id, a.balance, a.account_type, a.currency_id, a.created_at, a.updated_at
	`

	var acc models.Account
	err := r.db.QueryRowContext(ctx, query,
		req.Balance,
		userID,
		accID,
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
			return models.Account{}, ErrAccountNotFound
		}
		return models.Account{}, fmt.Errorf("failed to update account: %w", err)
	}

	return acc, nil
}

func (r *PostgresRepository) DeleteAccount(ctx context.Context, userID, accID int) (models.Account, error) {
	query := `
		DELETE FROM account
		WHERE _id = $1
		RETURNING _id, balance, account_type, currency_id, created_at, updated_at
	`

	var acc models.Account
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
			return models.Account{}, ErrAccountNotFound
		}
		return models.Account{}, fmt.Errorf("failed to delete account: %w", err)
	}

	return acc, nil
}
