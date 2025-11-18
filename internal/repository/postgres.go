package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/account"
	categoryrepo "github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/category"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/operation"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
	DB            *sql.DB
	AccountRepo   *account.PostgresRepository
	OperationRepo *operation.PostgresRepository
	CategoryRepo  *categoryrepo.PostgresRepository
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

	store.AccountRepo = account.NewPostgresRepository(db)
	store.OperationRepo = operation.NewPostgresRepository(db)
	store.CategoryRepo = categoryrepo.NewPostgresRepository(db)

	return store, nil
}

func (s *PostgresStore) Close() error {
	if err := s.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}
	return nil
}

// Реализация интерфейсов репозиториев


// AccountRepository
func (s *PostgresStore) GetAccountsByUser(ctx context.Context, userID int) ([]models.Account, error) {
	return s.AccountRepo.GetAccountsByUser(ctx, userID)
}

func (s *PostgresStore) CreateAccount(ctx context.Context, account models.Account, userID int) (models.Account, error) {
	return s.AccountRepo.CreateAccount(ctx, account, userID)
}

func (s *PostgresStore) UpdateAccount(ctx context.Context, req models.UpdateAccountRequest, userID, accID int) (models.Account, error) {
	return s.AccountRepo.UpdateAccount(ctx, req, userID, accID)
}
func (s *PostgresStore) DeleteAccount(ctx context.Context, userID, accID int) (models.Account, error) {
	return s.AccountRepo.DeleteAccount(ctx, userID, accID)
}


// OperationRepository
func (s *PostgresStore) GetOperationsByAccount(ctx context.Context, accountID int) ([]models.OperationInList, error) {
	return s.OperationRepo.GetOperationsByAccount(ctx, accountID)
}

func (s *PostgresStore) GetOperationsByUser(ctx context.Context, userID int) ([]models.Operation, error) {
	return s.OperationRepo.GetOperationsByUser(ctx, userID)
}

func (s *PostgresStore) GetOperationByID(ctx context.Context, accID int, opID int) (models.Operation, error) {
	return s.OperationRepo.GetOperationByID(ctx, accID, opID)
}

func (s *PostgresStore) CreateOperation(ctx context.Context, op models.Operation) (models.Operation, error) {
	return s.OperationRepo.CreateOperation(ctx, op)
}

func (s *PostgresStore) UpdateOperation(ctx context.Context, req models.UpdateOperationRequest, accID int, opID int) (models.Operation, error) {
	return s.OperationRepo.UpdateOperation(ctx, req, accID, opID)
}

func (s *PostgresStore) DeleteOperation(ctx context.Context, accID int, opID int) (models.Operation, error) {
	return s.OperationRepo.DeleteOperation(ctx, accID, opID)
}

// CategoryRepository
func (s *PostgresStore) GetCategoriesByUser(ctx context.Context, userID int) ([]models.Category, error) {
	categoriesDB, err := s.CategoryRepo.GetCategoriesByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories by user: %w", err)
	}

	var categories []models.Category
	for _, categoryDB := range categoriesDB {
		categories = append(categories, categoryrepo.CategoryDBToModel(categoryDB))
	}

	return categories, nil
}

func (s *PostgresStore) GetCategoryByID(ctx context.Context, userID, categoryID int) (models.Category, error) {
	categoryDB, err := s.CategoryRepo.GetCategoryByID(ctx, userID, categoryID)
	if err != nil {
		return models.Category{}, fmt.Errorf("failed to get category by ID: %w", err)
	}

	return categoryrepo.CategoryDBToModel(categoryDB), nil
}

func (s *PostgresStore) CreateCategory(ctx context.Context, category models.Category) (models.Category, error) {
	categoryDB := categoryrepo.CategoryModelToDB(category)
	id, err := s.CategoryRepo.CreateCategory(ctx, categoryDB)
	if err != nil {
		return models.Category{}, fmt.Errorf("failed to create category: %w", err)
	}
	category.ID = id
	return category, nil
}

func (s *PostgresStore) UpdateCategory(ctx context.Context, category models.Category) error {
	categoryDB := categoryrepo.CategoryModelToDB(category)
	if err := s.CategoryRepo.UpdateCategory(ctx, categoryDB); err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	return nil
}

func (s *PostgresStore) DeleteCategory(ctx context.Context, userID, categoryID int) error {
	if err := s.CategoryRepo.DeleteCategory(ctx, userID, categoryID); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}

func (s *PostgresStore) GetCategoryStats(ctx context.Context, userID, categoryID int) (int, error) {
	stats, err := s.CategoryRepo.GetCategoryStats(ctx, userID, categoryID)
	if err != nil {
		return 0, fmt.Errorf("failed to get category stats: %w", err)
	}
	return stats, nil
}
