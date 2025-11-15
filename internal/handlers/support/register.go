package support

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/support"
	"github.com/gorilla/mux"
)

func Register(router *mux.Router, supportUC *support.UseCase) {
	handler := NewHandler(supportUC)

	router.HandleFunc("/support/{user_id}", handler.CreateSupportRequest).Methods("POST")
	router.HandleFunc("/support/{user_id}", handler.GetUserSupportRequests).Methods("GET")
	router.HandleFunc("/support/{req_id}/status", handler.UpdateSupportStatus).Methods("PUT")
	router.HandleFunc("/support/stats", handler.GetSupportStats).Methods("GET")
}
