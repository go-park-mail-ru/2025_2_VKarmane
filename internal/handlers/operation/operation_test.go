package operation

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

type MockOperationUseCase struct {
	GetOps    func(ctx context.Context, accID int) ([]models.Operation, error)
	GetOpByID func(ctx context.Context, accID, opID int) (models.Operation, error)
	CreateOp  func(ctx context.Context, req models.CreateOperationRequest, accID int) (models.Operation, error)
	UpdateOp  func(ctx context.Context, req models.UpdateOperationRequest, accID, opID int) (models.Operation, error)
	DeleteOp  func(ctx context.Context, accID, opID int) (models.Operation, error)
}

func (m *MockOperationUseCase) GetAccountOperations(ctx context.Context, accID int) ([]models.Operation, error) {
	return m.GetOps(ctx, accID)
}
func (m *MockOperationUseCase) GetOperationByID(ctx context.Context, accID, opID int) (models.Operation, error) {
	return m.GetOpByID(ctx, accID, opID)
}
func (m *MockOperationUseCase) CreateOperation(ctx context.Context, req models.CreateOperationRequest, accID int) (models.Operation, error) {
	return m.CreateOp(ctx, req, accID)
}
func (m *MockOperationUseCase) UpdateOperation(ctx context.Context, req models.UpdateOperationRequest, accID, opID int) (models.Operation, error) {
	return m.UpdateOp(ctx, req, accID, opID)
}
func (m *MockOperationUseCase) DeleteOperation(ctx context.Context, accID, opID int) (models.Operation, error) {
	return m.DeleteOp(ctx, accID, opID)
}


func TestGetAccountOperationsUnauthorized(t *testing.T) {
	mockUC := &MockOperationUseCase{}
	h := NewHandler(mockUC)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/account/1/operations/", nil)
	rr := httptest.NewRecorder()

	h.GetAccountOperations(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetAccountOperationsSuccess(t *testing.T) {
	mockUC := &MockOperationUseCase{
		GetOps: func(ctx context.Context, accID int) ([]models.Operation, error) {
			return []models.Operation{{ID: 1, Name: "test", Sum: 100}}, nil
		},
	}
	h := NewHandler(mockUC)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/account/7/operations/", nil)
	req = mux.SetURLVars(req, map[string]string{"acc_id": "7"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	h.GetAccountOperations(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestGetOperationByIDNotFound(t *testing.T) {
	mockUC := &MockOperationUseCase{
		GetOpByID: func(ctx context.Context, accID, opID int) (models.Operation, error) {
			return models.Operation{}, errors.New("not found")
		},
	}
	h := NewHandler(mockUC)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/accounts/1/operations/7", nil)
	req = mux.SetURLVars(req, map[string]string{"acc_id": "1", "op_id": "7"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	h.GetOperationByID(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestCreateOperationSuccess(t *testing.T) {
	mockUC := &MockOperationUseCase{
		CreateOp: func(ctx context.Context, req models.CreateOperationRequest, accID int) (models.Operation, error) {
			return models.Operation{ID: 1, Name: req.Name, Sum: req.Sum}, nil
		},
	}
	h := NewHandler(mockUC)

	reqBody := models.CreateOperationRequest{
		AccountID:  1,
		CategoryID: 1,
		Type:       models.OperationIncome,
		Name:       "test",
		Description: "desc",
		Sum:        100,
	}
	bodyBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/accounts/7/operation/create", bytes.NewReader(bodyBytes))
	req = mux.SetURLVars(req, map[string]string{"acc_id": "7"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	h.CreateOperation(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
}

func TestDeleteOperationSuccess(t *testing.T) {
	mockUC := &MockOperationUseCase{
		GetOpByID: func(ctx context.Context, accID, opID int) (models.Operation, error) {
			return models.Operation{ID: opID, Name: "deleted"}, nil
		},
		DeleteOp: func(ctx context.Context, accID, opID int) (models.Operation, error) {
			return models.Operation{ID: opID, Name: "deleted"}, nil
		},
	}
	h := NewHandler(mockUC)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/accounts/7/operations/delete/1", nil)
	req = mux.SetURLVars(req, map[string]string{"acc_id": "7", "op_id": "1"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	h.DeleteOperation(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateOperationUnauthorized(t *testing.T) {
	mockUC := &MockOperationUseCase{}
	h := NewHandler(mockUC)

	reqBody := models.UpdateOperationRequest{Name: utilsPtr("Updated name")}
	bodyBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/accounts/7/operation/update/1", bytes.NewReader(bodyBytes))
	req = mux.SetURLVars(req, map[string]string{"acc_id": "7", "op_id": "1"})
	rr := httptest.NewRecorder()

	h.UpdateOperation(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestUpdateOperationNotFound(t *testing.T) {
	mockUC := &MockOperationUseCase{
		GetOpByID: func(ctx context.Context, accID, opID int) (models.Operation, error) {
			return models.Operation{}, errors.New("not found")
		},
	}
	h := NewHandler(mockUC)

	reqBody := models.UpdateOperationRequest{Name: utilsPtr("Updated name")}
	bodyBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/accounts/7/operation/update/1", bytes.NewReader(bodyBytes))
	req = mux.SetURLVars(req, map[string]string{"acc_id": "7", "op_id": "1"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	h.UpdateOperation(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestUpdateOperationSuccess(t *testing.T) {
	mockUC := &MockOperationUseCase{
		GetOpByID: func(ctx context.Context, accID, opID int) (models.Operation, error) {
			return models.Operation{ID: opID, AccountID: accID, Name: "Old name"}, nil
		},
		UpdateOp: func(ctx context.Context, req models.UpdateOperationRequest, accID, opID int) (models.Operation, error) {
			return models.Operation{ID: opID, AccountID: accID, Name: *req.Name}, nil
		},
	}
	h := NewHandler(mockUC)

	reqBody := models.UpdateOperationRequest{Name: utilsPtr("Updated name")}
	bodyBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/accounts/7/operation/update/1", bytes.NewReader(bodyBytes))
	req = mux.SetURLVars(req, map[string]string{"acc_id": "7", "op_id": "1"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	h.UpdateOperation(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func utilsPtr[T any](v T) *T {
	return &v
}
