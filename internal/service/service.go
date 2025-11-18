package service

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/storage/image"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/balance"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/category"
	imageservice "github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/image"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/operation"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

type Repository interface {
	Close() error
	AccountRepository
	OperationRepository
	CategoryRepository
}

type AccountRepository interface {
	GetAccountsByUser(ctx context.Context, userID int) ([]models.Account, error)
	CreateAccount(ctx context.Context, account models.Account, userID int) (models.Account, error)
	UpdateAccount(ctx context.Context, req models.UpdateAccountRequest, userID, accID int) (models.Account, error)
	DeleteAccount(ctx context.Context, userID, accID int) (models.Account, error)
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
	OpUC       operation.OperationService
	CategoryUC category.CategoryService
	ImageUC    imageservice.ImageService
}

func NewService(store Repository, jwtSecret string, imageStorage image.ImageStorage) *Service {
	realClock := clock.RealClock{}

	balanceService := balance.NewService(store, realClock)
	opService := operation.NewService(store, realClock)
	categoryService := category.NewService(store)
	imageService := imageservice.NewService(imageStorage)

	return &Service{
		BalanceUC:  balanceService,
		OpUC:       opService,
		CategoryUC: categoryService,
		ImageUC:    imageService,
	}
}
