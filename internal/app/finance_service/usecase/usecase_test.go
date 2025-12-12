package usecase

import (
	"context"
	"errors"
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



func TestGetAccountsByUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockFinanceService(ctrl)
	uc := NewFinanceUseCase(mockSvc)

	ctx := context.Background()
	expected := &finpb.ListAccountsResponse{}

	mockSvc.EXPECT().GetAccountsByUser(ctx, 1).Return(expected, nil)

	res, err := uc.GetAccountsByUser(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, expected, res)
}

func TestCreateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockFinanceService(ctrl)
	uc := NewFinanceUseCase(mockSvc)

	ctx := context.Background()
	req := models.CreateAccountRequest{UserID: 1}
	expected := &finpb.Account{}

	mockSvc.EXPECT().CreateAccount(ctx, req).Return(expected, nil)

	res, err := uc.CreateAccount(ctx, req)
	require.NoError(t, err)
	require.Equal(t, expected, res)
}

func TestDeleteAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockFinanceService(ctrl)
	uc := NewFinanceUseCase(mockSvc)

	ctx := context.Background()
	expected := &finpb.Account{}

	mockSvc.EXPECT().DeleteAccount(ctx, 1, 2).Return(expected, nil)

	res, err := uc.DeleteAccount(ctx, 1, 2)
	require.NoError(t, err)
	require.Equal(t, expected, res)
}

func TestCreateOperation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockFinanceService(ctrl)
	uc := NewFinanceUseCase(mockSvc)

	ctx := context.Background()
	req := models.CreateOperationRequest{UserID: 1}
	expected := &finpb.Operation{}

	mockSvc.EXPECT().CreateOperation(ctx, req, 2).Return(expected, nil)

	res, err := uc.CreateOperation(ctx, req, 2)
	require.NoError(t, err)
	require.Equal(t, expected, res)
}

func TestDeleteOperation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockFinanceService(ctrl)
	uc := NewFinanceUseCase(mockSvc)

	ctx := context.Background()
	expected := &finpb.Operation{}

	mockSvc.EXPECT().DeleteOperation(ctx, 2, 3).Return(expected, nil)

	res, err := uc.DeleteOperation(ctx, 2, 3)
	require.NoError(t, err)
	require.Equal(t, expected, res)
}

func TestCreateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockFinanceService(ctrl)
	uc := NewFinanceUseCase(mockSvc)

	ctx := context.Background()
	req := models.CreateCategoryRequest{UserID: 1}
	expected := &finpb.Category{}

	mockSvc.EXPECT().CreateCategory(ctx, req).Return(expected, nil)

	res, err := uc.CreateCategory(ctx, req)
	require.NoError(t, err)
	require.Equal(t, expected, res)
}

func TestDeleteCategoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockFinanceService(ctrl)
	uc := NewFinanceUseCase(mockSvc)

	ctx := context.Background()
	mockSvc.EXPECT().DeleteCategory(ctx, 1, 2).Return(errors.New("boom"))

	err := uc.DeleteCategory(ctx, 1, 2)
	require.Error(t, err)
	require.ErrorContains(t, err, "finance.DeleteCategory")
}


func ptrFloat64(v float64) *float64 { return &v }
