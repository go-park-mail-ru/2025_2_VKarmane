package operation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/handlers/category"
	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
	image "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/image/usecase"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	httputils "github.com/go-park-mail-ru/2025_2_VKarmane/pkg/http"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gorilla/mux"
)

type Handler struct {
	finClient     finpb.FinanceServiceClient
	imageUC       image.ImageUseCase
	kafkaProducer *kafka.Writer
	clock         clock.Clock
}

func NewHandler(finClient finpb.FinanceServiceClient, imageUC image.ImageUseCase, kafkaProducer *kafka.Writer, clck clock.Clock) *Handler {
	return &Handler{finClient: finClient, imageUC: imageUC, kafkaProducer: kafkaProducer, clock: clck}
}

func (h *Handler) getUserID(r *http.Request) (int, bool) {
	return middleware.GetUserIDFromContext(r.Context())
}

func (h *Handler) parseIDFromURL(r *http.Request, paramName string) (int, error) {
	vars := mux.Vars(r)
	idStr := vars[paramName]
	return strconv.Atoi(idStr)
}

// func OperationInListToResponse(op models.OperationInList) models.OperationInListResponse {
// 	return models.OperationInListResponse{
// 		ID:               op.ID,
// 		AccountID:        op.AccountID,
// 		CategoryID:       op.CategoryID,
// 		CategoryName:     op.CategoryName,
// 		Name:             op.Name,
// 		Type:             string(op.Type),
// 		CategoryHashedID: op.CategoryLogoHashedID,
// 		Description:      op.Description,
// 		Sum:              op.Sum,
// 		CurrencyID:       op.CurrencyID,
// 		CreatedAt:        op.CreatedAt,
// 		Date:             op.Date,
// 	}
// }

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
	id, ok := h.getUserID(r)
	if !ok {
		httputils.UnauthorizedError(w, r, "Требуется авторизация", models.ErrCodeUnauthorized)
		return
	}

	accID, err := h.parseIDFromURL(r, "acc_id")
	if err != nil {
		httputils.ValidationError(w, r, "Некорректный ID счета", "acc_id")
		return
	}

	name := r.URL.Query().Get("title")
	categoryIDStr := r.URL.Query().Get("category_id")
    var categoryID *int
    if categoryIDStr != "" {
        v, _ := strconv.Atoi(categoryIDStr)
        categoryID = &v
    }


	// ops, err := h.opUC.GetAccountOperations(r.Context(), accID)
	ops, err := h.finClient.GetOperationsByAccount(r.Context(), ProtoGetOperationsRequst(id, accID, categoryID, name))
	if err != nil {
		_, ok := status.FromError(err)
		log := logger.FromContext(r.Context())
		if !ok {
			if log != nil {
				log.Error("grpc GetOperationsByAccount unknown error", "error", err)
			}
			httputils.InternalError(w, r, "failed to get operations")
			return
		}

		if log != nil {
			log.Error("grpc GetOperationsByAccount error", "error", err)
		}
		httputils.InternalError(w, r, "failed to get operations")
		return
	}


	 operationsResponse := make([]models.OperationInListResponse, 1)
	if ops != nil {
		for _, op := range ops.Operations {
			operationsResponse = append(operationsResponse, MapOperationInListToResponse(op))
		}
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
	log := logger.FromContext(r.Context())
	id, ok := h.getUserID(r)
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

	// op, err := h.opUC.CreateOperation(r.Context(), req, accID)
	op, err := h.finClient.CreateOperation(r.Context(), CreateOperationRequestToProto(req, id, accID))
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			if log != nil {
				log.Error("grpc CreateOperation unknown error", "error", err)
			}
			httputils.InternalError(w, r, "failed to create operation")
			return
		}
		switch st.Code() {
		case codes.NotFound:
			if log != nil {
				log.Error("grpc CreateOperation exists error", "error", err)
			}
			httputils.ConflictError(w, r, "Счет не найден", models.ErrCodeAccountNotFound)
			return
		case codes.PermissionDenied:
			if log != nil {
				log.Error("grpc CreateOperation exists error", "error", err)
			}
			httputils.ConflictError(w, r, "Доступ запрещен", models.ErrCodeForbidden)
			return
		default:
			if log != nil {
				log.Error("grpc CreateOperation error", "error", err)
			}
			httputils.InternalError(w, r, "failed to create operation")
			return
		}
	}

	operationResponse := ProtoOperationToResponse(op)
	ctg, err := h.finClient.GetCategory(r.Context(), category.UserAndCtegoryIDToProto(id, operationResponse.CategoryID))
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			if log != nil {
				log.Error("grpc CreateCategory unknown error", "error", err)
			}
			httputils.InternalError(w, r, "failed to get category")
			return
		}
		switch st.Code() {
		case codes.NotFound:
			if log != nil {
				log.Error("grpc GetCategory invalid arg", "error", err)
			}
			httputils.Error(w, r, "failed to get category", http.StatusBadRequest)
			return
		default:
			if log != nil {
				log.Error("grpc GetCategory error", "error", err)
			}
			httputils.InternalError(w, r, "failed to get category")
			return
		}

	}
	ctgLogo, err := h.imageUC.GetImageURL(r.Context(), ctg.Category.LogoHashedId)
	if err != nil {
		httputils.InternalError(w, r, "Ошибка получения операций")
		return
	}
	ctgDTO := category.CategoryWithStatsToAPI(ctg)

	transactionSearch := OperationResponseToSearch(operationResponse, ctgDTO, ctgLogo)
	transactionSearch.Action = models.WRITE

	data, _ := json.Marshal(transactionSearch)
	if err = h.kafkaProducer.WriteMessages(r.Context(), kafka.Message{Value: data}); err != nil {
		httputils.InternalError(w, r, "Failed to create opeartion")
		if log != nil {
			log.Error("kafka CreateCategory unknown error", "error", err)
		}
		return
	}

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
	id, ok := h.getUserID(r)
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

	op, err := h.finClient.GetOperation(r.Context(), OperationAndUserIDToProtoID(opID, accID, id))
	if err != nil {
		st, ok := status.FromError(err)
		log := logger.FromContext(r.Context())
		if !ok {
			if log != nil {
				log.Error("grpc GetOperation unknown error", "error", err)
			}
			httputils.InternalError(w, r, "failed to get operation")
			return
		}
		switch st.Code() {
		case codes.NotFound:
			if log != nil {
				log.Error("grpc GetOperation not found error", "error", err)
			}
			httputils.ConflictError(w, r, "Операция не найдена", models.ErrCodeTransactionNotFound)
			return
		case codes.PermissionDenied:
			if log != nil {
				log.Error("grpc GetOperation forbidden", "error", err)
			}
			httputils.ConflictError(w, r, "Доступ запрещен", models.ErrCodeForbidden)
			return
		default:
			if log != nil {
				log.Error("grpc GetOperation error", "error", err)
			}
			httputils.InternalError(w, r, "failed to get operation")
			return
		}
	}

	operationResponse := ProtoOperationToResponse(op)
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
	id, ok := h.getUserID(r)
	log := logger.FromContext(r.Context())
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

	op, err := h.finClient.UpdateOperation(r.Context(), UpdateOperationRequestToProto(req, id, accID, opID))
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			if log != nil {
				log.Error("grpc UpdateOperation unknown error", "error", err)
			}
			httputils.InternalError(w, r, "failed to update operation")
			return
		}
		switch st.Code() {
		case codes.NotFound:
			if log != nil {
				log.Error("grpc UpdateOperation not found error", "error", err)
			}
			httputils.ConflictError(w, r, "Операция не найдена", models.ErrCodeTransactionNotFound)
			return
		case codes.PermissionDenied:
			if log != nil {
				log.Error("grpc UpdateOperation forbidden", "error", err)
			}
			httputils.ConflictError(w, r, "Доступ запрещен", models.ErrCodeForbidden)
			return
		default:
			if log != nil {
				log.Error("grpc UpdateOperation error", "error", err)
			}
			httputils.InternalError(w, r, "failed to update operation")
			return
		}
	}

	operationResponse := ProtoOperationToResponse(op)
	log.Info(fmt.Sprintf("%d",operationResponse.CategoryID))
	ctg, err := h.finClient.GetCategory(r.Context(), category.UserAndCtegoryIDToProto(id, operationResponse.CategoryID))
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			if log != nil {
				log.Error("grpc CreateCategory unknown error", "error", err)
			}
			httputils.InternalError(w, r, "failed update operation")
			return
		}
		switch st.Code() {
		case codes.NotFound:
			if log != nil {
				log.Error("grpc GetCategory invalid arg", "error", err)
			}
			httputils.Error(w, r, "failed to update operation", http.StatusBadRequest)
			return
		default:
			if log != nil {
				log.Error("grpc GetCategory error", "error", err)
			}
			httputils.InternalError(w, r, "failed to update operation")
			return
		}

	}

	ctgLogo, err := h.imageUC.GetImageURL(r.Context(), ctg.Category.LogoHashedId)
	if err != nil {
		httputils.InternalError(w, r, "Ошибка получения операций")
		return
	}
	ctgDTO := category.CategoryWithStatsToAPI(ctg)

	transactionSearch := OperationResponseToSearch(operationResponse, ctgDTO, ctgLogo)
	transactionSearch.Action = models.UPDATE

	data, _ := json.Marshal(transactionSearch)
	if err = h.kafkaProducer.WriteMessages(r.Context(), kafka.Message{Value: data}); err != nil {
		httputils.InternalError(w, r, "Failed to update opeartion")
		if log != nil {
			log.Error("kafka GetCategory unknown error", "error", err)
		}
		return
	}

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
	id, ok := h.getUserID(r)
	log := logger.FromContext(r.Context())
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

	op, err := h.finClient.DeleteOperation(r.Context(), OperationAndUserIDToProtoID(opID, accID, id))
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			if log != nil {
				log.Error("grpc DeleteOperation unknown error", "error", err)
			}
			httputils.InternalError(w, r, "failed to delete operation")
			return
		}
		switch st.Code() {
		case codes.NotFound:
			if log != nil {
				log.Error("grpc DeleteOperation not found error", "error", err)
			}
			httputils.ConflictError(w, r, "Операция не найдена", models.ErrCodeTransactionNotFound)
			return
		case codes.PermissionDenied:
			if log != nil {
				log.Error("grpc DeleteOperation forbidden", "error", err)
			}
			httputils.ConflictError(w, r, "Доступ запрещен", models.ErrCodeForbidden)
			return
		default:
			if log != nil {
				log.Error("grpc DeleteOperation error", "error", err)
			}
			httputils.InternalError(w, r, "failed to delete operation")
			return
		}
	}

	operationResponse := ProtoOperationToResponse(op)
	transactionSearch := OperationResponseToSearchDelete(operationResponse)
	transactionSearch.Action = models.DELETE

	data, _ := json.Marshal(transactionSearch)
	if err = h.kafkaProducer.WriteMessages(r.Context(), kafka.Message{Value: data}); err != nil {
		httputils.InternalError(w, r, "Failed to delete opeartion")
		if log != nil {
			log.Error("kafka GetCategory unknown error", "error", err)
		}
		return
	}
	httputils.Success(w, r, operationResponse)
}
