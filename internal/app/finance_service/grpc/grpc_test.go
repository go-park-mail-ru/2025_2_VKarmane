package grpc

import (
	"context"
	"errors"
	"testing"

	finerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/errors"
	finmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/models"
	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestFinanceServer_CreateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockFinanceUseCase(ctrl)
	server := NewFinanceServer(mockUC)

	req := &finpb.CreateAccountRequest{UserId: 1, Balance: 12}
	expected := &finpb.Account{Id: 1, Balance: 12}

	mockUC.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(expected, nil)

	resp, err := server.CreateAccount(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp)
}

func TestFinanceServer_GetAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockFinanceUseCase(ctrl)
	server := NewFinanceServer(mockUC)

	req := &finpb.AccountRequest{UserId: 1, AccountId: 2}
	expected := &finpb.Account{Id: 2, Balance: 2}

	mockUC.EXPECT().GetAccountByID(gomock.Any(), 1, 2).Return(expected, nil)

	resp, err := server.GetAccount(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expected.Id, resp.Id)
}

func TestFinanceServer_DeleteAccount_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockFinanceUseCase(ctrl)
	server := NewFinanceServer(mockUC)

	req := &finpb.AccountRequest{UserId: 1, AccountId: 2}

	mockUC.EXPECT().DeleteAccount(gomock.Any(), 1, 2).Return(nil, errors.New("some error"))

	resp, err := server.DeleteAccount(context.Background(), req)
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestFinanceServer_UpdateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockFinanceUseCase(ctrl)
	server := NewFinanceServer(mockUC)

	req := &finpb.UpdateAccountRequest{UserId: 1, AccountId: 1, Balance: 12.0}
	expected := &finpb.Account{Id: 1, Balance: 12.0}

	mockUC.EXPECT().
		UpdateAccount(gomock.Any(), gomock.AssignableToTypeOf(finmodels.UpdateAccountRequest{})).
		DoAndReturn(func(_ context.Context, req finmodels.UpdateAccountRequest) (*finpb.Account, error) {
			assert.Equal(t, 1, req.UserID)
			assert.Equal(t, 1, req.AccountID)
			assert.Equal(t, 12.0, req.Balance)
			return &finpb.Account{Id: int32(req.AccountID), Balance: req.Balance}, nil
		})

	resp, err := server.UpdateAccount(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expected.Balance, resp.Balance)
}

func TestFinanceServer_CreateOperation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockFinanceUseCase(ctrl)
	server := NewFinanceServer(mockUC)

	req := &finpb.CreateOperationRequest{AccountId: 1, Sum: 100}
	expected := &finpb.Operation{Id: 1, AccountId: 1, Sum: 100}

	mockUC.EXPECT().CreateOperation(gomock.Any(), gomock.Any(), 1).Return(expected, nil)

	resp, err := server.CreateOperation(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expected.Sum, resp.Sum)
}

func TestFinanceServer_GetOperation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockFinanceUseCase(ctrl)
	server := NewFinanceServer(mockUC)

	req := &finpb.OperationRequest{AccountId: 1, OperationId: 2}
	expected := &finpb.Operation{Id: 2, AccountId: 1, Sum: 50}

	mockUC.EXPECT().GetOperationByID(gomock.Any(), 1, 2).Return(expected, nil)

	resp, err := server.GetOperation(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expected.AccountId, resp.AccountId)
}

func TestFinanceServer_DeleteOperation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockFinanceUseCase(ctrl)
	server := NewFinanceServer(mockUC)

	req := &finpb.OperationRequest{AccountId: 1, OperationId: 2}
	expected := &finpb.Operation{Id: 2, AccountId: 1, Sum: 50}

	mockUC.EXPECT().DeleteOperation(gomock.Any(), 1, 2).Return(expected, nil)

	resp, err := server.DeleteOperation(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expected.AccountId, resp.AccountId)
}

func TestFinanceServer_CreateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockFinanceUseCase(ctrl)
	server := NewFinanceServer(mockUC)

	req := &finpb.CreateCategoryRequest{UserId: 1, Name: "Food"}
	expected := &finpb.Category{Id: 1, UserId: 1, Name: "Food"}

	mockUC.EXPECT().CreateCategory(gomock.Any(), gomock.Any()).Return(expected, nil)

	resp, err := server.CreateCategory(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expected.Name, resp.Name)
}

func TestFinanceServer_UpdateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockFinanceUseCase(ctrl)
	server := NewFinanceServer(mockUC)
	name := "updated"

	req := &finpb.UpdateCategoryRequest{UserId: 1, CategoryId: 1, Name: &name}
	expected := &finpb.Category{Id: 1, UserId: 1, Name: "updated"}

	mockUC.EXPECT().
		UpdateCategory(gomock.Any(), gomock.AssignableToTypeOf(finmodels.Category{})).
		DoAndReturn(func(_ context.Context, cat finmodels.Category) (*finpb.Category, error) {
			assert.Equal(t, 1, cat.ID)
			assert.Equal(t, 1, cat.UserID)
			assert.Equal(t, "updated", cat.Name) // проверяем фактическое значение
			return &finpb.Category{Id: int32(cat.ID), UserId: int32(cat.UserID), Name: cat.Name}, nil
		})

	resp, err := server.UpdateCategory(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expected.Name, resp.Name)
}

func TestFinanceServer_GetAccountsByUser_ErrorMap(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockFinanceUseCase(ctrl)
	server := NewFinanceServer(mockUC)

	mockUC.EXPECT().GetAccountsByUser(gomock.Any(), 1).Return(nil, finerrors.ErrAccountNotFound)

	resp, err := server.GetAccountsByUser(context.Background(), &finpb.UserID{UserId: 1})
	assert.Nil(t, resp)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
}

func TestFinanceServer_CreateOperation_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockFinanceUseCase(ctrl)
	server := NewFinanceServer(mockUC)

	mockUC.EXPECT().
		CreateOperation(gomock.Any(), gomock.AssignableToTypeOf(finmodels.CreateOperationRequest{}), 2).
		Return(nil, finerrors.ErrForbidden)

	resp, err := server.CreateOperation(context.Background(), &finpb.CreateOperationRequest{AccountId: 2})
	assert.Nil(t, resp)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.PermissionDenied, st.Code())
}

func TestFinanceServer_GetCategory_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockFinanceUseCase(ctrl)
	server := NewFinanceServer(mockUC)

	mockUC.EXPECT().GetCategoryByID(gomock.Any(), 1, 2).Return(nil, finerrors.ErrCategoryNotFound)

	resp, err := server.GetCategory(context.Background(), &finpb.CategoryRequest{UserId: 1, CategoryId: 2})
	assert.Nil(t, resp)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
}

func TestFinanceServer_GetCategoriesByUser_Internal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockFinanceUseCase(ctrl)
	server := NewFinanceServer(mockUC)

	mockUC.EXPECT().GetCategoriesByUser(gomock.Any(), 1).Return(nil, errors.New("boom"))

	resp, err := server.GetCategoriesByUser(context.Background(), &finpb.UserID{UserId: 1})
	assert.Nil(t, resp)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.Internal, st.Code())
}
