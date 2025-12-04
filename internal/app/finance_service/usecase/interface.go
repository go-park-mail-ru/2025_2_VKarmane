package usecase

import (
	"context"

	finmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/models"
	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
)

type FinanceService interface {
	// Account methods
	GetAccountsByUser(ctx context.Context, userID int) (*finpb.ListAccountsResponse, error)
	GetAccountByID(ctx context.Context, userID, accountID int) (*finpb.Account, error)
	CreateAccount(ctx context.Context, req finmodels.CreateAccountRequest) (*finpb.Account, error)
	UpdateAccount(ctx context.Context, req finmodels.UpdateAccountRequest) (*finpb.Account, error)
	DeleteAccount(ctx context.Context, userID, accountID int) (*finpb.Account, error)
	AddUserToAccount(ctx context.Context, userID, accountID int) (*finpb.SharingsResponse, error)

	// Operation methods
	// GetOperationsByAccount(ctx context.Context, accountID, categoryID int, opName string) (*finpb.ListOperationsResponse, error)
	GetOperationsByAccount(ctx context.Context, req []byte) (*finpb.ListOperationsResponse, error)
	GetOperationByID(ctx context.Context, accID, opID int) (*finpb.Operation, error)
	CreateOperation(ctx context.Context, req finmodels.CreateOperationRequest, accountID int) (*finpb.Operation, error)
	UpdateOperation(ctx context.Context, req finmodels.UpdateOperationRequest) (*finpb.Operation, error)
	DeleteOperation(ctx context.Context, accID, opID int) (*finpb.Operation, error)

	// Category methods
	CreateCategory(ctx context.Context, req finmodels.CreateCategoryRequest) (*finpb.Category, error)
	GetCategoriesByUser(ctx context.Context, userID int) (*finpb.ListCategoriesResponse, error)
	GetCategoriesWithStatsByUser(ctx context.Context, userID int) (*finpb.ListCategoriesWithStatsResponse, error)
	GetCategoryByID(ctx context.Context, userID, categoryID int) (*finpb.CategoryWithStats, error)
	UpdateCategory(ctx context.Context, category finmodels.Category) (*finpb.Category, error)
	DeleteCategory(ctx context.Context, userID, categoryID int) error
}
