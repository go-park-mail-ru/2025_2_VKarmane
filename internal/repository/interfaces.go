package repository

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type Repository interface {
	Close() error

	// UserRepository методы
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
	GetUserByID(ctx context.Context, id int) (models.User, error)
	UpdateUser(ctx context.Context, user models.User) error

	// AccountRepository методы
	GetAccountsByUser(ctx context.Context, userID int) ([]models.Account, error)

	// BudgetRepository методы
	GetBudgetsByUser(ctx context.Context, userID int) ([]models.Budget, error)

	// OperationRepository методы
	GetOperationsByAccount(ctx context.Context, accountID int) ([]models.Operation, error)
	GetOperationsByUser(ctx context.Context, userID int) ([]models.Operation, error)

	// CategoryRepository методы
	GetCategoriesByUser(ctx context.Context, userID int) ([]models.Category, error)
	GetCategoryByID(ctx context.Context, userID, categoryID int) (models.Category, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
	GetUserByID(ctx context.Context, id int) (models.User, error)
	UpdateUser(ctx context.Context, user models.User) error
}

type AccountRepository interface {
	GetAccountsByUser(ctx context.Context, userID int) ([]models.Account, error)
}

type BudgetRepository interface {
	GetBudgetsByUser(ctx context.Context, userID int) ([]models.Budget, error)
}

type OperationRepository interface {
	GetOperationsByAccount(ctx context.Context, accountID int) ([]models.Operation, error)
	GetOperationsByUser(ctx context.Context, userID int) ([]models.Operation, error)
}

type CategoryRepository interface {
	GetCategoriesByUser(ctx context.Context, userID int) ([]models.Category, error)
	GetCategoryByID(ctx context.Context, userID, categoryID int) (models.Category, error)
}
