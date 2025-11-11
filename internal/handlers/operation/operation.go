package operation

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/image"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	httputils "github.com/go-park-mail-ru/2025_2_VKarmane/pkg/http"
	serviceerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/errors"

	"github.com/gorilla/mux"
)

type Handler struct {
	opUC    OperationUseCase
	imageUC image.ImageUseCase
	clock   clock.Clock
}

func NewHandler(opUC OperationUseCase, imageUC image.ImageUseCase, clck clock.Clock) *Handler {
	return &Handler{opUC: opUC, imageUC: imageUC, clock: clck}
}

func (h *Handler) getUserID(r *http.Request) (int, bool) {
	return middleware.GetUserIDFromContext(r.Context())
}

func (h *Handler) parseIDFromURL(r *http.Request, paramName string) (int, error) {
	vars := mux.Vars(r)
	idStr := vars[paramName]
	return strconv.Atoi(idStr)
}

func OperationToResponse(op models.Operation) models.OperationResponse {
	return models.OperationResponse{
		ID:           op.ID,
		AccountID:    op.AccountID,
		CategoryID:   op.CategoryID,
		CategoryName: op.CategoryName,
		Type:         string(op.Type),
		Status:       string(op.Status),
		Description:  op.Description,
		ReceiptURL:   op.ReceiptURL,
		Name:         op.Name,
		Sum:          op.Sum,
		CurrencyID:   op.CurrencyID,
		CreatedAt:    op.CreatedAt,
		Date:         op.Date,
	}
}

func OperationInListToResponse(op models.OperationInList) models.OperationInListResponse {
	return models.OperationInListResponse{
		ID:               op.ID,
		AccountID:        op.AccountID,
		CategoryID:       op.CategoryID,
		CategoryName:     op.CategoryName,
		Name:             op.Name,
		Type:             string(op.Type),
		CategoryHashedID: op.CategoryLogoHashedID,
		Description:      op.Description,
		Sum:              op.Sum,
		CurrencyID:       op.CurrencyID,
		CreatedAt:        op.CreatedAt,
		Date:             op.Date,
	}
}

// GetAccountOperations godoc
// @Summary Получение операций по счету
// @Description Возвращает список всех операций для указанного счета
// @Tags operations
// @Produce json
// @Security ApiKeyAuth
// @Param acc_id path int true "ID счета"
// @Success 200 {object} map[string]interface{} "Список операций"
// @Failure 400 {object} models.ErrorResponse "Некорректный ID счета (INVALID_REQUEST)"
// @Failure 401 {object} models.ErrorResponse "Требуется аутентификация (UNAUTHORIZED, TOKEN_MISSING, TOKEN_INVALID, TOKEN_EXPIRED)"
// @Failure 403 {object} models.ErrorResponse "Доступ запрещен (FORBIDDEN, ACCESS_DENIED)"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера (INTERNAL_ERROR, DATABASE_ERROR)"
// @Router /operations/account/{acc_id} [get]
func (h *Handler) GetAccountOperations(w http.ResponseWriter, r *http.Request) {
	_, ok := h.getUserID(r)
	if !ok {
		httputils.UnauthorizedError(w, r, "Требуется авторизация", models.ErrCodeUnauthorized)
		return
	}

	accID, err := h.parseIDFromURL(r, "acc_id")
	if err != nil {
		httputils.ValidationError(w, r, "Некорректный ID счета", "acc_id")
		return
	}

	ops, err := h.opUC.GetAccountOperations(r.Context(), accID)
	if err != nil {
		httputils.InternalError(w, r, "Ошибка получения операций")
		return
	}

	// Преобразуем операции в OperationResponse
	var operationsResponse []models.OperationInListResponse
	if ops != nil {
		for _, op := range ops {
			operationsResponse = append(operationsResponse, OperationInListToResponse(op))
		}
	}

	if operationsResponse == nil {
		operationsResponse = []models.OperationInListResponse{}
	}

	for idx, op := range operationsResponse {
		if op.CategoryHashedID != "" {
			opCategoryLogo, err := h.imageUC.GetImageURL(r.Context(), operationsResponse[idx].CategoryHashedID)
			if err != nil {
				httputils.InternalError(w, r, "Ошибка получения операций")
				return
			}
			operationsResponse[idx].CategoryLogo = opCategoryLogo
		}
	}

	response := map[string]interface{}{
		"operations": operationsResponse,
	}
	httputils.Success(w, r, response)
}

// CreateOperation godoc
// @Summary Создание новой операции
// @Description Создает новую операцию для указанного счета
// @Tags operations
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param acc_id path int true "ID счета"
// @Param request body models.CreateOperationRequest true "Данные для создания операции"
// @Success 201 {object} map[string]interface{} "Операция успешно создана"
// @Failure 400 {object} models.ErrorResponse "Некорректные данные (INVALID_REQUEST, MISSING_FIELDS, INVALID_AMOUNT)"
// @Failure 401 {object} models.ErrorResponse "Требуется аутентификация (UNAUTHORIZED, TOKEN_MISSING, TOKEN_INVALID, TOKEN_EXPIRED)"
// @Failure 403 {object} models.ErrorResponse "Доступ запрещен (FORBIDDEN, ACCESS_DENIED)"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера (INTERNAL_ERROR, DATABASE_ERROR)"
// @Router /operations/account/{acc_id} [post]
func (h *Handler) CreateOperation(w http.ResponseWriter, r *http.Request) {
	_, ok := h.getUserID(r)
	if !ok {
		httputils.UnauthorizedError(w, r, "Требуется авторизация", models.ErrCodeUnauthorized)
		return
	}

	accID, err := h.parseIDFromURL(r, "acc_id")
	if err != nil {
		httputils.ValidationError(w, r, "Некорректный ID счета", "acc_id")
		return
	}

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

	op, err := h.opUC.CreateOperation(r.Context(), req, accID)
	if err != nil {
		if errors.Is(err, serviceerrors.ErrForbidden) {
			httputils.Error(w, r, "Доступ к счету запрещен", 403)
			return
		}
		httputils.InternalError(w, r, "Ошибка создания операции")
		return
	}

	operationResponse := OperationToResponse(op)
	httputils.Created(w, r, operationResponse)
}

// GetOperationByID godoc
// @Summary Получение операции по ID
// @Description Возвращает информацию об операции по её ID
// @Tags operations
// @Produce json
// @Security ApiKeyAuth
// @Param acc_id path int true "ID счета"
// @Param op_id path int true "ID операции"
// @Success 200 {object} models.OperationResponse "Операция"
// @Failure 400 {object} models.ErrorResponse "Некорректный ID (INVALID_REQUEST)"
// @Failure 401 {object} models.ErrorResponse "Требуется аутентификация (UNAUTHORIZED, TOKEN_MISSING, TOKEN_INVALID, TOKEN_EXPIRED)"
// @Failure 403 {object} models.ErrorResponse "Доступ запрещен (FORBIDDEN, ACCESS_DENIED)"
// @Failure 404 {object} models.ErrorResponse "Операция не найдена (NOT_FOUND)"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера (INTERNAL_ERROR, DATABASE_ERROR)"
// @Router /operations/account/{acc_id}/operation/{op_id} [get]
func (h *Handler) GetOperationByID(w http.ResponseWriter, r *http.Request) {
	_, ok := h.getUserID(r)
	if !ok {
		httputils.UnauthorizedError(w, r, "Требуется авторизация", models.ErrCodeUnauthorized)
		return
	}

	accID, err := h.parseIDFromURL(r, "acc_id")
	if err != nil {
		httputils.ValidationError(w, r, "Некорректный ID счета", "acc_id")
		return
	}

	opID, err := h.parseIDFromURL(r, "op_id")
	if err != nil {
		httputils.ValidationError(w, r, "Некорректный ID операции", "op_id")
		return
	}

	op, err := h.opUC.GetOperationByID(r.Context(), accID, opID)
	if err != nil {
		if errors.Is(err, serviceerrors.ErrForbidden) {

			httputils.Error(w, r, "Доступ к операции запрещен", 403)
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			httputils.Error(w, r, "Операция не найдена", 404)
			return
		}
		httputils.InternalError(w, r, "Ошибка получения операции")
		return
	}

	operationResponse := OperationToResponse(op)
	httputils.Success(w, r, operationResponse)
}

// UpdateOperation godoc
// @Summary Изменение операции
// @Description Обновляет информацию об операции
// @Tags operations
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param acc_id path int true "ID счета"
// @Param op_id path int true "ID операции"
// @Param request body models.UpdateOperationRequest true "Данные для обновления операции"
// @Success 200 {object} models.OperationResponse "Операция успешно обновлена"
// @Failure 400 {object} models.ErrorResponse "Некорректные данные (INVALID_REQUEST, INVALID_AMOUNT)"
// @Failure 401 {object} models.ErrorResponse "Требуется аутентификация (UNAUTHORIZED, TOKEN_MISSING, TOKEN_INVALID, TOKEN_EXPIRED)"
// @Failure 403 {object} models.ErrorResponse "Доступ запрещен (FORBIDDEN, ACCESS_DENIED)"
// @Failure 404 {object} models.ErrorResponse "Операция не найдена (NOT_FOUND)"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера (INTERNAL_ERROR, DATABASE_ERROR)"
// @Router /operations/account/{acc_id}/operation/{op_id} [put]
func (h *Handler) UpdateOperation(w http.ResponseWriter, r *http.Request) {
	_, ok := h.getUserID(r)
	if !ok {
		httputils.UnauthorizedError(w, r, "Требуется авторизация", models.ErrCodeUnauthorized)
		return
	}

	accID, err := h.parseIDFromURL(r, "acc_id")
	if err != nil {
		httputils.ValidationError(w, r, "Некорректный ID счета", "acc_id")
		return
	}

	opID, err := h.parseIDFromURL(r, "op_id")
	if err != nil {
		httputils.ValidationError(w, r, "Некорректный ID операции", "op_id")
		return
	}

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

	if err := req.Validate(); err != nil {
		httputils.ValidationError(w, r, err.Error(), "sum")
		return
	}

	op, err := h.opUC.UpdateOperation(r.Context(), req, accID, opID)
	if err != nil {
		if errors.Is(err, serviceerrors.ErrForbidden) {
			httputils.Error(w, r, "Доступ к операции запрещен", 403)
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			httputils.Error(w, r, "Операция не найдена", 404)
			return
		}
		httputils.InternalError(w, r, "Ошибка обновления операции")
		return
	}

	operationResponse := OperationToResponse(op)
	httputils.Success(w, r, operationResponse)
}

// DeleteOperation godoc
// @Summary Удаление операции
// @Description Помечает операцию как удаленную (reverted)
// @Tags operations
// @Produce json
// @Security ApiKeyAuth
// @Param acc_id path int true "ID счета"
// @Param op_id path int true "ID операции"
// @Success 200 {object} models.OperationResponse "Операция успешно удалена"
// @Failure 400 {object} models.ErrorResponse "Некорректный ID (INVALID_REQUEST)"
// @Failure 401 {object} models.ErrorResponse "Требуется аутентификация (UNAUTHORIZED, TOKEN_MISSING, TOKEN_INVALID, TOKEN_EXPIRED)"
// @Failure 403 {object} models.ErrorResponse "Доступ запрещен (FORBIDDEN, ACCESS_DENIED)"
// @Failure 404 {object} models.ErrorResponse "Операция не найдена (NOT_FOUND)"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера (INTERNAL_ERROR, DATABASE_ERROR)"
// @Router /operations/account/{acc_id}/operation/{op_id} [delete]
func (h *Handler) DeleteOperation(w http.ResponseWriter, r *http.Request) {
	_, ok := h.getUserID(r)
	if !ok {
		httputils.UnauthorizedError(w, r, "Требуется авторизация", models.ErrCodeUnauthorized)
		return
	}

	accID, err := h.parseIDFromURL(r, "acc_id")
	if err != nil {
		httputils.ValidationError(w, r, "Некорректный ID счета", "acc_id")
		return
	}

	opID, err := h.parseIDFromURL(r, "op_id")
	if err != nil {
		httputils.ValidationError(w, r, "Некорректный ID операции", "op_id")
		return
	}

	op, err := h.opUC.DeleteOperation(r.Context(), accID, opID)
	if err != nil {
		if errors.Is(err, serviceerrors.ErrForbidden) {
			httputils.Error(w, r, "Доступ к операции запрещен", 403)
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			httputils.Error(w, r, "Операция не найдена", 404)
			return
		}
		httputils.InternalError(w, r, "Ошибка удаления операции")
		return
	}

	operationResponse := OperationToResponse(op)
	httputils.Success(w, r, operationResponse)
}
