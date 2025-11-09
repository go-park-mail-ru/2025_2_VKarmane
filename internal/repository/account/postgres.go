package account

import (
	"context"
	"database/sql"
	"fmt"

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

func (r *PostgresRepository) GetAccountsByUser(ctx context.Context, userID int) ([]models.Account, error) {
	query := `
		SELECT a._id, a.balance, a.account_type, a.currency_id, a.created_at, a.updated_at
		FROM account a
		JOIN sharings s ON a._id = s.account_id
		WHERE s.user_id = $1
		ORDER BY a.created_at DESC
	`
	stmt, err := r.db.PrepareContext(ctx, query)
		if err != nil {
			return []models.Account{}, fmt.Errorf("failed to prepare stmt: %w", err)
		}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userID)
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

func (r *PostgresRepository) CreateAccount(ctx context.Context, account AccountDB) (int, error) {
	query := `
		INSERT INTO account (balance, account_type, currency_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING _id
	`
	stmt, err := r.db.PrepareContext(ctx, query)
		if err != nil {
			return 0, fmt.Errorf("failed to prepare stmt: %w", err)
		}
	defer stmt.Close()

	var id int
	err = stmt.QueryRowContext(ctx,
		account.Balance,
		account.Type,
		account.CurrencyID,
		account.CreatedAt,
		account.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create account: %w", err)
	}

	return id, nil
}

func (r *PostgresRepository) CreateUserAccount(ctx context.Context, userAccount UserAccountDB) error {
	query := `
		INSERT INTO sharings (account_id, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`
	stmt, err := r.db.PrepareContext(ctx, query)
		if err != nil {
			return  fmt.Errorf("failed to prepare stmt: %w", err)
		}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		userAccount.AccountID,
		userAccount.UserID,
		userAccount.CreatedAt,
		userAccount.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user account: %w", err)
	}

	return nil
}

func (r *PostgresRepository) UpdateAccountBalance(ctx context.Context, accountID int, newBalance float64) error {
	query := `
		UPDATE account 
		SET balance = $1, updated_at = NOW()
		WHERE _id = $2
	`
	stmt, err := r.db.PrepareContext(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to prepare stmt: %w", err)
		}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, newBalance, accountID)
	if err != nil {
		return fmt.Errorf("failed to update account balance: %w", err)
	}

	return nil
}
