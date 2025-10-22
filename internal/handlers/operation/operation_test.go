package operation

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	"github.com/go-park-mail-ru/2025_2_VKarmane/mocks"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountOperationsUnauthorized(t *testing.T) {
	mockUC := mocks.NewOperationUseCase(t)
	realClock := clock.RealClock{}
	h := NewHandler(mockUC, realClock)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/account/1/operations/", nil)
	rr := httptest.NewRecorder()

	h.GetAccountOperations(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetAccountOperationsSuccess(t *testing.T) {
	mockUC := mocks.NewOperationUseCase(t)
	realClock := clock.RealClock{}
	h := NewHandler(mockUC, realClock)

	mockUC.On("GetAccountOperations", mock.Anything, 7).
		Return([]models.Operation{{ID: 1, Name: "test", Sum: 100}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/account/7/operations/", nil)
	req = mux.SetURLVars(req, map[string]string{"acc_id": "7"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	h.GetAccountOperations(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	mockUC.AssertExpectations(t)
}

func TestGetOperationByIDNotFound(t *testing.T) {
	mockUC := mocks.NewOperationUseCase(t)
	realClock := clock.RealClock{}
	h := NewHandler(mockUC, realClock)

	mockUC.On("GetOperationByID", mock.Anything, 1, 7).
		Return(models.Operation{}, errors.New("not found"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/accounts/1/operations/7", nil)
	req = mux.SetURLVars(req, map[string]string{"acc_id": "1", "op_id": "7"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	h.GetOperationByID(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
	mockUC.AssertExpectations(t)
}

func TestCreateOperationSuccess(t *testing.T) {
	mockUC := mocks.NewOperationUseCase(t)
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}

	h := NewHandler(mockUC, fixedClock)

	reqBody := models.CreateOperationRequest{
		AccountID:   1,
		CategoryID:  1,
		Type:        models.OperationIncome,
		Name:        "test",
		Description: "desc",
		Sum:         100,
		CreatedAt:   h.clock.Now(),
	}
	expectedOp := models.Operation{
		ID:        1,
		Name:      reqBody.Name,
		Sum:       reqBody.Sum,
		CreatedAt: reqBody.CreatedAt,
	}

	mockUC.On("CreateOperation", mock.Anything, reqBody, 7).
		Return(expectedOp, nil)

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/accounts/7/operation/create", bytes.NewReader(bodyBytes))
	req = mux.SetURLVars(req, map[string]string{"acc_id": "7"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	h.CreateOperation(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
	mockUC.AssertExpectations(t)
}

func TestDeleteOperationSuccess(t *testing.T) {
	mockUC := mocks.NewOperationUseCase(t)
	realClock := clock.RealClock{}
	h := NewHandler(mockUC, realClock)

	mockUC.On("GetOperationByID", mock.Anything, 7, 1).
		Return(models.Operation{ID: 1, Name: "deleted"}, nil)
	mockUC.On("DeleteOperation", mock.Anything, 7, 1).
		Return(models.Operation{ID: 1, Name: "deleted"}, nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/accounts/7/operations/delete/1", nil)
	req = mux.SetURLVars(req, map[string]string{"acc_id": "7", "op_id": "1"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	h.DeleteOperation(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	mockUC.AssertExpectations(t)
}

func TestUpdateOperationUnauthorized(t *testing.T) {
	mockUC := mocks.NewOperationUseCase(t)
	realClock := clock.RealClock{}
	h := NewHandler(mockUC, realClock)

	reqBody := models.UpdateOperationRequest{Name: utilsPtr("Updated name")}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/accounts/7/operation/update/1", bytes.NewReader(bodyBytes))
	req = mux.SetURLVars(req, map[string]string{"acc_id": "7", "op_id": "1"})
	rr := httptest.NewRecorder()

	h.UpdateOperation(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestUpdateOperationNotFound(t *testing.T) {
	mockUC := mocks.NewOperationUseCase(t)
	realClock := clock.RealClock{}
	h := NewHandler(mockUC, realClock)

	mockUC.On("GetOperationByID", mock.Anything, 7, 1).
		Return(models.Operation{}, errors.New("not found"))

	reqBody := models.UpdateOperationRequest{Name: utilsPtr("Updated name")}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/accounts/7/operation/update/1", bytes.NewReader(bodyBytes))
	req = mux.SetURLVars(req, map[string]string{"acc_id": "7", "op_id": "1"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	h.UpdateOperation(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
	mockUC.AssertExpectations(t)
}

func TestUpdateOperationSuccess(t *testing.T) {
	mockUC := mocks.NewOperationUseCase(t)
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}

	h := NewHandler(mockUC, fixedClock)

	mockUC.On("GetOperationByID", mock.Anything, 7, 1).
		Return(models.Operation{ID: 1, AccountID: 7, Name: "Old name"}, nil)

	reqBody := models.UpdateOperationRequest{
		Name:      utilsPtr("Updated name"),
		CreatedAt: utilsPtr(h.clock.Now()),
	}

	mockUC.On("UpdateOperation", mock.Anything, reqBody, 7, 1).
		Return(models.Operation{ID: 1, AccountID: 7, Name: *reqBody.Name, CreatedAt: *reqBody.CreatedAt}, nil)

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/accounts/7/operation/update/1", bytes.NewReader(bodyBytes))
	req = mux.SetURLVars(req, map[string]string{"acc_id": "7", "op_id": "1"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	h.UpdateOperation(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	mockUC.AssertExpectations(t)
}

func utilsPtr[T any](v T) *T {
	return &v
}
