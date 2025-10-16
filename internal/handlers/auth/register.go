package auth

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
)

func Register(publicRouter *mux.Router, protectedRouter *mux.Router, uc AuthUseCase, log logger.Logger) {
	h := NewHandler(uc, log)

	publicRouter.HandleFunc("/auth/register", h.Register).Methods(http.MethodPost)
	publicRouter.HandleFunc("/auth/login", h.Login).Methods(http.MethodPost)

	protectedRouter.HandleFunc("/auth/logout", h.Logout).Methods(http.MethodPost)
	protectedRouter.HandleFunc("/profile", h.GetProfile).Methods(http.MethodGet)
}
