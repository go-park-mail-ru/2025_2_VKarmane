package image

import (
	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/image"
)

func Register(router *mux.Router, imageUC image.ImageUseCase) {
	handler := NewHandler(imageUC)

	router.HandleFunc("/images/upload", handler.UploadImage).Methods("POST")
	router.HandleFunc("/images/url", handler.GetImageURL).Methods("GET")
	router.HandleFunc("/images", handler.GetImage).Methods("GET")
	router.HandleFunc("/images", handler.DeleteImage).Methods("DELETE")
}
