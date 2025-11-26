package service

import (
	"context"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/models"
	mock_repo "github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetAccountByID(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock_repo.NewMockFinanceRepository(ctrl)
	svc := NewService(mockRepo, fixedClock)

	ctx := context.Background()
	userID, accountID := 1, 2
	acc := models.Account{ID: accountID, Balance: 100}

	mockRepo.EXPECT().GetAccountByID(ctx, userID, accountID).Return(acc, nil)

	res, err := svc.GetAccountByID(ctx, userID, accountID)
	require.NoError(t, err)
	require.Equal(t, int32(accountID), res.Id)
	require.Equal(t, float64(100), res.Balance)
}

func TestCreateAccount(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock_repo.NewMockFinanceRepository(ctrl)
	svc := NewService(mockRepo, fixedClock)

	ctx := context.Background()
	req := models.CreateAccountRequest{UserID: 1, Balance: 50, Type: "cash", CurrencyID: 2}
	acc := models.Account{ID: 5, Balance: 50, Type: "cash", CurrencyID: 2, CreatedAt: fixedClock.FixedTime, UpdatedAt: fixedClock.FixedTime}

	mockRepo.EXPECT().CreateAccount(ctx, gomock.Any(), req.UserID).Return(acc, nil)

	res, err := svc.CreateAccount(ctx, req)
	require.NoError(t, err)
	require.Equal(t, int32(5), res.Id)
	require.Equal(t, float64(50), res.Balance)
}

func TestUpdateAccount(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock_repo.NewMockFinanceRepository(ctrl)
	svc := NewService(mockRepo, fixedClock)

	ctx := context.Background()
	req := models.UpdateAccountRequest{UserID: 1, AccountID: 2, Balance: 200}
	acc := models.Account{ID: 2, Balance: 200}

	mockRepo.EXPECT().UpdateAccount(ctx, req).Return(acc, nil)

	res, err := svc.UpdateAccount(ctx, req)
	require.NoError(t, err)
	require.Equal(t, float64(200), res.Balance)
}

func TestDeleteAccount(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock_repo.NewMockFinanceRepository(ctrl)
	svc := NewService(mockRepo, fixedClock)

	ctx := context.Background()
	userID, accountID := 1, 2
	acc := models.Account{ID: accountID, Balance: 100}

	mockRepo.EXPECT().DeleteAccount(ctx, userID, accountID).Return(acc, nil)

	res, err := svc.DeleteAccount(ctx, userID, accountID)
	require.NoError(t, err)
	require.Equal(t, int32(accountID), res.Id)
}

func TestCreateOperation_UpdateAccountBalance(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock_repo.NewMockFinanceRepository(ctrl)
	svc := NewService(mockRepo, fixedClock)

	ctx := context.Background()
	req := models.CreateOperationRequest{
		UserID: 1, AccountID: 2, Type: models.OperationIncome, Sum: 100,
	}

	mockRepo.EXPECT().CreateOperation(ctx, gomock.Any()).Return(models.Operation{ID: 1, AccountID: 2, Sum: 100, Type: models.OperationIncome}, nil)

	op, err := svc.CreateOperation(ctx, req, req.AccountID)
	require.NoError(t, err)
	require.Equal(t, int32(1), op.Id)
	require.Equal(t, float64(100), op.Sum)
}

func TestUpdateOperation(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock_repo.NewMockFinanceRepository(ctrl)
	svc := NewService(mockRepo, fixedClock)

	ctx := context.Background()
	req := models.UpdateOperationRequest{UserID: 1, AccountID: 2, OperationID: 3}
	op := models.Operation{ID: 3, AccountID: 2}

	mockRepo.EXPECT().UpdateOperation(ctx, req, req.AccountID, req.OperationID).Return(op, nil)

	res, err := svc.UpdateOperation(ctx, req)
	require.NoError(t, err)
	require.Equal(t, int32(3), res.Id)
}

func TestCreateCategory_DefaultLogo(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock_repo.NewMockFinanceRepository(ctrl)
	svc := NewService(mockRepo, fixedClock)

	ctx := context.Background()
	req := models.CreateCategoryRequest{UserID: 1, Name: "food"}
	cat := models.Category{ID: 5, UserID: 1, Name: "food", LogoHashedID: "c1dfd96eea8cc2b62785275bca38ac261256e278"}

	mockRepo.EXPECT().CreateCategory(ctx, gomock.Any()).Return(cat, nil)

	res, err := svc.CreateCategory(ctx, req)
	require.NoError(t, err)
	require.Equal(t, int32(5), res.Id)
	require.Equal(t, "food", res.Name)
	require.Equal(t, "c1dfd96eea8cc2b62785275bca38ac261256e278", res.LogoHashedId)
}

func TestUpdateCategory(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock_repo.NewMockFinanceRepository(ctrl)
	svc := NewService(mockRepo, fixedClock)

	ctx := context.Background()
	category := models.Category{ID: 1, UserID: 2, Name: "transport"}

	mockRepo.EXPECT().UpdateCategory(ctx, category).Return(nil)
	mockRepo.EXPECT().GetCategoryByID(ctx, category.UserID, category.ID).Return(category, nil)

	res, err := svc.UpdateCategory(ctx, category)
	require.NoError(t, err)
	require.Equal(t, "transport", res.Name)
}

func TestDeleteCategory(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock_repo.NewMockFinanceRepository(ctrl)
	svc := NewService(mockRepo, fixedClock)

	ctx := context.Background()
	userID, categoryID := 1, 2
	mockRepo.EXPECT().DeleteCategory(ctx, userID, categoryID).Return(nil)

	err := svc.DeleteCategory(ctx, userID, categoryID)
	require.NoError(t, err)
}

func TestGetCategoriesWithStatsByUser(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock_repo.NewMockFinanceRepository(ctrl)
	svc := NewService(mockRepo, fixedClock)

	ctx := context.Background()
	userID := 1
	cats := []models.CategoryWithStats{
		{Category: models.Category{ID: 1, UserID: 1, Name: "food"}, OperationsCount: 10},
	}

	mockRepo.EXPECT().GetCategoriesWithStatsByUser(ctx, userID).Return(cats, nil)

	res, err := svc.GetCategoriesWithStatsByUser(ctx, userID)
	require.NoError(t, err)
	require.Len(t, res.Categories, 1)
	require.Equal(t, int32(10), res.Categories[0].OperationsCount)
}
