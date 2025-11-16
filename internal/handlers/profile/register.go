package profile

import (
	"github.com/gorilla/mux"

	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/image"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/profile"
)

func Register(router *mux.Router, profileUC profile.ProfileUseCase, imageUC image.ImageUseCase, authClient authpb.AuthServiceClient) {
	handler := NewHandler(profileUC, imageUC, authClient)

	router.HandleFunc("/profile", handler.GetProfile).Methods("GET")
	router.HandleFunc("/profile/edit", handler.UpdateProfile).Methods("PUT")
}
