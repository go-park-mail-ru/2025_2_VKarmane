package balance

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	// accounterrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/account"
	serviceerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/errors"
	httputils "github.com/go-park-mail-ru/2025_2_VKarmane/pkg/http"
)

type Handler struct {
	balanceUC BalanceUseCase
	clock     clock.Clock
}

func NewHandler(balanceUC BalanceUseCase, clck clock.Clock) *Handler {
	return &Handler{balanceUC: balanceUC, clock: clck}
}

func (h *Handler) getUserID(r *http.Request) (int, bool) {
	return middleware.GetUserIDFromContext(r.Context())
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

	accounts, err := h.balanceUC.GetBalanceForUser(r.Context(), userID)
	if err != nil {
		httputils.InternalError(w, r, "Failed to get accounts for user")
		return
	}

	if accounts == nil {
		accounts = []models.Account{}
	}

	accountsDTO := AccountsToBalanceAPI(userID, accounts)
	httputils.Success(w, r, accountsDTO)
}

func (h *Handler) parseIDFromURL(r *http.Request, paramName string) (int, error) {
	vars := mux.Vars(r)
	idStr := vars[paramName]
	return strconv.Atoi(idStr)
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

	account, err := h.balanceUC.GetAccountByID(r.Context(), userID, accountID)
	if err != nil {
		httputils.InternalError(w, r, "Failed to get account")
		return
	}

	// Преобразуем в формат для фронтенда (camel_case)
	accountDTO := AccountToAPI(account)
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

	account, err := h.balanceUC.CreateAccount(r.Context(), req, userID)
	if err != nil {
		httputils.InternalError(w, r, "Failed to create account")
		return
	}
	accountDTO := AccountToAPI(account)
	httputils.Success(w, r, accountDTO)
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

	acc, err := h.balanceUC.UpdateAccount(r.Context(), req, userID, accID)
	if err != nil {
		if errors.Is(err, serviceerrors.ErrForbidden) {
			httputils.Error(w, r, "Доступ к счету запрещен", http.StatusForbidden)
			return
		}
		httputils.InternalError(w, r, "failed to update budget")
		return
	}
	accDTO := AccountToAPI(acc)
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

	acc, err := h.balanceUC.DeleteAccount(r.Context(), userID, accID)
	if err != nil {
		if errors.Is(err, serviceerrors.ErrForbidden) {
			httputils.Error(w, r, "Доступ к бюджету запрещен", 403)
			return
		}
		// if errors.Is(err, budgeterrors.ErrBudgetNotFound) {
		// 	httputils.NotFoundError(w, r, "Бюджет не найден")
		// 	return
		// }
		httputils.InternalError(w, r, "failed to delete budget")
		return
	}
	accDTO := AccountToAPI(acc)
	httputils.Success(w, r, accDTO)
}

