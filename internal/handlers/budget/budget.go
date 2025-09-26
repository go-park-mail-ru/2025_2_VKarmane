package budget

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/budget"
	"github.com/gorilla/mux"
)

type Handler struct {
	budgetUC *budget.UseCase
}

func NewHandler(budgetUC *budget.UseCase) *Handler {
	return &Handler{budgetUC: budgetUC}
}

func (h *Handler) getUserID(r *http.Request) (int, bool) {
	return middleware.GetUserIDFromContext(r.Context())
}

func (h *Handler) parseIDFromURL(r *http.Request, paramName string) (int, error) {
	vars := mux.Vars(r)
	idStr := vars[paramName]
	return strconv.Atoi(idStr)
}

func (h *Handler) sendJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	log.Printf("Error: %s", message)
	http.Error(w, message, statusCode)
}

func (h *Handler) GetListBudgets(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r)
	if !ok {
		h.sendErrorResponse(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	budgets, err := h.budgetUC.GetBudgetsForUser(userID)
	if err != nil {
		log.Printf("Error getting budgets for user %d: %v", userID, err)
		h.sendErrorResponse(w, "Failed to get budgets for user", http.StatusInternalServerError)
		return
	}

	budgetsDTO := BudgetsToAPI(userID, budgets)
	h.sendJSONResponse(w, budgetsDTO)
}

func (h *Handler) GetBudgetByID(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromURL(r, "id")
	if err != nil {
		log.Printf("Invalid budget ID format: %v", err)
		h.sendErrorResponse(w, "Invalid budget ID format", http.StatusBadRequest)
		return
	}

	userID, ok := h.getUserID(r)
	if !ok {
		h.sendErrorResponse(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	budget, err := h.budgetUC.GetBudgetByID(userID, id)
	if err != nil {
		log.Printf("Error getting budget %d for user %d: %v", id, userID, err)
		h.sendErrorResponse(w, "Budget not found", http.StatusNotFound)
		return
	}

	budgetDTO := BudgetToAPI(budget)
	h.sendJSONResponse(w, budgetDTO)
}
