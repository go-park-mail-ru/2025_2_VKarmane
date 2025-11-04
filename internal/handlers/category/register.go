package category

import (
	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/category"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/image"
)

func Register(router *mux.Router, categoryUC category.CategoryUseCase, imageUC image.ImageUseCase) {
	handler := NewHandler(categoryUC, imageUC)

	router.HandleFunc("/categories", handler.GetCategories).Methods("GET")
	router.HandleFunc("/categories", handler.CreateCategory).Methods("POST")
	router.HandleFunc("/categories/{id}", handler.GetCategoryByID).Methods("GET")
	router.HandleFunc("/categories/{id}", handler.UpdateCategory).Methods("PUT")
	router.HandleFunc("/categories/{id}", handler.DeleteCategory).Methods("DELETE")
}
