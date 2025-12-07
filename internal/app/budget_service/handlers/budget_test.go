package budget

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	bdgpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/proto"
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

func TestGetListBudgets_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockBudgetServiceClient(ctrl)
	h := NewHandler(clock.RealClock{}, mockClient)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/budgets", nil)
	rr := httptest.NewRecorder()

	h.GetListBudgets(rr, req)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetListBudgets_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockBudgetServiceClient(ctrl)
	h := NewHandler(clock.RealClock{}, mockClient)

	budgetsResp := &bdgpb.ListBudgetsResponse{
		Budgets: []*bdgpb.Budget{
			{Id: 1},
			{Id: 2},
		},
	}

	mockClient.EXPECT().
		GetListBudgets(gomock.Any(), &bdgpb.UserID{UserID: 42}).
		Return(budgetsResp, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/budgets", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 42))
	rr := httptest.NewRecorder()

	h.GetListBudgets(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)

	var resp map[string][]models.Budget
	err := json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err)
	require.Contains(t, resp, "budgets")
	require.Len(t, resp["budgets"], 2)
}

func TestGetBudgetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockBudgetServiceClient(ctrl)
	h := NewHandler(clock.RealClock{}, mockClient)

	mockClient.EXPECT().
		GetBudget(gomock.Any(), &bdgpb.BudgetRequest{UserID: 1, BudgetID: 5}).
		Return(&bdgpb.Budget{Id: 5}, nil)

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

	mockClient := mocks.NewMockBudgetServiceClient(ctrl)
	h := NewHandler(clock.RealClock{}, mockClient)

	mockClient.EXPECT().
		GetBudget(gomock.Any(), &bdgpb.BudgetRequest{UserID: 1, BudgetID: 5}).
		Return(nil, status.Error(codes.NotFound, "not found"))

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

	mockClient := mocks.NewMockBudgetServiceClient(ctrl)
	h := NewHandler(clock.RealClock{}, mockClient)

	reqBody := models.CreateBudgetRequest{
		Description: "Test budget",
		Amount:      100,
		CreatedAt:   time.Now(),
		PeriodStart: time.Now(),
		PeriodEnd:   time.Now().Add(30 * 24 * time.Hour),
	}

	mockClient.EXPECT().
		CreateBudget(gomock.Any(), ModelCreateReqtoProtoReq(reqBody, 42)).
		Return(&bdgpb.Budget{Id: 10, Sum: 100}, nil)

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

	mockClient := mocks.NewMockBudgetServiceClient(ctrl)
	h := NewHandler(clock.RealClock{}, mockClient)

	reqBody := models.CreateBudgetRequest{Amount: 50}

	mockClient.EXPECT().
		CreateBudget(gomock.Any(), ModelCreateReqtoProtoReq(reqBody, 1)).
		Return(nil, status.Error(codes.AlreadyExists, "active budget exists"))

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

	mockClient := mocks.NewMockBudgetServiceClient(ctrl)
	h := NewHandler(clock.RealClock{}, mockClient)

	sum := 200.5
	reqBody := models.UpdatedBudgetRequest{Amount: &sum, PeriodStart: &time.Time{}, PeriodEnd: &time.Time{}}

	mockClient.EXPECT().
		UpdateBudget(gomock.Any(), ModelUpdateReqtoProtoReq(reqBody, 5, 1)).
		Return(&bdgpb.Budget{Id: 5, Sum: sum}, nil)

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

	mockClient := mocks.NewMockBudgetServiceClient(ctrl)
	h := NewHandler(clock.RealClock{}, mockClient)

	mockClient.EXPECT().
		DeleteBudget(gomock.Any(), IDsToBudgetRequest(5, 1)).
		Return(&bdgpb.Budget{Id: 5}, nil)

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

	mockClient := mocks.NewMockBudgetServiceClient(ctrl)
	h := NewHandler(clock.RealClock{}, mockClient)

	mockClient.EXPECT().
		DeleteBudget(gomock.Any(), IDsToBudgetRequest(5, 1)).
		Return(nil, status.Error(codes.NotFound, "not found"))

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/budget/5", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "5"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	h.DeleteBudget(rr, req)
	require.Equal(t, http.StatusNotFound, rr.Code)
}
