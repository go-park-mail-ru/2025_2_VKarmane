package account

import (
	"net/http"

	"github.com/gorilla/mux"

	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

func Register(r *mux.Router, finClient finpb.FinanceServiceClient) {
	realClock := clock.RealClock{}
	h := NewHandler(finClient, realClock)

	// Старый формат
	r.HandleFunc("/accounts", h.GetAccounts).Methods(http.MethodGet)
	r.HandleFunc("/accounts", h.CreateAccount).Methods(http.MethodPost)
	r.HandleFunc("/account/{id}", h.GetAccountByID).Methods(http.MethodGet)
	r.HandleFunc("/account/{id}", h.UpdateAccount).Methods(http.MethodPut)
	r.HandleFunc("/account/{id}", h.DeleteAccount).Methods(http.MethodDelete)
}
