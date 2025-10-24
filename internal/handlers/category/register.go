package category

import (
	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/category"
)

func Register(router *mux.Router, categoryUC category.CategoryUseCase) {
	handler := NewHandler(categoryUC)

	router.HandleFunc("/categories", handler.GetCategories).Methods("GET")
}
