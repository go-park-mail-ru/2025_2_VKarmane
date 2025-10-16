package balance

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Register(r *mux.Router, uc BalanceUseCase) {
	h := NewHandler(uc)

	r.HandleFunc("/balance", h.GetListBalance).Methods(http.MethodGet)
	r.HandleFunc("/balance/{id}", h.GetBalanceByAccountID).Methods(http.MethodGet)
}
