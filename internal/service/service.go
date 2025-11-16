package service

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/storage/image"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/balance"
	budgetService "github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/budget"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/category"
	imageservice "github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/image"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/operation"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

type Repository interface {
	Close() error
	UserRepository
	AccountRepository
	BudgetRepository
	OperationRepository
	CategoryRepository
}

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
	GetUserByID(ctx context.Context, id int) (models.User, error)
	UpdateUser(ctx context.Context, user models.User) error
	EditUserByID(ctx context.Context, req models.UpdateProfileRequest, userID int) (models.User, error)
}

type AccountRepository interface {
	GetAccountsByUser(ctx context.Context, userID int) ([]models.Account, error)
	CreateAccount(ctx context.Context, account models.Account, userID int) (models.Account, error)
	UpdateAccount(ctx context.Context, req models.UpdateAccountRequest, userID, accID int) (models.Account, error)
	DeleteAccount(ctx context.Context, userID, accID int) (models.Account, error)
}

type BudgetRepository interface {
	GetBudgetsByUser(ctx context.Context, userID int) ([]models.Budget, error)
	GetAccountsByUser(ctx context.Context, userID int) ([]models.Account, error)
	GetOperationsByAccount(ctx context.Context, accountID int) ([]models.OperationInList, error)
	CreateBudget(ctx context.Context, budget models.Budget) (models.Budget, error)
	UpdateBudget(ctx context.Context, req models.UpdatedBudgetRequest, userID, budgetID int) (models.Budget, error)
	DeleteBudget(ctx context.Context, budgetID int) (models.Budget, error)
}

type OperationRepository interface {
	GetOperationsByAccount(ctx context.Context, accountID int) ([]models.OperationInList, error)
	GetOperationsByUser(ctx context.Context, userID int) ([]models.Operation, error)
	GetOperationByID(ctx context.Context, accID int, opID int) (models.Operation, error)
	CreateOperation(ctx context.Context, op models.Operation) (models.Operation, error)
	UpdateOperation(ctx context.Context, req models.UpdateOperationRequest, accID int, opID int) (models.Operation, error)
	DeleteOperation(ctx context.Context, accID int, opID int) (models.Operation, error)
}

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category models.Category) (models.Category, error)
	GetCategoriesByUser(ctx context.Context, userID int) ([]models.Category, error)
	GetCategoryByID(ctx context.Context, userID, categoryID int) (models.Category, error)
	UpdateCategory(ctx context.Context, category models.Category) error
	DeleteCategory(ctx context.Context, userID, categoryID int) error
	GetCategoryStats(ctx context.Context, userID, categoryID int) (int, error)
}

type Service struct {
	BalanceUC  balance.BalanceService
	BudgetUC   budgetService.BudgetService
	OpUC       operation.OperationService
	CategoryUC category.CategoryService
	ImageUC    imageservice.ImageService
}

func NewService(store Repository, jwtSecret string, imageStorage image.ImageStorage) *Service {
	realClock := clock.RealClock{}

	balanceService := balance.NewService(store, realClock)
	budgetService := budgetService.NewService(store, realClock)
	opService := operation.NewService(store, realClock)
	categoryService := category.NewService(store)
	imageService := imageservice.NewService(imageStorage)

	return &Service{
		BalanceUC:  balanceService,
		BudgetUC:   budgetService,
		OpUC:       opService,
		CategoryUC: categoryService,
		ImageUC:    imageService,
	}
}
