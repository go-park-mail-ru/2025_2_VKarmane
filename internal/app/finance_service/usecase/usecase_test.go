package usecase

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/models"
	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetAccountByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockFinanceService(ctrl)
	uc := NewFinanceUseCase(mockSvc)

	ctx := context.Background()
	userID, accID := 10, 7
	expected := &finpb.Account{}

	mockSvc.EXPECT().GetAccountByID(ctx, userID, accID).Return(expected, nil)

	res, err := uc.GetAccountByID(ctx, userID, accID)
	require.NoError(t, err)
	require.Equal(t, expected, res)
}

func TestUpdateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockFinanceService(ctrl)
	uc := NewFinanceUseCase(mockSvc)
	balance := 200.

	ctx := context.Background()
	req := models.UpdateAccountRequest{UserID: 10, AccountID: 7, Balance: &balance}
	expected := &finpb.Account{}

	mockSvc.EXPECT().UpdateAccount(ctx, req).Return(expected, nil)

	res, err := uc.UpdateAccount(ctx, req)
	require.NoError(t, err)
	require.Equal(t, expected, res)
}

func TestGetOperationsByAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockFinanceService(ctrl)
	uc := NewFinanceUseCase(mockSvc)

	ctx := context.Background()
	accountID := 1
	ctgID := []int{1}
	opName := "name"
	opType := "type"
	accType := "type"
	date := "date"

	expected := &finpb.ListOperationsResponse{}

	mockSvc.EXPECT().
		GetOperationsByAccount(ctx, gomock.Any()).
		Return(expected, nil)

	res, err := uc.GetOperationsByAccount(ctx, accountID, ctgID, opName, opType, accType, date)
	require.NoError(t, err)
	require.Equal(t, expected, res)
}

func TestGetOperationByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockFinanceService(ctrl)
	uc := NewFinanceUseCase(mockSvc)

	ctx := context.Background()
	accID, opID := 1, 2
	expected := &finpb.Operation{}

	mockSvc.EXPECT().GetOperationByID(ctx, accID, opID).Return(expected, nil)

	res, err := uc.GetOperationByID(ctx, accID, opID)
	require.NoError(t, err)
	require.Equal(t, expected, res)
}

func TestUpdateOperation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockFinanceService(ctrl)
	uc := NewFinanceUseCase(mockSvc)

	ctx := context.Background()
	req := models.UpdateOperationRequest{
		UserID: 1, AccountID: 2, OperationID: 3, Sum: ptrFloat64(500),
	}
	expected := &finpb.Operation{}

	mockSvc.EXPECT().UpdateOperation(ctx, req).Return(expected, nil)

	res, err := uc.UpdateOperation(ctx, req)
	require.NoError(t, err)
	require.Equal(t, expected, res)
}

func TestGetCategoriesByUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockFinanceService(ctrl)
	uc := NewFinanceUseCase(mockSvc)

	ctx := context.Background()
	userID := 1
	expected := &finpb.ListCategoriesResponse{}

	mockSvc.EXPECT().GetCategoriesByUser(ctx, userID).Return(expected, nil)

	res, err := uc.GetCategoriesByUser(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, expected, res)
}

func TestGetCategoriesWithStatsByUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockFinanceService(ctrl)
	uc := NewFinanceUseCase(mockSvc)

	ctx := context.Background()
	userID := 1
	expected := &finpb.ListCategoriesWithStatsResponse{}

	mockSvc.EXPECT().GetCategoriesWithStatsByUser(ctx, userID).Return(expected, nil)

	res, err := uc.GetCategoriesWithStatsByUser(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, expected, res)
}

func TestGetCategoryByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockFinanceService(ctrl)
	uc := NewFinanceUseCase(mockSvc)

	ctx := context.Background()
	userID, categoryID := 1, 2
	expected := &finpb.CategoryWithStats{}

	mockSvc.EXPECT().GetCategoryByID(ctx, userID, categoryID).Return(expected, nil)

	res, err := uc.GetCategoryByID(ctx, userID, categoryID)
	require.NoError(t, err)
	require.Equal(t, expected, res)
}

func TestUpdateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockFinanceService(ctrl)
	uc := NewFinanceUseCase(mockSvc)

	ctx := context.Background()
	category := models.Category{ID: 2, UserID: 1, Name: "cat"}
	expected := &finpb.Category{}

	mockSvc.EXPECT().UpdateCategory(ctx, category).Return(expected, nil)

	res, err := uc.UpdateCategory(ctx, category)
	require.NoError(t, err)
	require.Equal(t, expected, res)
}

func ptrFloat64(v float64) *float64 { return &v }
