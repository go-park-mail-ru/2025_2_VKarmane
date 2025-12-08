package category

import (
	"github.com/gorilla/mux"

	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
	image "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/image/usecase"
	kafkautils "github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/kafka"
)

func Register(router *mux.Router, finClient finpb.FinanceServiceClient, imageUC image.ImageUseCase, kafkaProducer kafkautils.KafkaProducer) {
	handler := NewHandler(finClient, imageUC, kafkaProducer)

	router.HandleFunc("/categories", handler.GetCategories).Methods("GET")
	router.HandleFunc("/categories", handler.CreateCategory).Methods("POST")
	router.HandleFunc("/categories/report", handler.GetCategoriesReport).Methods("GET")
	router.HandleFunc("/categories/{id}", handler.GetCategoryByID).Methods("GET")
	router.HandleFunc("/categories/{id}", handler.UpdateCategory).Methods("PUT")
	router.HandleFunc("/categories/{id}", handler.DeleteCategory).Methods("DELETE")
}
