package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/account"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/budget"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/operation"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/user"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
	DB            *sql.DB
	UserRepo      *user.PostgresRepository
	AccountRepo   *account.PostgresRepository
	BudgetRepo    *budget.PostgresRepository
	OperationRepo *operation.PostgresRepository
}

func NewPostgresStore(dsn string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	store := &PostgresStore{
		DB: db,
	}

	store.UserRepo = user.NewPostgresRepository(db)
	store.AccountRepo = account.NewPostgresRepository(db)
	store.BudgetRepo = budget.NewPostgresRepository(db)
	store.OperationRepo = operation.NewPostgresRepository(db)

	return store, nil
}

func (s *PostgresStore) Close() error {
	return s.DB.Close()
}

// Реализация интерфейсов репозиториев

// UserRepository
func (s *PostgresStore) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	return s.UserRepo.CreateUserModel(ctx, user)
}

func (s *PostgresStore) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	return s.UserRepo.GetUserByLoginModel(ctx, login)
}

func (s *PostgresStore) GetUserByID(ctx context.Context, id int) (models.User, error) {
	return s.UserRepo.GetUserByIDModel(ctx, id)
}

// AccountRepository
func (s *PostgresStore) GetAccountsByUser(ctx context.Context, userID int) ([]models.Account, error) {
	return s.AccountRepo.GetAccountsByUser(ctx, userID)
}

// BudgetRepository
func (s *PostgresStore) GetBudgetsByUser(ctx context.Context, userID int) ([]models.Budget, error) {
	return s.BudgetRepo.GetBudgetsByUser(ctx, userID)
}

// OperationRepository
func (s *PostgresStore) GetOperationsByAccount(ctx context.Context, accountID int) ([]models.Operation, error) {
	return s.OperationRepo.GetOperationsByAccount(ctx, accountID)
}

func (s *PostgresStore) GetOperationsByUser(ctx context.Context, userID int) ([]models.Operation, error) {
	return s.OperationRepo.GetOperationsByUser(ctx, userID)
}
