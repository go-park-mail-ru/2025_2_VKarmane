package account

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	"github.com/gorilla/mux"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetAccounts_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockFinanceServiceClient(ctrl)
	handler := NewHandler(mockClient, clock.RealClock{})

	req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
	rr := httptest.NewRecorder()

	handler.GetAccounts(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetAccounts_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockFinanceServiceClient(ctrl)
	handler := NewHandler(mockClient, clock.RealClock{})

	resp := &finpb.ListAccountsResponse{
		Accounts: []*finpb.Account{
			{
				Id:         1,
				Balance:    1000.5,
				Type:       "debit",
				CurrencyId: 1,
			},
		},
	}

	mockClient.EXPECT().
		GetAccountsByUser(gomock.Any(), &finpb.UserID{UserId: 1}).
		Return(resp, nil)

	req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.GetAccounts(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestGetAccountByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockFinanceServiceClient(ctrl)
	handler := NewHandler(mockClient, clock.RealClock{})

	mockClient.EXPECT().
		GetAccount(gomock.Any(), &finpb.AccountRequest{UserId: 1, AccountId: 5}).
		Return(nil, status.Error(codes.NotFound, "not found"))

	req := httptest.NewRequest(http.MethodGet, "/balance/5", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	req = mux.SetURLVars(req, map[string]string{"id": "5"})

	rr := httptest.NewRecorder()
	handler.GetAccountByID(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateAccount_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockFinanceServiceClient(ctrl)
	handler := NewHandler(mockClient, clock.RealClock{})

	req := httptest.NewRequest(http.MethodPost, "/balance", bytes.NewBuffer([]byte(`invalid`)))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))

	rr := httptest.NewRecorder()
	handler.CreateAccount(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdateAccount_Forbidden(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockFinanceServiceClient(ctrl)
	handler := NewHandler(mockClient, clock.RealClock{})
	balance := 100.0

	body, _ := json.Marshal(models.UpdateAccountRequest{Balance: &balance})

	mockClient.EXPECT().
		UpdateAccount(gomock.Any(),
			&finpb.UpdateAccountRequest{
				UserId:    1,
				AccountId: 5,
				Balance:   &balance,
			},
		).
		Return(nil, status.Error(codes.PermissionDenied, "forbidden"))

	req := httptest.NewRequest(http.MethodPut, "/balance/5", bytes.NewBuffer(body))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	req = mux.SetURLVars(req, map[string]string{"id": "5"})

	rr := httptest.NewRecorder()
	handler.UpdateAccount(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestDeleteAccount_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockFinanceServiceClient(ctrl)
	handler := NewHandler(mockClient, clock.RealClock{})

	mockClient.EXPECT().
		DeleteAccount(gomock.Any(),
			&finpb.AccountRequest{UserId: 1, AccountId: 7},
		).
		Return(nil, status.Error(codes.NotFound, "not found"))

	req := httptest.NewRequest(http.MethodDelete, "/balance/7", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	req = mux.SetURLVars(req, map[string]string{"id": "7"})

	rr := httptest.NewRecorder()
	handler.DeleteAccount(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetAccountByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockFinanceServiceClient(ctrl)

	h := NewHandler(mockClient, clock.RealClock{})

	grpcResp := &finpb.Account{Id: 7, Balance: 999, Type: "debit", CurrencyId: 1}

	mockClient.EXPECT().
		GetAccount(gomock.Any(), &finpb.AccountRequest{UserId: 10, AccountId: 7}).
		Return(grpcResp, nil)

	req := httptest.NewRequest(http.MethodGet, "/balance/7", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 10))
	req = mux.SetURLVars(req, map[string]string{"id": "7"})

	rr := httptest.NewRecorder()

	h.GetAccountByID(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestCreateAccount_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockFinanceServiceClient(ctrl)
	h := NewHandler(mockClient, clock.RealClock{})

	body, _ := json.Marshal(map[string]interface{}{
		"balance":     200.5,
		"type":        "debit",
		"currency_id": 1,
	})

	mockClient.EXPECT().
		CreateAccount(gomock.Any(), &finpb.CreateAccountRequest{
			UserId:     10,
			Balance:    200.5,
			Type:       "debit",
			CurrencyId: 1,
		}).
		Return(&finpb.Account{Id: 3, Balance: 200.5, Type: "debit", CurrencyId: 1}, nil)

	req := httptest.NewRequest(http.MethodPost, "/balance", bytes.NewBuffer(body))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 10))

	rr := httptest.NewRecorder()

	h.CreateAccount(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
}

func TestUpdateAccount_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockFinanceServiceClient(ctrl)
	h := NewHandler(mockClient, clock.RealClock{})

	body, _ := json.Marshal(map[string]interface{}{"balance": 999})
	balance := 999.

	mockClient.EXPECT().
		UpdateAccount(gomock.Any(), &finpb.UpdateAccountRequest{
			UserId:    10,
			AccountId: 5,
			Balance:   &balance,
		}).
		Return(&finpb.Account{Id: 5, Balance: 999, Type: "debit", CurrencyId: 1}, nil)

	req := httptest.NewRequest(http.MethodPut, "/balance/5", bytes.NewBuffer(body))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 10))
	req = mux.SetURLVars(req, map[string]string{"id": "5"})

	rr := httptest.NewRecorder()

	h.UpdateAccount(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteAccount_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockFinanceServiceClient(ctrl)
	h := NewHandler(mockClient, clock.RealClock{})

	mockClient.EXPECT().
		DeleteAccount(gomock.Any(), &finpb.AccountRequest{UserId: int32(10), AccountId: int32(7)}).
		Return(&finpb.Account{}, nil)

	req := httptest.NewRequest(http.MethodDelete, "/balance/7", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 10))
	req = mux.SetURLVars(req, map[string]string{"id": "7"})

	rr := httptest.NewRecorder()

	h.DeleteAccount(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}
