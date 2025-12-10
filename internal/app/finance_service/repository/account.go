package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/errors"
	finmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type PostgresRepository struct {
	db *sql.DB
}

type ESRepository struct {
	es *elasticsearch.Client
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

func NewESConnection(host, port string) (*elasticsearch.Client, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("http://%s:%s", host, port),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open ElasticSearch: %w", err)
	}
	return es, nil
}

func NewESRepository(es *elasticsearch.Client) *ESRepository {
	return &ESRepository{
		es: es,
	}
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) GetAccountsByUser(ctx context.Context, userID int) ([]finmodels.Account, error) {
	query := `
			SELECT a._id, a.balance, a.account_type, a.currency_id, 
		a.created_at, a.updated_at, a.account_name, a.account_description
		FROM account a
		JOIN sharings s ON a._id = s.account_id
		JOIN 
		WHERE s.user_id = $1

		ORDER BY a.balance DESC, a.created_at DESC
		LIMIT 3;
	`
	rows, err := r.db.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, MapPgAccountError(err)
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
			&account.Name,
			&account.Description,
		)
		if err != nil {
			return nil, MapPgAccountError(err)
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (r *PostgresRepository) GetAccountByID(ctx context.Context, userID, accountID int) (finmodels.Account, error) {
	query := `
		SELECT a._id, a.balance, a.account_type, a.currency_id, a.created_at, a.updated_at, a.account_name, a.account_description
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
		&account.Name,
		&account.Description,
	)

	if err != nil {
		return finmodels.Account{}, MapPgAccountError(err)
	}

	return account, nil
}

func (r *PostgresRepository) CreateAccount(ctx context.Context, account finmodels.Account, userID int) (finmodels.Account, error) {
	query := `
		INSERT INTO account (account_name, account_description, balance, account_type, currency_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING _id, created_at, updated_at, account_name, account_description
	`
	err := r.db.QueryRowContext(ctx, query,
		account.Name,
		account.Description,
		account.Balance,
		account.Type,
		account.CurrencyID,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt, &account.Name, &account.Description)

	if err != nil {
		return finmodels.Account{}, MapPgAccountError(err)
	}

	if err := r.CreateUserAccount(ctx, userID, account.ID); err != nil {
		return finmodels.Account{}, MapPgAccountError(err)
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
		return MapPgAccountError(err)
	}
	return nil
}

func (r *PostgresRepository) UpdateAccount(ctx context.Context, req finmodels.UpdateAccountRequest) (finmodels.Account, error) {
	query := `
		UPDATE account a
		SET 
			balance = COALESCE($1, a.balance),
			account_name = COALESCE($2, a.account_name),
			account_description = COALESCE($3, a.account_description),
			updated_at = NOW()
		FROM sharings s
		WHERE a._id = s.account_id AND s.user_id = $4 AND a._id = $5
		RETURNING a._id, a.account_name, a.account_description, a.balance, a.account_type, a.currency_id, a.created_at, a.updated_at
	`

	var acc finmodels.Account
	err := r.db.QueryRowContext(ctx, query,
		req.Balance,
		req.Name,
		req.Description,
		req.UserID,
		req.AccountID,
	).Scan(
		&acc.ID,
		&acc.Name,
		&acc.Description,
		&acc.Balance,
		&acc.Type,
		&acc.CurrencyID,
		&acc.CreatedAt,
		&acc.UpdatedAt,
	)

	if err != nil {
		return finmodels.Account{}, MapPgAccountError(err)
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
		return MapPgAccountError(err)
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
		return finmodels.Account{}, MapPgAccountError(err)
	}

	return acc, nil
}

func (r *PostgresRepository) AddUserToAccount(ctx context.Context, userID, accountID int) (finmodels.SharingAccount, error) {
	var accountType string

	err := r.db.QueryRowContext(ctx, `SELECT account_type FROM account where _id = $1`, accountID).Scan(&accountType)
	if err != nil {
		return finmodels.SharingAccount{}, MapPgAccountError(err)
	}

	if accountType == models.PrivateAccount {
		return finmodels.SharingAccount{}, errors.ErrPrivateAccount
	}

	query := `
		INSERT INTO sharings
		(account_id, user_id, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING _id, account_id, user_id, created_at
	`
	sh := finmodels.SharingAccount{}

	err = r.db.QueryRowContext(ctx, query, accountID, userID).
		Scan(&sh.ID, &sh.AccountID, &sh.UserID, &sh.CreatedAt)

	if err != nil {
		return finmodels.SharingAccount{}, MapPgAccountError(err)
	}
	return sh, nil

}
