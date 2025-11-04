package operation

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
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetAccountOperations_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOperationUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	req := httptest.NewRequest(http.MethodGet, "/operations/account/1", nil)
	req = mux.SetURLVars(req, map[string]string{"acc_id": "1"})
	rr := httptest.NewRecorder()

	handler.GetAccountOperations(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetAccountOperations_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOperationUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	operations := []models.Operation{
		{ID: 1, AccountID: 1, Name: "Test", Sum: 100},
	}

	mockUC.EXPECT().GetAccountOperations(gomock.Any(), 1).Return(operations, nil)

	req := httptest.NewRequest(http.MethodGet, "/operations/account/1", nil)
	req = mux.SetURLVars(req, map[string]string{"acc_id": "1"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.GetAccountOperations(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestGetAccountOperations_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOperationUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	req := httptest.NewRequest(http.MethodGet, "/operations/account/invalid", nil)
	req = mux.SetURLVars(req, map[string]string{"acc_id": "invalid"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.GetAccountOperations(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateOperation_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOperationUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	createReq := models.CreateOperationRequest{
		AccountID: 1,
		Type:      models.OperationExpense,
		Name:      "Test",
		Sum:       100,
	}

	operation := models.Operation{
		ID:        1,
		AccountID: 1,
		Name:      "Test",
		Sum:       100,
		Type:      models.OperationExpense,
		Status:    models.OperationFinished,
	}

	mockUC.EXPECT().CreateOperation(gomock.Any(), gomock.Any(), 1).Return(operation, nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/operations/account/1", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"acc_id": "1"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.CreateOperation(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
}

func TestGetOperationByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOperationUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	operation := models.Operation{
		ID:        1,
		AccountID: 1,
		Name:      "Test",
		Sum:       100,
	}

	mockUC.EXPECT().GetOperationByID(gomock.Any(), 1, 1).Return(operation, nil)

	req := httptest.NewRequest(http.MethodGet, "/operations/account/1/operation/1", nil)
	req = mux.SetURLVars(req, map[string]string{"acc_id": "1", "op_id": "1"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.GetOperationByID(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateOperation_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOperationUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	name := "Updated"
	sum := 200.0
	updateReq := models.UpdateOperationRequest{
		Name: &name,
		Sum:  &sum,
	}

	operation := models.Operation{
		ID:        1,
		AccountID: 1,
		Name:      "Updated",
		Sum:       200,
	}

	mockUC.EXPECT().UpdateOperation(gomock.Any(), gomock.Any(), 1, 1).Return(operation, nil)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/operations/account/1/operation/1", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"acc_id": "1", "op_id": "1"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.UpdateOperation(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteOperation_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOperationUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	operation := models.Operation{
		ID:        1,
		AccountID: 1,
		Status:    models.OperationReverted,
	}

	mockUC.EXPECT().DeleteOperation(gomock.Any(), 1, 1).Return(operation, nil)

	req := httptest.NewRequest(http.MethodDelete, "/operations/account/1/operation/1", nil)
	req = mux.SetURLVars(req, map[string]string{"acc_id": "1", "op_id": "1"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.DeleteOperation(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestCreateOperation_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOperationUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	req := httptest.NewRequest(http.MethodPost, "/operations/account/1", nil)
	req = mux.SetURLVars(req, map[string]string{"acc_id": "1"})
	rr := httptest.NewRecorder()

	handler.CreateOperation(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetOperationByID_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOperationUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	req := httptest.NewRequest(http.MethodGet, "/operations/account/1/operation/1", nil)
	req = mux.SetURLVars(req, map[string]string{"acc_id": "1", "op_id": "1"})
	rr := httptest.NewRecorder()

	handler.GetOperationByID(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestUpdateOperation_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOperationUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	req := httptest.NewRequest(http.MethodPut, "/operations/account/1/operation/1", nil)
	req = mux.SetURLVars(req, map[string]string{"acc_id": "1", "op_id": "1"})
	rr := httptest.NewRecorder()

	handler.UpdateOperation(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestDeleteOperation_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOperationUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	req := httptest.NewRequest(http.MethodDelete, "/operations/account/1/operation/1", nil)
	req = mux.SetURLVars(req, map[string]string{"acc_id": "1", "op_id": "1"})
	rr := httptest.NewRecorder()

	handler.DeleteOperation(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestCreateOperation_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOperationUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	req := httptest.NewRequest(http.MethodPost, "/operations/account/1", bytes.NewBufferString("invalid"))
	req = mux.SetURLVars(req, map[string]string{"acc_id": "1"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.CreateOperation(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdateOperation_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOperationUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	req := httptest.NewRequest(http.MethodPut, "/operations/account/1/operation/1", bytes.NewBufferString("invalid"))
	req = mux.SetURLVars(req, map[string]string{"acc_id": "1", "op_id": "1"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.UpdateOperation(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetOperationByID_InvalidOpID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOperationUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	req := httptest.NewRequest(http.MethodGet, "/operations/account/1/operation/invalid", nil)
	req = mux.SetURLVars(req, map[string]string{"acc_id": "1", "op_id": "invalid"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.GetOperationByID(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdateOperation_InvalidAccID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOperationUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	req := httptest.NewRequest(http.MethodPut, "/operations/account/invalid/operation/1", nil)
	req = mux.SetURLVars(req, map[string]string{"acc_id": "invalid", "op_id": "1"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.UpdateOperation(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdateOperation_InvalidOpID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOperationUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	req := httptest.NewRequest(http.MethodPut, "/operations/account/1/operation/invalid", nil)
	req = mux.SetURLVars(req, map[string]string{"acc_id": "1", "op_id": "invalid"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.UpdateOperation(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDeleteOperation_InvalidAccID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOperationUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	req := httptest.NewRequest(http.MethodDelete, "/operations/account/invalid/operation/1", nil)
	req = mux.SetURLVars(req, map[string]string{"acc_id": "invalid", "op_id": "1"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.DeleteOperation(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDeleteOperation_InvalidOpID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockOperationUseCase(ctrl)
	handler := NewHandler(mockUC, clock.RealClock{})

	req := httptest.NewRequest(http.MethodDelete, "/operations/account/1/operation/invalid", nil)
	req = mux.SetURLVars(req, map[string]string{"acc_id": "1", "op_id": "invalid"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.DeleteOperation(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}
