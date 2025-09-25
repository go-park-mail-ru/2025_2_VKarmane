package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *Handler) GetListBalance(w http.ResponseWriter, r *http.Request) {
	userID := 1

	balance, err := h.svc.GetBalanceForUser(userID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(balance); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func (h *Handler) GetBalanceByAccountID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}

	userID := 1

	balance, err := h.svc.GetBalanceForUser(userID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}

	for _, account := range balance.Accounts {
		if account.ID == id {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(account); err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
			return
		}
	}
	http.Error(w, "not found", http.StatusNotFound)
}
