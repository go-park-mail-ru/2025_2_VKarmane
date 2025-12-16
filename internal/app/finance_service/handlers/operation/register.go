package operation

import (
	"net/http"

	"github.com/gorilla/mux"

	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
	image "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/image/usecase"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	kafkautils "github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/kafka"
)

func Register(r *mux.Router, finClient finpb.FinanceServiceClient, imageUC image.ImageUseCase, kafkaProducer kafkautils.KafkaProducer) {
	realClock := clock.RealClock{}
	handler := NewHandler(finClient, imageUC, kafkaProducer, realClock)

	// Старый формат
	r.HandleFunc("/operations/account/{acc_id}", handler.GetAccountOperations).Methods(http.MethodGet)
	r.HandleFunc("/operations/account/{acc_id}", handler.CreateOperation).Methods(http.MethodPost)
	r.HandleFunc("/operations/account/{acc_id}/operation/{op_id}", handler.GetOperationByID).Methods(http.MethodGet)
	r.HandleFunc("/operations/account/{acc_id}/operation/{op_id}", handler.UpdateOperation).Methods(http.MethodPut)
	r.HandleFunc("/operations/account/{acc_id}/operation/{op_id}", handler.DeleteOperation).Methods(http.MethodDelete)

	// Новый формат (для фронтенда)
	r.HandleFunc("/operations/import", handler.UploadCVSData).Methods(http.MethodPost)
	r.HandleFunc("/operations/export", handler.GetCSVData).Methods(http.MethodGet)
	r.HandleFunc("/account/{acc_id}/operations", handler.GetAccountOperations).Methods(http.MethodGet)
	r.HandleFunc("/account/{acc_id}/operations", handler.CreateOperation).Methods(http.MethodPost)
	r.HandleFunc("/account/{acc_id}/operations/{op_id}", handler.GetOperationByID).Methods(http.MethodGet)
	r.HandleFunc("/account/{acc_id}/operations/{op_id}", handler.UpdateOperation).Methods(http.MethodPut)
	r.HandleFunc("/account/{acc_id}/operations/{op_id}", handler.DeleteOperation).Methods(http.MethodDelete)
}
