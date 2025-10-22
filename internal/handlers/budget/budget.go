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
