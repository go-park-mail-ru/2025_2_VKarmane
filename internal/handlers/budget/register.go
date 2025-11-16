package budget

import (
	"net/http"

	"github.com/gorilla/mux"

	bdgpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

func Register(r *mux.Router, uc BudgetUseCase, budgetClient bdgpb.BudgetServiceClient) {
	realClock := clock.RealClock{}
	handler := NewHandler(uc, realClock, budgetClient)

	r.HandleFunc("/budgets", handler.GetListBudgets).Methods(http.MethodGet)
	r.HandleFunc("/budgets", handler.CreateBudget).Methods(http.MethodPost)
	r.HandleFunc("/budgets/{id}", handler.GetBudgetByID).Methods(http.MethodGet)
	r.HandleFunc("/budgets/{id}", handler.UpdateBudget).Methods(http.MethodPut)
	r.HandleFunc("/budgets/{id}", handler.DeleteBudget).Methods(http.MethodDelete)

}
