package profile

import (
	"github.com/gorilla/mux"

	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/image"
)

func Register(router *mux.Router, imageUC image.ImageUseCase, authClient authpb.AuthServiceClient) {
	handler := NewHandler(imageUC, authClient)

	router.HandleFunc("/profile", handler.GetProfile).Methods("GET")
	router.HandleFunc("/profile/edit", handler.UpdateProfile).Methods("PUT")
}
