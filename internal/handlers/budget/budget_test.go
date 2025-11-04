package budget

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetListBudgetsUnauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockBudgetUseCase(ctrl)
	clock := clock.RealClock{}
	h := NewHandler(m, clock)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/budgets", nil)
	rr := httptest.NewRecorder()
	h.GetListBudgets(rr, req)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetBudgetByIDSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockBudgetUseCase(ctrl)
	clock := clock.RealClock{}
	h := NewHandler(m, clock)

	m.EXPECT().GetBudgetByID(gomock.Any(), 1, 5).Return(models.Budget{ID: 5}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/budget/5", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "5"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()
	h.GetBudgetByID(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
}

func TestGetListBudgets_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockBudgetUseCase(ctrl)
	clock := clock.RealClock{}
	h := NewHandler(m, clock)

	budgets := []models.Budget{{ID: 1}, {ID: 2}}
	m.EXPECT().GetBudgetsForUser(gomock.Any(), 1).Return(budgets, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/budgets", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()
	h.GetListBudgets(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
}

func TestGetBudgetByID_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockBudgetUseCase(ctrl)
	clock := clock.RealClock{}
	h := NewHandler(m, clock)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/budget/invalid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()
	h.GetBudgetByID(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)
}
