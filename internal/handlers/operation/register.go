package operation

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Register(r *mux.Router, uc OperationUseCase) {
	handler := NewHandler(uc)

	r.HandleFunc("/account/{acc_id}/operations", handler.GetAccountOperations).Methods(http.MethodGet)
	r.HandleFunc("/account/{acc_id}/operations", handler.CreateOperation).Methods(http.MethodPost)
	r.HandleFunc("/account/{acc_id}/operations/{op_id}", handler.GetOperationByID).Methods(http.MethodGet)
	r.HandleFunc("/account/{acc_id}/operations/update/{op_id}", handler.UpdateOperation).Methods(http.MethodPost)
	r.HandleFunc("/account/{acc_id}/operations/delete/{op_id}", handler.DeleteOperation).Methods(http.MethodPost)
}
