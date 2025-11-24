package usecase

import (
	"context"

	pkgerrors "github.com/pkg/errors"

	finmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/models"
	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
)

type UseCase struct {
	financeService FinanceService
}

func NewFinanceUseCase(svc FinanceService) *UseCase {
	return &UseCase{
		financeService: svc,
	}
}

// Account methods
func (uc *UseCase) GetAccountsByUser(ctx context.Context, userID int) (*finpb.ListAccountsResponse, error) {
	log := logger.FromContext(ctx)
	accounts, err := uc.financeService.GetAccountsByUser(ctx, userID)
	if err != nil {
		if log != nil {
			log.Error("Failed to get accounts for user", "error", err, "user_id", userID)
		}
		return nil, pkgerrors.Wrap(err, "finance.GetAccountsByUser")
	}
	return accounts, nil
}

func (uc *UseCase) GetAccountByID(ctx context.Context, userID, accountID int) (*finpb.Account, error) {
	log := logger.FromContext(ctx)
	account, err := uc.financeService.GetAccountByID(ctx, userID, accountID)
	if err != nil {
		if log != nil {
			log.Error("Failed to get account by ID", "error", err, "user_id", userID, "account_id", accountID)
		}
		return nil, pkgerrors.Wrap(err, "finance.GetAccountByID")
	}
	return account, nil
}

func (uc *UseCase) CreateAccount(ctx context.Context, req finmodels.CreateAccountRequest) (*finpb.Account, error) {
	log := logger.FromContext(ctx)
	account, err := uc.financeService.CreateAccount(ctx, req)
	if err != nil {
		if log != nil {
			log.Error("Failed to create account", "error", err, "user_id", req.UserID)
		}
		return nil, pkgerrors.Wrap(err, "finance.CreateAccount")
	}
	return account, nil
}

func (uc *UseCase) UpdateAccount(ctx context.Context, req finmodels.UpdateAccountRequest) (*finpb.Account, error) {
	log := logger.FromContext(ctx)
	account, err := uc.financeService.UpdateAccount(ctx, req)
	if err != nil {
		if log != nil {
			log.Error("Failed to update account", "error", err, "user_id", req.UserID, "account_id", req.AccountID)
		}
		return nil, pkgerrors.Wrap(err, "finance.UpdateAccount")
	}
	return account, nil
}

func (uc *UseCase) DeleteAccount(ctx context.Context, userID, accountID int) (*finpb.Account, error) {
	log := logger.FromContext(ctx)
	account, err := uc.financeService.DeleteAccount(ctx, userID, accountID)
	if err != nil {
		if log != nil {
			log.Error("Failed to delete account", "error", err, "user_id", userID, "account_id", accountID)
		}
		return nil, pkgerrors.Wrap(err, "finance.DeleteAccount")
	}
	return account, nil
}

// Operation methods
func (uc *UseCase) GetOperationsByAccount(ctx context.Context, accountID int) (*finpb.ListOperationsResponse, error) {
	log := logger.FromContext(ctx)
	operations, err := uc.financeService.GetOperationsByAccount(ctx, accountID)
	if err != nil {
		if log != nil {
			log.Error("Failed to get operations for account", "error", err, "account_id", accountID)
		}
		return nil, pkgerrors.Wrap(err, "finance.GetOperationsByAccount")
	}
	return operations, nil
}

func (uc *UseCase) GetOperationByID(ctx context.Context, accID, opID int) (*finpb.Operation, error) {
	log := logger.FromContext(ctx)
	operation, err := uc.financeService.GetOperationByID(ctx, accID, opID)
	if err != nil {
		if log != nil {
			log.Error("Failed to get operation by ID", "error", err, "account_id", accID, "operation_id", opID)
		}
		return nil, pkgerrors.Wrap(err, "finance.GetOperationByID")
	}
	return operation, nil
}

func (uc *UseCase) CreateOperation(ctx context.Context, req finmodels.CreateOperationRequest, accountID int) (*finpb.Operation, error) {
	log := logger.FromContext(ctx)
	operation, err := uc.financeService.CreateOperation(ctx, req, accountID)
	if err != nil {
		if log != nil {
			log.Error("Failed to create operation", "error", err, "user_id", req.UserID, "account_id", accountID)
		}
		return nil, pkgerrors.Wrap(err, "finance.CreateOperation")
	}
	return operation, nil
}

func (uc *UseCase) UpdateOperation(ctx context.Context, req finmodels.UpdateOperationRequest) (*finpb.Operation, error) {
	log := logger.FromContext(ctx)
	operation, err := uc.financeService.UpdateOperation(ctx, req)
	if err != nil {
		if log != nil {
			log.Error("Failed to update operation", "error", err, "user_id", req.UserID, "operation_id", req.OperationID)
		}
		return nil, pkgerrors.Wrap(err, "finance.UpdateOperation")
	}
	return operation, nil
}

func (uc *UseCase) DeleteOperation(ctx context.Context, accID, opID int) (*finpb.Operation, error) {
	log := logger.FromContext(ctx)
	operation, err := uc.financeService.DeleteOperation(ctx, accID, opID)
	if err != nil {
		if log != nil {
			log.Error("Failed to delete operation", "error", err, "account_id", accID, "operation_id", opID)
		}
		return nil, pkgerrors.Wrap(err, "finance.DeleteOperation")
	}
	return operation, nil
}

// Category methods
func (uc *UseCase) CreateCategory(ctx context.Context, req finmodels.CreateCategoryRequest) (*finpb.Category, error) {
	log := logger.FromContext(ctx)
	category, err := uc.financeService.CreateCategory(ctx, req)
	if err != nil {
		if log != nil {
			log.Error("Failed to create category", "error", err, "user_id", req.UserID)
		}
		return nil, pkgerrors.Wrap(err, "finance.CreateCategory")
	}
	return category, nil
}

func (uc *UseCase) GetCategoriesByUser(ctx context.Context, userID int) (*finpb.ListCategoriesResponse, error) {
	log := logger.FromContext(ctx)
	categories, err := uc.financeService.GetCategoriesByUser(ctx, userID)
	if err != nil {
		if log != nil {
			log.Error("Failed to get categories for user", "error", err, "user_id", userID)
		}
		return nil, pkgerrors.Wrap(err, "finance.GetCategoriesByUser")
	}
	return categories, nil
}

func (uc *UseCase) GetCategoriesWithStatsByUser(ctx context.Context, userID int) (*finpb.ListCategoriesWithStatsResponse, error) {
	log := logger.FromContext(ctx)
	categories, err := uc.financeService.GetCategoriesWithStatsByUser(ctx, userID)
	if err != nil {
		if log != nil {
			log.Error("Failed to get categories with stats for user", "error", err, "user_id", userID)
		}
		return nil, pkgerrors.Wrap(err, "finance.GetCategoriesWithStatsByUser")
	}
	return categories, nil
}

func (uc *UseCase) GetCategoryByID(ctx context.Context, userID, categoryID int) (*finpb.CategoryWithStats, error) {
	log := logger.FromContext(ctx)
	category, err := uc.financeService.GetCategoryByID(ctx, userID, categoryID)
	if err != nil {
		if log != nil {
			log.Error("Failed to get category by ID", "error", err, "user_id", userID, "category_id", categoryID)
		}
		return nil, pkgerrors.Wrap(err, "finance.GetCategoryByID")
	}
	return category, nil
}

func (uc *UseCase) UpdateCategory(ctx context.Context, category finmodels.Category) (*finpb.Category, error) {
	log := logger.FromContext(ctx)
	updatedCategory, err := uc.financeService.UpdateCategory(ctx, category)
	if err != nil {
		if log != nil {
			log.Error("Failed to update category", "error", err, "user_id", category.UserID, "category_id", category.ID)
		}
		return nil, pkgerrors.Wrap(err, "finance.UpdateCategory")
	}
	return updatedCategory, nil
}

func (uc *UseCase) DeleteCategory(ctx context.Context, userID, categoryID int) error {
	log := logger.FromContext(ctx)
	err := uc.financeService.DeleteCategory(ctx, userID, categoryID)
	if err != nil {
		if log != nil {
			log.Error("Failed to delete category", "error", err, "user_id", userID, "category_id", categoryID)
		}
		return pkgerrors.Wrap(err, "finance.DeleteCategory")
	}
	return nil
}
