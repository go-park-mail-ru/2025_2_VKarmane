package account

import (
	"encoding/json"
	"net/http"
	"strconv"

	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	httputils "github.com/go-park-mail-ru/2025_2_VKarmane/pkg/http"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gorilla/mux"
)

type Handler struct {
	finClient finpb.FinanceServiceClient
	clock     clock.Clock
}

func NewHandler(finClient finpb.FinanceServiceClient, clck clock.Clock) *Handler {
	return &Handler{finClient: finClient, clock: clck}
}

func (h *Handler) getUserID(r *http.Request) (int, bool) {
	return middleware.GetUserIDFromContext(r.Context())
}

func (h *Handler) parseIDFromURL(r *http.Request, paramName string) (int, error) {
	vars := mux.Vars(r)
	idStr := vars[paramName]
	return strconv.Atoi(idStr)
}

// GetAccounts godoc
// @Summary Получение списка аккаунтов пользователя
// @Description Возвращает список всех аккаунтов пользователя в формате для фронтенда
// @Tags accounts
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "Список аккаунтов пользователя"
// @Failure 401 {object} models.ErrorResponse "Требуется аутентификация (UNAUTHORIZED, TOKEN_MISSING, TOKEN_INVALID, TOKEN_EXPIRED)"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера (INTERNAL_ERROR, DATABASE_ERROR)"
// @Router /accounts [get]
func (h *Handler) GetAccounts(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r)
	if !ok {
		httputils.Error(w, r, "User not authenticated", http.StatusUnauthorized)
		return
	}

	accounts, err := h.finClient.GetAccountsByUser(r.Context(), UserIDToProtoID(userID))
	if err != nil {
		_, ok := status.FromError(err)
		log := logger.FromContext(r.Context())
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

	accountsDTO := AccountResponseListProtoToApit(accounts, userID)
	httputils.Success(w, r, accountsDTO)
}

// GetAccountByID godoc
// @Summary Получение аккаунта по ID
// @Description Возвращает информацию об аккаунте пользователя по его ID
// @Tags accounts
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID аккаунта"
// @Success 200 {object} models.Account "Аккаунт пользователя"
// @Failure 400 {object} models.ErrorResponse "Некорректный ID аккаунта (INVALID_REQUEST)"
// @Failure 401 {object} models.ErrorResponse "Требуется аутентификация (UNAUTHORIZED, TOKEN_MISSING, TOKEN_INVALID, TOKEN_EXPIRED)"
// @Failure 404 {object} models.ErrorResponse "Аккаунт не найден (ACCOUNT_NOT_FOUND)"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера (INTERNAL_ERROR, DATABASE_ERROR)"
// @Router /balance/{id} [get]
func (h *Handler) GetAccountByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r)
	if !ok {
		httputils.Error(w, r, "User not authenticated", http.StatusUnauthorized)
		return
	}

	accountID, err := h.parseIDFromURL(r, "id")
	if err != nil {
		httputils.ValidationError(w, r, "Некорректный ID аккаунта", "id")
		return
	}

	account, err := h.finClient.GetAccount(r.Context(), UserIDAndAccountIDToProtoID(userID, accountID))
	if err != nil {
		st, ok := status.FromError(err)
		log := logger.FromContext(r.Context())
		if !ok {
			if log != nil {
				log.Error("grpc GetAccount unknown error", "error", err)
			}
			httputils.InternalError(w, r, "failed to get account")
			return
		}
		switch st.Code() {
		case codes.PermissionDenied:
			if log != nil {
				log.Error("grpc GetAccount invalid arg", "error", err)
			}
			httputils.Error(w, r, "failed to get account", http.StatusForbidden)
			return
		case codes.NotFound:
			if log != nil {
				log.Error("grpc GetAccount not found error", "error", err)
			}
			httputils.Error(w, r, "Счет не найден", http.StatusBadRequest)
			return
		default:
			if log != nil {
				log.Error("grpc GetAccount error", "error", err)
			}
			httputils.InternalError(w, r, "failed to get account")
			return
		}
	}

	// Преобразуем в формат для фронтенда (camel_case)
	accountDTO := ProtoAccountToApi(account)
	httputils.Success(w, r, accountDTO)
}

func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r)
	if !ok {
		httputils.Error(w, r, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var req models.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ValidationError(w, r, "Некорректный формат данных", "body")
		return
	}

	validationErrors := utils.ValidateStruct(req)
	if len(validationErrors) > 0 {
		httputils.ValidationErrors(w, r, validationErrors)
		return
	}

	account, err := h.finClient.CreateAccount(r.Context(), AccountCreateRequestToProto(userID, req))
	if err != nil {
		_, ok := status.FromError(err)
		log := logger.FromContext(r.Context())
		if !ok {
			if log != nil {
				log.Error("grpc CreateAccount unknown error", "error", err)
			}
			httputils.InternalError(w, r, "failed to create account")
			return
		}
		if log != nil {
			log.Error("grpc CreateAccount error", "error", err)
		}
		httputils.InternalError(w, r, "failed to create account")
		return
	}
	accountDTO := ProtoAccountToApi(account)
	httputils.Created(w, r, accountDTO)
}

func (h *Handler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r)
	if !ok {
		httputils.Error(w, r, "User not authenticated", http.StatusUnauthorized)
		return
	}

	accID, err := h.parseIDFromURL(r, "id")
	if err != nil {
		httputils.ValidationError(w, r, "Некорректный ID счета", "budgetID")
		return
	}
	var req models.UpdateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ValidationError(w, r, "Некорректный формат данных", "body")
		return
	}

	acc, err := h.finClient.UpdateAccount(r.Context(), AccountUpdateRequestToProto(userID, accID, req))
	if err != nil {
		st, ok := status.FromError(err)
		log := logger.FromContext(r.Context())
		if !ok {
			if log != nil {
				log.Error("grpc UpdateAccount unknown error", "error", err)
			}
			httputils.InternalError(w, r, "failed to update account")
			return
		}
		switch st.Code() {
		case codes.PermissionDenied:
			if log != nil {
				log.Error("grpc UpdateAccount invalid arg", "error", err)
			}
			httputils.Error(w, r, "failed to update account", http.StatusForbidden)
			return
		case codes.NotFound:
			if log != nil {
				log.Error("grpc UpdateAccount not found error", "error", err)
			}
			httputils.Error(w, r, "Счет не найден", http.StatusBadRequest)
			return
		default:
			if log != nil {
				log.Error("grpc UpdateAccount error", "error", err)
			}
			httputils.InternalError(w, r, "failed to update account")
			return
		}
	}
	accDTO := ProtoAccountToApi(acc)
	httputils.Success(w, r, accDTO)
}

func (h *Handler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r)
	if !ok {
		httputils.Error(w, r, "User not authenticated", http.StatusUnauthorized)
		return
	}

	accID, err := h.parseIDFromURL(r, "id")
	if err != nil {
		httputils.ValidationError(w, r, "Некорректный ID бюджета", "budgetID")
		return
	}

	acc, err := h.finClient.DeleteAccount(r.Context(), UserIDAndAccountIDToProtoID(userID, accID))
	if err != nil {
		st, ok := status.FromError(err)
		log := logger.FromContext(r.Context())
		if !ok {
			if log != nil {
				log.Error("grpc DeleteAccount unknown error", "error", err)
			}
			httputils.InternalError(w, r, "failed to delete account")
			return
		}
		switch st.Code() {
		case codes.PermissionDenied:
			if log != nil {
				log.Error("grpc DeleteAccount invalid arg", "error", err)
			}
			httputils.Error(w, r, "failed to delete account", http.StatusForbidden)
			return
		case codes.NotFound:
			if log != nil {
				log.Error("grpc DeleteAccount not found error", "error", err)
			}
			httputils.Error(w, r, "Счет не найден", http.StatusBadRequest)
			return
		default:
			if log != nil {
				log.Error("grpc DeleteAccount error", "error", err)
			}
			httputils.InternalError(w, r, "failed to delete account")
			return
		}
	}
	accDTO := ProtoAccountToApi(acc)
	httputils.Success(w, r, accDTO)
}
