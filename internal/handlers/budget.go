package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *Handler) GetListBudgets(w http.ResponseWriter, r *http.Request) {
	userID := 1

	budgets, err := h.svc.GetBudgetsForUser(userID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(budgets); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func (h *Handler) GetBudgetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	userID := 1

	budgets, err := h.svc.GetBudgetsForUser(userID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}

	for _, budget := range budgets {
		if budget.ID == id {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(budget); err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
			return
		}
	}
	http.Error(w, "not found", http.StatusNotFound)
}
