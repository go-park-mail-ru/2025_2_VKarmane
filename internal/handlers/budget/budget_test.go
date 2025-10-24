package budget

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	"github.com/go-park-mail-ru/2025_2_VKarmane/mocks"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetListBudgetsUnauthorized(t *testing.T) {
	m := mocks.NewBudgetUseCase(t)
	clock := clock.RealClock{}
	h := NewHandler(m, clock)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/budgets", nil)
	rr := httptest.NewRecorder()
	h.GetListBudgets(rr, req)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetBudgetByIDSuccess(t *testing.T) {
	m := mocks.NewBudgetUseCase(t)
	clock := clock.RealClock{}
	h := NewHandler(m, clock)

	m.On("GetBudgetByID", mock.Anything, 1, 5).Return(models.Budget{ID: 5}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/budget/5", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "5"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()
	h.GetBudgetByID(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
}
