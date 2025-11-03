package operation

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

func Register(r *mux.Router, uc OperationUseCase) {
	realClock := clock.RealClock{}
	handler := NewHandler(uc, realClock)

	r.HandleFunc("/operations/account/{acc_id}", handler.GetAccountOperations).Methods(http.MethodGet)
	r.HandleFunc("/operations/account/{acc_id}", handler.CreateOperation).Methods(http.MethodPost)
	r.HandleFunc("/operations/account/{acc_id}/operation/{op_id}", handler.GetOperationByID).Methods(http.MethodGet)
	r.HandleFunc("/operations/account/{acc_id}/operation/{op_id}", handler.UpdateOperation).Methods(http.MethodPut)
	r.HandleFunc("/operations/account/{acc_id}/operation/{op_id}", handler.DeleteOperation).Methods(http.MethodDelete)
}
