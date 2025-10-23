package operation

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

func Register(r *mux.Router, uc OperationUseCase) {
	realClock := clock.RealClock{}
	handler := NewHandler(uc, realClock)

	r.HandleFunc("/account/{acc_id}/operations", handler.GetAccountOperations).Methods(http.MethodGet)
	r.HandleFunc("/account/{acc_id}/operations", handler.CreateOperation).Methods(http.MethodPost)
	r.HandleFunc("/account/{acc_id}/operations/{op_id}", handler.GetOperationByID).Methods(http.MethodGet)
	r.HandleFunc("/account/{acc_id}/operations/{op_id}", handler.UpdateOperation).Methods(http.MethodPut)
	r.HandleFunc("/account/{acc_id}/operations/{op_id}", handler.DeleteOperation).Methods(http.MethodDelete)
}
