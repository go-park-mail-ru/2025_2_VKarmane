package auth

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/proto"
)

func Register(publicRouter *mux.Router, protectedRouter *mux.Router, log logger.Logger, authClient authpb.AuthServiceClient) {
	realClock := clock.RealClock{}
	h := NewHandler(realClock, log, authClient)

	publicRouter.HandleFunc("/auth/csrf", h.GetCSRFToken).Methods(http.MethodGet)
	publicRouter.HandleFunc("/auth/register", h.Register).Methods(http.MethodPost)
	publicRouter.HandleFunc("/auth/login", h.Login).Methods(http.MethodPost)

	protectedRouter.HandleFunc("/auth/logout", h.Logout).Methods(http.MethodPost)
}
