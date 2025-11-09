package budget

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

func Register(r *mux.Router, uc BudgetUseCase) {
	realClock := clock.RealClock{}
	handler := NewHandler(uc, realClock)

	r.HandleFunc("/budgets", handler.GetListBudgets).Methods(http.MethodGet)
	r.HandleFunc("/budgets/{id}", handler.GetBudgetByID).Methods(http.MethodGet)
}
