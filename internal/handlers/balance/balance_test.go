package balance

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/mocks"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetListBalanceUnauthorized(t *testing.T) {
	m := mocks.NewBalanceUseCase(t)
	h := NewHandler(m)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/balance", nil)
	rr := httptest.NewRecorder()
	h.GetListBalance(rr, req)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetBalanceByAccountIDSuccess(t *testing.T) {
	m := mocks.NewBalanceUseCase(t)
	h := NewHandler(m)

	m.On("GetAccountByID", mock.Anything, 1, 7).Return(models.Account{ID: 7}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/balance/7", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "7"})
	rr := httptest.NewRecorder()

	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	h.GetBalanceByAccountID(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
}
