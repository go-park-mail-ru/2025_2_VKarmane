package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers/balance"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers/budget"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase"
	"github.com/gorilla/mux"
)

type Handler struct {
	balanceHandler *balance.Handler
	budgetHandler  *budget.Handler
}

func NewHandler(uc *usecase.UseCase) *Handler {
	return &Handler{
		balanceHandler: balance.NewHandler(uc.BalanceUC),
		budgetHandler:  budget.NewHandler(uc.BudgetUC),
	}
}

func (h *Handler) getUserID(_ *http.Request) int {
	// TODO: Extract from JWT token or session
	return 1
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

func (h *Handler) GetListBalance(w http.ResponseWriter, r *http.Request) {
	h.balanceHandler.GetListBalance(w, r)
}

func (h *Handler) GetBalanceByAccountID(w http.ResponseWriter, r *http.Request) {
	h.balanceHandler.GetBalanceByAccountID(w, r)
}

func (h *Handler) GetListBudgets(w http.ResponseWriter, r *http.Request) {
	h.budgetHandler.GetListBudgets(w, r)
}

func (h *Handler) GetBudgetByID(w http.ResponseWriter, r *http.Request) {
	h.budgetHandler.GetBudgetByID(w, r)
}
