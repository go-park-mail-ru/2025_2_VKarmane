package balance

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	accounterrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/account/repository"
	serviceerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/errors"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetAccounts_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockBalanceUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
	rr := httptest.NewRecorder()

	handler.GetAccounts(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetAccounts_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockBalanceUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	accounts := []models.Account{
		{ID: 1, Balance: 1000.50, Type: "debit"},
	}

	mockUC.EXPECT().GetBalanceForUser(gomock.Any(), 1).Return(accounts, nil)

	req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.GetAccounts(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestGetAccounts_EmptyAccounts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockBalanceUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	mockUC.EXPECT().GetBalanceForUser(gomock.Any(), 1).Return([]models.Account{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.GetAccounts(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestCreateAccount_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler := NewHandler(mocks.NewMockBalanceUseCase(ctrl), clock.RealClock{})

	req := httptest.NewRequest(http.MethodPost, "/balance", bytes.NewBuffer([]byte(`{}`)))
	rr := httptest.NewRecorder()

	handler.CreateAccount(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestCreateAccount_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler := NewHandler(mocks.NewMockBalanceUseCase(ctrl), clock.RealClock{})

	req := httptest.NewRequest(http.MethodPost, "/balance", bytes.NewBuffer([]byte(`invalid`)))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.CreateAccount(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateAccount_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockBalanceUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	reqBody := models.CreateAccountRequest{
		Type:       "debit",
		CurrencyID: 1,
	}
	body, _ := json.Marshal(reqBody)

	mockUC.EXPECT().CreateAccount(gomock.Any(), reqBody, 1).Return(models.Account{ID: 1, Type: "debit"}, nil)

	req := httptest.NewRequest(http.MethodPost, "/balance", bytes.NewBuffer(body))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.CreateAccount(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateAccount_Forbidden(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockBalanceUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	reqBody := models.UpdateAccountRequest{Balance: 12}
	body, _ := json.Marshal(reqBody)

	mockUC.EXPECT().
		UpdateAccount(gomock.Any(), reqBody, 1, 5).
		Return(models.Account{}, serviceerrors.ErrForbidden)

	req := httptest.NewRequest(http.MethodPut, "/balance/5", bytes.NewBuffer(body))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	req = mux.SetURLVars(req, map[string]string{"id": "5"})
	rr := httptest.NewRecorder()

	handler.UpdateAccount(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestUpdateAccount_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockBalanceUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	reqBody := models.UpdateAccountRequest{Balance: 12}
	body, _ := json.Marshal(reqBody)

	mockUC.EXPECT().
		UpdateAccount(gomock.Any(), reqBody, 1, 5).
		Return(models.Account{}, accounterrors.ErrAccountNotFound)

	req := httptest.NewRequest(http.MethodPut, "/balance/5", bytes.NewBuffer(body))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	req = mux.SetURLVars(req, map[string]string{"id": "5"})
	rr := httptest.NewRecorder()

	handler.UpdateAccount(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestUpdateAccount_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockBalanceUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	reqBody := models.UpdateAccountRequest{Balance: 1111}
	body, _ := json.Marshal(reqBody)

	acc := models.Account{ID: 5, Type: "debit", Balance: 150}
	mockUC.EXPECT().
		UpdateAccount(gomock.Any(), reqBody, 1, 5).
		Return(acc, nil)

	req := httptest.NewRequest(http.MethodPut, "/balance/5", bytes.NewBuffer(body))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	req = mux.SetURLVars(req, map[string]string{"id": "5"})
	rr := httptest.NewRecorder()

	handler.UpdateAccount(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteAccount_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockBalanceUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	acc := models.Account{ID: 7, Type: "credit", Balance: 1000}
	mockUC.EXPECT().DeleteAccount(gomock.Any(), 1, 7).Return(acc, nil)

	req := httptest.NewRequest(http.MethodDelete, "/balance/7", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	req = mux.SetURLVars(req, map[string]string{"id": "7"})
	rr := httptest.NewRecorder()

	handler.DeleteAccount(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteAccount_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockBalanceUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	mockUC.EXPECT().DeleteAccount(gomock.Any(), 1, 7).Return(models.Account{}, accounterrors.ErrAccountNotFound)

	req := httptest.NewRequest(http.MethodDelete, "/balance/7", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	req = mux.SetURLVars(req, map[string]string{"id": "7"})
	rr := httptest.NewRecorder()

	handler.DeleteAccount(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestDeleteAccount_Forbidden(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockBalanceUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	mockUC.EXPECT().DeleteAccount(gomock.Any(), 1, 7).Return(models.Account{}, serviceerrors.ErrForbidden)

	req := httptest.NewRequest(http.MethodDelete, "/balance/7", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	req = mux.SetURLVars(req, map[string]string{"id": "7"})
	rr := httptest.NewRecorder()

	handler.DeleteAccount(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}
