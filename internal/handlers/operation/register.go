package operation

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Register(r *mux.Router, uc OperationUseCase) {
	handler := NewHandler(uc)

	r.HandleFunc("/operations/account/{acc_id}", handler.GetAccountOperations).Methods(http.MethodGet)
	r.HandleFunc("/operations/account/{acc_id}", handler.CreateOperation).Methods(http.MethodPost)
}
