package service

import (
	"context"

	finmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/models"
)

type FinanceRepository interface {
	// Account methods
	GetAccountsByUser(ctx context.Context, userID int) ([]finmodels.Account, error)
	GetAccountByID(ctx context.Context, userID, accountID int) (finmodels.Account, error)
	CreateAccount(ctx context.Context, account finmodels.Account, userID int) (finmodels.Account, error)
	UpdateAccount(ctx context.Context, req finmodels.UpdateAccountRequest) (finmodels.Account, error)
	DeleteAccount(ctx context.Context, userID, accID int) (finmodels.Account, error)
	UpdateAccountBalance(ctx context.Context, accountID int, newBalance float64) error
	AddUserToAccount(ctx context.Context, userID, accountID int) (finmodels.SharingAccount, error)

	// Operation methods
	GetOperationsByAccount(ctx context.Context, accountID int) ([]finmodels.OperationInList, error)
	GetOperationByID(ctx context.Context, accID int, opID int) (finmodels.Operation, error)
	CreateOperation(ctx context.Context, op finmodels.Operation) (finmodels.Operation, error)
	UpdateOperation(ctx context.Context, req finmodels.UpdateOperationRequest, accID int, opID int) (finmodels.Operation, error)
	DeleteOperation(ctx context.Context, accID int, opID int) (finmodels.Operation, error)

	// Category methods
	CreateCategory(ctx context.Context, category finmodels.Category) (finmodels.Category, error)
	GetCategoriesByUser(ctx context.Context, userID int) ([]finmodels.Category, error)
	GetCategoriesWithStatsByUser(ctx context.Context, userID int) ([]finmodels.CategoryWithStats, error)
	GetCategoryByID(ctx context.Context, userID, categoryID int) (finmodels.Category, error)
	GetCategoryByName(ctx context.Context, userID int, categoryName string) (finmodels.Category, error)
	UpdateCategory(ctx context.Context, category finmodels.Category) error
	DeleteCategory(ctx context.Context, userID, categoryID int) error
	GetCategoryStats(ctx context.Context, userID, categoryID int) (int, error)
}
