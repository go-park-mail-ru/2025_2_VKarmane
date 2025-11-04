package balance

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
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
