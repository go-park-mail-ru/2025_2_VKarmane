package auth

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

func Register(publicRouter *mux.Router, protectedRouter *mux.Router, uc AuthUseCase, log logger.Logger) {
	realClock := clock.RealClock{}
	h := NewHandler(uc, realClock, log)

	publicRouter.HandleFunc("/auth/register", h.Register).Methods(http.MethodPost)
	publicRouter.HandleFunc("/auth/login", h.Login).Methods(http.MethodPost)

	protectedRouter.HandleFunc("/auth/logout", h.Logout).Methods(http.MethodPost)
	protectedRouter.HandleFunc("/profile/edit", h.EditUser).Methods(http.MethodPut)
	protectedRouter.HandleFunc("/profile", h.GetProfile).Methods(http.MethodGet)
}
