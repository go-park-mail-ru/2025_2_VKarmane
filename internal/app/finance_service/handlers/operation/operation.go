package operation

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	lg "log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/handlers/account"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/handlers/category"
	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
	image "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/image/usecase"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	kafkautils "github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/kafka"
	httputils "github.com/go-park-mail-ru/2025_2_VKarmane/pkg/http"
)

type Handler struct {
	finClient     finpb.FinanceServiceClient
	imageUC       image.ImageUseCase
	kafkaProducer kafkautils.KafkaProducer
	clock         clock.Clock
}

func NewHandler(finClient finpb.FinanceServiceClient, imageUC image.ImageUseCase, kafkaProducer kafkautils.KafkaProducer, clck clock.Clock) *Handler {
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

func (h *Handler) checkCSVHeaders(headers []string) error {
	expected := []string{"date", "category_name", "account_id", "sum-expense", "to", "sum-income", "description"}
	for i, _ := range headers {
		if headers[i] != expected[i] {
			return errors.New("invalid file headers")
		}
	}
	return nil
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

	urlQueries := r.URL.Query()

	// ops, err := h.opUC.GetAccountOperations(r.Context(), accID)
	ops, err := h.finClient.GetOperationsByAccount(r.Context(), ProtoGetOperationsRequest(id, accID, urlQueries))
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

	operationsResponse := []models.OperationInListResponse{}
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
		case codes.FailedPrecondition:
			if log != nil {
				log.Error("grpc CreateOperation exists error", "error", err)
			}
			httputils.Error(w, r, "Баланс счета не может быть отрицательным", http.StatusBadRequest)
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

	var ctgLogo string
	ctgDTO := models.CategoryWithStats{}
	if operationResponse.CategoryID != 0 {
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
		ctgLogo, err = h.imageUC.GetImageURL(r.Context(), ctg.Category.LogoHashedId)
		if err != nil {
			httputils.InternalError(w, r, "Ошибка получения операций")
			return
		}

		ctgDTO = category.CategoryWithStatsToAPI(ctg)
	}

	transactionSearch := OperationResponseToSearch(operationResponse, ctgDTO, ctgLogo)
	transactionSearch.Action = models.WRITE

	data, _ := json.Marshal(transactionSearch)
	if err = h.kafkaProducer.WriteMessages(r.Context(), kafkautils.KafkaMessage{Payload: data, Type: models.TRANSACTIONS}); err != nil {
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
		case codes.FailedPrecondition:
			if log != nil {
				log.Error("grpc CreateOperation exists error", "error", err)
			}
			httputils.Error(w, r, "Баланс счета не может быть отрицательным", http.StatusBadRequest)
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
	if err = h.kafkaProducer.WriteMessages(r.Context(), kafkautils.KafkaMessage{Payload: data, Type: models.TRANSACTIONS}); err != nil {
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
		case codes.FailedPrecondition:
			if log != nil {
				log.Error("grpc CreateOperation exists error", "error", err)
			}
			httputils.Error(w, r, "Баланс счета не может быть отрицательным", http.StatusBadRequest)
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
	if err = h.kafkaProducer.WriteMessages(r.Context(), kafkautils.KafkaMessage{Payload: data, Type: models.TRANSACTIONS}); err != nil {
		httputils.InternalError(w, r, "Failed to delete opeartion")
		if log != nil {
			log.Error("kafka GetCategory unknown error", "error", err)
		}
		return
	}
	httputils.Success(w, r, operationResponse)
}

func (h *Handler) UploadCVSData(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r)
	log := logger.FromContext(r.Context())
	if !ok {
		httputils.Error(w, r, "User not authenticated", http.StatusUnauthorized)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		httputils.Error(w, r, "Failed to upload file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		httputils.Error(w, r, "Failed to read file", http.StatusBadRequest)
		return
	}
	err = h.checkCSVHeaders(headers)
	if err != nil {
		httputils.Error(w, r, "Invalid file headers", http.StatusBadRequest)
		return
	}

	var reqs []models.CreateOperationRequest
	var accID int

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			if log != nil {
				log.Error("Failed to read file", "error", err)
			}

			errMsg := fmt.Sprintf("Failed to read file: %v", err)
			httputils.Error(w, r, errMsg, http.StatusBadRequest)
			return
		}

		var ctgID int
		categoryName := row[1]
		accountIDText := row[2]
		if categoryName != "" {
			ctg, err := h.finClient.GetCategoryByName(r.Context(), category.UserIDCategoryNameToProto(userID, categoryName))
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
				default:
					if log != nil {
						log.Error("grpc GetCategory error", "error", err)
					}
				}

			}
			ctgID = int(ctg.Category.Id)
		}

		accID, _ = strconv.Atoi(accountIDText)

		opType := "expense"
		incomeSum := row[5]
		if incomeSum != "" {
			opType = "income"
		}

		sum, err := strconv.ParseFloat(row[3], 5)
		if err != nil || sum == 0 {
			sum, _ = strconv.ParseFloat(row[5], 5)
		}

		date := row[0]
		t, err := time.Parse("2006-01-02", date)
		if err != nil {
			httputils.InternalError(w, r, "failed parse date")
			return
		}

		req := models.CreateOperationRequest{
			AccountID:   accID,
			CategoryID:  &ctgID,
			Type:        models.OperationType(opType),
			Name:        row[4],
			Description: row[6],
			Sum:         sum,
			Date:        &t,
		}
		lg.Printf("req: %v", req)
		reqs = append(reqs, req)
	}

	createdOps, _ := ProcessCSV(r.Context(), reqs, userID, accID, h.finClient)

	if len(createdOps) == 0 {
		httputils.Error(w, r, "Failed to create operations", http.StatusBadRequest)
		return
	}

	for _, op := range createdOps {
		var ctgLogo string
		var ctgDTO models.CategoryWithStats

		if op.CategoryID != 0 {
			ctg, err := h.finClient.GetCategory(r.Context(),
				category.UserAndCtegoryIDToProto(userID, op.CategoryID),
			)
			if err == nil {
				ctgLogo, _ = h.imageUC.GetImageURL(r.Context(), ctg.Category.LogoHashedId)
				ctgDTO = category.CategoryWithStatsToAPI(ctg)
			}
		}

		searchObj := OperationResponseToSearch(op, ctgDTO, ctgLogo)
		searchObj.Action = models.WRITE

		data, _ := json.Marshal(searchObj)

		_ = h.kafkaProducer.WriteMessages(
			r.Context(),
			kafkautils.KafkaMessage{
				Payload: data,
				Type:    models.TRANSACTIONS,
			},
		)
	}

	httputils.Created(w, r, "Данные успешно загружены")
}

func (h *Handler) GetCSVData(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r)
	log := logger.FromContext(r.Context())
	if !ok {
		httputils.Error(w, r, "User not authenticated", http.StatusUnauthorized)
		return
	}

	accs, err := h.finClient.GetAccountsByUser(r.Context(), account.UserIDToProtoID(userID))
	if err != nil {
		_, ok := status.FromError(err)
		if !ok {
			if log != nil {
				log.Error("grpc GetAccountsByUser unknown error", "error", err)
			}
			httputils.InternalError(w, r, "failed to get accounts")
			return
		}
		if log != nil {
			log.Error("grpc GetAccountsByUser error", "error", err)
		}
		httputils.InternalError(w, r, "failed to get account")
		return
	}

	accountsDTO := account.AccountResponseListProtoToApit(accs, userID)

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=\"operations.csv\"")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	writer.Write([]string{
		"date",
		"category_name",
		"account_id",
		"sum-expense",
		"to",
		"sum-income",
		"description",
	})

	for _, acc := range accountsDTO.Accounts {
		opsResp, err := h.finClient.GetOperationsByAccount(r.Context(), ProtoGetOperationsRequest(userID, acc.ID, r.URL.Query()))
		if err != nil {
			if log != nil {
				log.Error("grpc GetOperationsByAccount error", "error", err)
			}
			continue
		}
		for _, op := range opsResp.Operations {
			date := ""
			if op.Date != nil {
				date = op.Date.AsTime().Format("2006-01-02")
			}
			var (
				categoryName string
				sumExpense   string
				sumIncome    string
				toField      string
			)
			if op.Type == "expense" {
				categoryName = op.CategoryName
				sumExpense = fmt.Sprintf("%.2f", op.Sum)
				sumIncome = "0"
				toField = op.Name
			} else {
				categoryName = ""
				sumExpense = "0"
				sumIncome = fmt.Sprintf("%.2f", op.Sum)
				toField = op.Name
			}
			writer.Write([]string{
				date,
				categoryName,
				strconv.Itoa(int(op.AccountId)),
				sumExpense,
				toField,
				sumIncome,
				op.Description,
			})
		}
	}

}
