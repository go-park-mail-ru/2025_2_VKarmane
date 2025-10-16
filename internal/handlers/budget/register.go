package budget

import (
	"github.com/gorilla/mux"

	"net/http"
)

func Register(r *mux.Router, uc BudgetUseCase) {
	handler := NewHandler(uc)

	r.HandleFunc("/budgets", handler.GetListBudgets).Methods(http.MethodGet)
	r.HandleFunc("/budget/{id}", handler.GetBudgetByID).Methods(http.MethodGet)
}
