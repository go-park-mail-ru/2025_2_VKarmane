package budget

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	httputils "github.com/go-park-mail-ru/2025_2_VKarmane/pkg/http"
	"github.com/gorilla/mux"
)

type Handler struct {
	budgetUC BudgetUseCase
	clock    clock.Clock
}

func NewHandler(budgetUC BudgetUseCase, clck clock.Clock) *Handler {
	return &Handler{budgetUC: budgetUC, clock: clck}
}

func (h *Handler) getUserID(r *http.Request) (int, bool) {
	return middleware.GetUserIDFromContext(r.Context())
}

func (h *Handler) parseIDFromURL(r *http.Request, paramName string) (int, error) {
	vars := mux.Vars(r)
	idStr := vars[paramName]
	return strconv.Atoi(idStr)
}

// GetListBudgets godoc
// @Summary Получение списка бюджетов пользователя
// @Description Возвращает список всех бюджетов пользователя с расчетом фактических расходов
// @Tags budget
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "Список бюджетов пользователя"
// @Failure 401 {object} models.ErrorResponse "Требуется аутентификация (UNAUTHORIZED, TOKEN_MISSING, TOKEN_INVALID, TOKEN_EXPIRED)"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера (INTERNAL_ERROR, DATABASE_ERROR)"
// @Router /budgets [get]
func (h *Handler) GetListBudgets(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r)
	if !ok {
		httputils.Error(w, r, "User not authenticated", http.StatusUnauthorized)
		return
	}

	budgets, err := h.budgetUC.GetBudgetsForUser(r.Context(), userID)
	if err != nil {
		httputils.InternalError(w, r, "Failed to get budgets for user")
		return
	}

	budgetsDTO := BudgetsToAPI(userID, budgets)
	httputils.Success(w, r, budgetsDTO)
}

// GetBudgetByID godoc
// @Summary Получение конкретного бюджета
// @Description Возвращает информацию о конкретном бюджете пользователя
// @Tags budget
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID бюджета"
// @Success 200 {object} map[string]interface{} "Информация о бюджете"
// @Failure 400 {object} models.ErrorResponse "Некорректный ID бюджета (INVALID_REQUEST)"
// @Failure 401 {object} models.ErrorResponse "Требуется аутентификация (UNAUTHORIZED, TOKEN_MISSING, TOKEN_INVALID, TOKEN_EXPIRED)"
// @Failure 404 {object} models.ErrorResponse "Бюджет не найден (BUDGET_NOT_FOUND)"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера (INTERNAL_ERROR, DATABASE_ERROR)"
// @Router /budget/{id} [get]
func (h *Handler) GetBudgetByID(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromURL(r, "id")
	if err != nil {
		httputils.ValidationError(w, r, "Invalid budget ID format", "id")
		return
	}

	userID, ok := h.getUserID(r)
	if !ok {
		httputils.Error(w, r, "User not authenticated", http.StatusUnauthorized)
		return
	}

	budget, err := h.budgetUC.GetBudgetByID(r.Context(), userID, id)
	if err != nil {
		httputils.NotFoundError(w, r, "Budget not found")
		return
	}

	budgetDTO := BudgetToAPI(budget)
	httputils.Success(w, r, budgetDTO)
}
