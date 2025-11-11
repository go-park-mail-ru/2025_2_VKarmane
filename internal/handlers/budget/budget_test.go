package budget

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
	budgeterrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/budget"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetListBudgets_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockBudgetUseCase(ctrl)
	h := NewHandler(mockUC, clock.RealClock{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/budgets", nil)
	rr := httptest.NewRecorder()

	h.GetListBudgets(rr, req)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetListBudgets_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockBudgetUseCase(ctrl)
	h := NewHandler(mockUC, clock.RealClock{})

	budgets := []models.Budget{{ID: 1}, {ID: 2}}
	mockUC.EXPECT().GetBudgetsForUser(gomock.Any(), 42).Return(budgets, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/budgets", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 42))
	rr := httptest.NewRecorder()

	h.GetListBudgets(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
}

func TestGetBudgetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockBudgetUseCase(ctrl)
	h := NewHandler(mockUC, clock.RealClock{})

	mockUC.EXPECT().GetBudgetByID(gomock.Any(), 1, 5).Return(models.Budget{ID: 5}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/budget/5", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "5"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	h.GetBudgetByID(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
}

func TestGetBudgetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockBudgetUseCase(ctrl)
	h := NewHandler(mockUC, clock.RealClock{})

	mockUC.EXPECT().GetBudgetByID(gomock.Any(), 1, 5).Return(models.Budget{}, budgeterrors.ErrBudgetNotFound)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/budget/5", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "5"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	h.GetBudgetByID(rr, req)
	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestCreateBudget_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockBudgetUseCase(ctrl)
	h := NewHandler(mockUC, clock.RealClock{})

	reqBody := models.CreateBudgetRequest{
		CategoryID: 1,
		Description: "Test budget",
		Amount: 100,
	}
	budget := models.Budget{ID: 10, Amount: 100, CategoryID: 1}

	mockUC.EXPECT().CreateBudget(gomock.Any(), reqBody, 42).Return(budget, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/budget", bytes.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 42))
	rr := httptest.NewRecorder()

	h.CreateBudget(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
}

func TestCreateBudget_Conflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockBudgetUseCase(ctrl)
	h := NewHandler(mockUC, clock.RealClock{})

	reqBody := models.CreateBudgetRequest{CategoryID: 1, Amount: 50}
	mockUC.EXPECT().CreateBudget(gomock.Any(), reqBody, 1).Return(models.Budget{}, budgeterrors.ErrActiveBudgetExists)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/budget", bytes.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	h.CreateBudget(rr, req)
	require.Equal(t, http.StatusConflict, rr.Code)
}

func TestUpdateBudget_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockBudgetUseCase(ctrl)
	h := NewHandler(mockUC, clock.RealClock{})

	sum  := 200.5
	reqBody := models.UpdatedBudgetRequest{Amount: &sum}
	updatedBudget := models.Budget{ID: 5, Amount: sum}

	mockUC.EXPECT().UpdateBudget(gomock.Any(), reqBody, 1, 5).Return(updatedBudget, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/budget/5", bytes.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": "5"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	h.UpdateBudget(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteBudget_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockBudgetUseCase(ctrl)
	h := NewHandler(mockUC, clock.RealClock{})

	mockUC.EXPECT().DeleteBudget(gomock.Any(), 1, 5).Return(models.Budget{ID: 5}, nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/budget/5", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "5"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	h.DeleteBudget(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteBudget_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockBudgetUseCase(ctrl)
	h := NewHandler(mockUC, clock.RealClock{})

	mockUC.EXPECT().DeleteBudget(gomock.Any(), 1, 5).Return(models.Budget{}, budgeterrors.ErrBudgetNotFound)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/budget/5", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "5"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	h.DeleteBudget(rr, req)
	require.Equal(t, http.StatusNotFound, rr.Code)
}
