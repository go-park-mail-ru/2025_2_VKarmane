package balance

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Register(r *mux.Router, uc BalanceUseCase) {
	h := NewHandler(uc)

	r.HandleFunc("/accounts", h.GetAccounts).Methods(http.MethodGet)
}
