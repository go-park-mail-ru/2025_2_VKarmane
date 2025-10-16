package operation

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/operation"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
	httputils "github.com/go-park-mail-ru/2025_2_VKarmane/pkg/http"
	"github.com/gorilla/mux"
)

type Handler struct {
	opUC OperationUseCase
}

func NewHandler(opUC OperationUseCase) *Handler {
	return &Handler{opUC: opUC}
}

func (h *Handler) getUserID(r *http.Request) (int, bool) {
	return middleware.GetUserIDFromContext(r.Context())
}

func (h *Handler) parseIDFromURL(r *http.Request, paramName string) (int, error) {
	vars := mux.Vars(r)
	idStr := vars[paramName]
	return strconv.Atoi((idStr))
}

func (h *Handler) GetAccountOperations(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r)
	if !ok {
		httputils.Error(w, r, "User not authenticated", http.StatusUnauthorized)
		return
	}

	accID, err := h.parseIDFromURL(r, "acc_id")
	if err != nil {
		httputils.ValidationError(w, r, "Invalid account ID format", "id")
		return
	}

	ops, err := h.opUC.GetAccountOperations(r.Context(), accID)
	if err != nil {
		switch {
		case errors.Is(err, operation.ErrForbidden):
			httputils.Error(w, r, "У пользователя нет доступа", http.StatusForbidden)
		default:
			httputils.InternalError(w, r, "Failed to get operations for user")
		}
		return
	}

	opsDTO := OperationsToApi(userID, ops)
	httputils.Success(w, r, opsDTO)
}

func (h *Handler) GetOperationByID(w http.ResponseWriter, r *http.Request) {
	opID, err := h.parseIDFromURL(r, "op_id")
	if err != nil {
		httputils.ValidationError(w, r, "Invalid operation ID format", "id")
		return
	}
	accID, err := h.parseIDFromURL(r, "acc_id")
	if err != nil {
		httputils.ValidationError(w, r, "Invalid acc ID format", "id")
		return
	}

	if _, ok := h.getUserID(r); !ok {
		httputils.Error(w, r, "User not authenticated", http.StatusUnauthorized)
		return
	}

	op, err := h.opUC.GetOperationByID(r.Context(), accID, opID)
	if err != nil {
		httputils.NotFoundError(w, r, "Operation not found")
		return
	}

	transactionDTO := OperationToApi(op)
	httputils.Success(w, r, transactionDTO)
}

func (h *Handler) UpdateOperation(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateOperationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ValidationError(w, r, "Некорректный формат данных", "body")
		return
	}

	validationErrors := utils.ValidateStruct(req)
	if len(validationErrors) > 0 {
		httputils.ValidationErrors(w, r, validationErrors)
		return
	}

	opID, err := h.parseIDFromURL(r, "op_id")
	if err != nil {
		httputils.ValidationError(w, r, "Invalid operation ID format", "id")
		return
	}
	accID, err := h.parseIDFromURL(r, "acc_id")
	if err != nil {
		httputils.ValidationError(w, r, "Invalid acc ID format", "id")
		return
	}

	if _, ok := h.getUserID(r); !ok {
		httputils.Error(w, r, "User not authenticated", http.StatusUnauthorized)
		return
	}

	if _, err = h.opUC.GetOperationByID(r.Context(), accID, opID); err != nil {
		httputils.NotFoundError(w, r, "Operation not found")
		return
	}

	response, err := h.opUC.UpdateOperation(r.Context(), req, accID, opID)
	if err != nil {
		switch {
		case errors.Is(err, operation.ErrForbidden):
			httputils.Error(w, r, "У пользователя нет доступа", http.StatusForbidden)
		default:
			httputils.Error(w, r, "Failed to update operation", http.StatusBadRequest)
		}
		return
	}

	opDTO := OperationToApi(response)
	httputils.Success(w, r, opDTO)
}

func (h *Handler) DeleteOperation(w http.ResponseWriter, r *http.Request) {
	opID, err := h.parseIDFromURL(r, "op_id")
	if err != nil {
		httputils.ValidationError(w, r, "Invalid operation ID format", "id")
		return
	}
	accID, err := h.parseIDFromURL(r, "acc_id")
	if err != nil {
		httputils.ValidationError(w, r, "Invalid acc ID format", "id")
		return
	}

	if _, ok := h.getUserID(r); !ok {
		httputils.Error(w, r, "User not authenticated", http.StatusUnauthorized)
		return
	}

	if _, err = h.opUC.GetOperationByID(r.Context(), accID, opID); err != nil {
		httputils.NotFoundError(w, r, "Operation not found")
		return
	}

	response, err := h.opUC.DeleteOperation(r.Context(), accID, opID)
	if err != nil {
		switch {
		case errors.Is(err, operation.ErrForbidden):
			httputils.Error(w, r, "У пользователя нет доступа", http.StatusForbidden)
		default:
			httputils.Error(w, r, "Failed to delete operation", http.StatusBadRequest)
		}
		return
	}

	opDTO := OperationToApi(response)
	httputils.Success(w, r, opDTO)
}

func (h *Handler) CreateOperation(w http.ResponseWriter, r *http.Request) {

	var req models.CreateOperationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ValidationError(w, r, "Некорректный формат данных", "body")
		return
	}

	validationErrors := utils.ValidateStruct(req)
	if len(validationErrors) > 0 {
		httputils.ValidationErrors(w, r, validationErrors)
		return
	}

	if _, ok := h.getUserID(r); !ok {
		httputils.Error(w, r, "User not authenticated", http.StatusUnauthorized)
		return
	}
	accID, err := h.parseIDFromURL(r, "acc_id")
	if err != nil {
		httputils.ValidationError(w, r, "Invalid acc ID format", "id")
		return
	}

	response, err := h.opUC.CreateOperation(r.Context(), req, accID)
	if err != nil {
		switch {
		case errors.Is(err, operation.ErrForbidden):
			httputils.Error(w, r, "У пользователя нет доступа", http.StatusForbidden)
		default:
			httputils.Error(w, r, "Failed to create operation", http.StatusBadRequest)
		}
		return
	}

	opDTO := OperationToApi(response)
	httputils.Created(w, r, opDTO)
}
