package balance

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
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

	// Возвращаем аккаунты в формате, который ожидает фронтенд
	accountsResponse := map[string]interface{}{
		"accounts": accounts,
	}
	httputils.Success(w, r, accountsResponse)
}
