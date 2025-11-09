package operation

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/image"
)

func Register(r *mux.Router, uc OperationUseCase, imageUC image.ImageUseCase) {
	realClock := clock.RealClock{}
	handler := NewHandler(uc, imageUC, realClock)

	// Старый формат
	r.HandleFunc("/operations/account/{acc_id}", handler.GetAccountOperations).Methods(http.MethodGet)
	r.HandleFunc("/operations/account/{acc_id}", handler.CreateOperation).Methods(http.MethodPost)
	r.HandleFunc("/operations/account/{acc_id}/operation/{op_id}", handler.GetOperationByID).Methods(http.MethodGet)
	r.HandleFunc("/operations/account/{acc_id}/operation/{op_id}", handler.UpdateOperation).Methods(http.MethodPut)
	r.HandleFunc("/operations/account/{acc_id}/operation/{op_id}", handler.DeleteOperation).Methods(http.MethodDelete)

	// Новый формат (для фронтенда)
	r.HandleFunc("/account/{acc_id}/operations", handler.GetAccountOperations).Methods(http.MethodGet)
	r.HandleFunc("/account/{acc_id}/operations", handler.CreateOperation).Methods(http.MethodPost)
	r.HandleFunc("/account/{acc_id}/operations/{op_id}", handler.GetOperationByID).Methods(http.MethodGet)
	r.HandleFunc("/account/{acc_id}/operations/{op_id}", handler.UpdateOperation).Methods(http.MethodPut)
	r.HandleFunc("/account/{acc_id}/operations/{op_id}", handler.DeleteOperation).Methods(http.MethodDelete)
}
