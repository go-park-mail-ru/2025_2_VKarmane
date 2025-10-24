package balance

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

func Register(r *mux.Router, uc BalanceUseCase) {
	realClock := clock.RealClock{}
	h := NewHandler(uc, realClock)

	r.HandleFunc("/accounts", h.GetAccounts).Methods(http.MethodGet)
}
