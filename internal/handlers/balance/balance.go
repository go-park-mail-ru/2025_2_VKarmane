package balance

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/balance"
	"github.com/gorilla/mux"
)

type Handler struct {
	balanceUC *balance.UseCase
}

func NewHandler(balanceUC *balance.UseCase) *Handler {
	return &Handler{balanceUC: balanceUC}
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
	userID := h.getUserID(r)

	accounts, err := h.balanceUC.GetBalanceForUser(userID)
	if err != nil {
		log.Printf("Error getting balance for user %d: %v", userID, err)
		h.sendErrorResponse(w, "Failed to get balance for user", http.StatusInternalServerError)
		return
	}

	balanceDTO := AccountsToBalanceAPI(userID, accounts)
	h.sendJSONResponse(w, balanceDTO)
}

func (h *Handler) GetBalanceByAccountID(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromURL(r, "id")
	if err != nil {
		log.Printf("Invalid account ID format: %v", err)
		h.sendErrorResponse(w, "Invalid account ID format", http.StatusBadRequest)
		return
	}

	userID := h.getUserID(r)

	account, err := h.balanceUC.GetAccountByID(userID, id)
	if err != nil {
		log.Printf("Error getting account %d for user %d: %v", id, userID, err)
		h.sendErrorResponse(w, "Account not found", http.StatusNotFound)
		return
	}

	accountDTO := AccountToAPI(account)
	h.sendJSONResponse(w, accountDTO)
}
