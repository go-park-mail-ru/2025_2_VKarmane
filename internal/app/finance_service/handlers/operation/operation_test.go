package operation

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGetAccountOperations_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFin := mocks.NewMockFinanceServiceClient(ctrl)
	mockImgUC := mocks.NewMockImageUseCase(ctrl)

	handler := NewHandler(mockFin, mockImgUC, nil, clock.RealClock{})

	// Мокируем gRPC ответ
	opsResp := &finpb.ListOperationsResponse{
		Operations: []*finpb.OperationInList{
			{
				Id:                   1,
				AccountId:            1,
				Name:                 "Test",
				Sum:                  100,
				CategoryLogoHashedId: "img-123",
			},
		},
	}
	mockFin.EXPECT().
		GetOperationsByAccount(gomock.Any(), gomock.Any()).
		Return(opsResp, nil)

	mockImgUC.EXPECT().
		GetImageURL(gomock.Any(), "img-123").
		Return("https://example.com/img-123", nil)

	req := httptest.NewRequest(http.MethodGet, "/operations/account/1", nil)
	req = mux.SetURLVars(req, map[string]string{"acc_id": "1"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.GetAccountOperations(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestCreateOperation_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFin := mocks.NewMockFinanceServiceClient(ctrl)
	mockImgUC := mocks.NewMockImageUseCase(ctrl)
	mockKafka := mocks.NewMockKafkaProducer(ctrl)

	handler := NewHandler(mockFin, mockImgUC, mockKafka, clock.RealClock{})


	opResp := &finpb.Operation{
		Id:        1,
		AccountId: 1,
		Name:      "Test",
		Sum:       100,
		Type:      string(models.OperationExpense),
		Status:    string(models.OperationFinished),
		CreatedAt: timestamppb.New(time.Now()),
		Date:      timestamppb.New(time.Now()),
	}

	mockFin.EXPECT().
		CreateOperation(gomock.Any(), gomock.Any()).
		Return(opResp, nil)


	categoryResp := &finpb.CategoryWithStats{
		Category: &finpb.Category{
			Id:          1,
			LogoHashedId: "img-123",
		},
	}

	mockFin.EXPECT().
		GetCategory(gomock.Any(), gomock.Any()).
		Return(categoryResp, nil)


	mockImgUC.EXPECT().
		GetImageURL(gomock.Any(), "img-123").
		Return("https://example.com/img-123", nil)


	mockKafka.EXPECT().
		WriteMessages(gomock.Any(), gomock.Any()).
		Return(nil)

	reqBody := models.CreateOperationRequest{
		AccountID: 1,
		Type:      models.OperationExpense,
		Name:      "Test",
		Sum:       100,
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/operations/account/1", bytes.NewBuffer(bodyBytes))
	req = mux.SetURLVars(req, map[string]string{"acc_id": "1"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.CreateOperation(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
}


func TestGetOperationByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFin := mocks.NewMockFinanceServiceClient(ctrl)
	mockImgUC := mocks.NewMockImageUseCase(ctrl)

	handler := NewHandler(mockFin, mockImgUC,nil, clock.RealClock{})

	opResp := &finpb.Operation{
		Id:        1,
		AccountId: 1,
		Name:      "Test",
		Sum:       100,
		CreatedAt: timestamppb.New(time.Now()),
		Date:      timestamppb.New(time.Now()),
	}

	mockFin.EXPECT().
		GetOperation(gomock.Any(), gomock.Any()).
		Return(opResp, nil)

	req := httptest.NewRequest(http.MethodGet, "/operations/account/1/operation/1", nil)
	req = mux.SetURLVars(req, map[string]string{"acc_id": "1", "op_id": "1"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.GetOperationByID(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateOperation_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFin := mocks.NewMockFinanceServiceClient(ctrl)
	mockImgUC := mocks.NewMockImageUseCase(ctrl)
	mockKafka := mocks.NewMockKafkaProducer(ctrl)

	handler := NewHandler(mockFin, mockImgUC, mockKafka, clock.RealClock{})

	updateReq := models.UpdateOperationRequest{
		Name: ptrString("Updated"),
		Sum:  ptrFloat64(200),
	}

	opResp := &finpb.Operation{
		Id:        1,
		AccountId: 1,
		Name:      "Updated",
		Sum:       200,
		CreatedAt: timestamppb.New(time.Now()),
		Date:      timestamppb.New(time.Now()),
	}

	mockFin.EXPECT().
		UpdateOperation(gomock.Any(), gomock.Any()).
		Return(opResp, nil)

	categoryResp := &finpb.CategoryWithStats{
		Category: &finpb.Category{
			Id:           1,
			LogoHashedId: "img-123",
		},
	}
	mockFin.EXPECT().
		GetCategory(gomock.Any(), gomock.Any()).
		Return(categoryResp, nil)

	mockImgUC.EXPECT().
		GetImageURL(gomock.Any(), "img-123").
		Return("https://example.com/img-123", nil)

	mockKafka.EXPECT().
		WriteMessages(gomock.Any(), gomock.Any()).
		Return(nil)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/operations/account/1/operation/1", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"acc_id": "1", "op_id": "1"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.UpdateOperation(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}


func TestDeleteOperation_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFin := mocks.NewMockFinanceServiceClient(ctrl)
	mockImgUC := mocks.NewMockImageUseCase(ctrl)
	mockKafka := mocks.NewMockKafkaProducer(ctrl)

	handler := NewHandler(mockFin, mockImgUC, mockKafka, clock.RealClock{})

	opResp := &finpb.Operation{
		Id:        1,
		AccountId: 1,
		Status:    string(models.OperationReverted),
		CreatedAt: timestamppb.New(time.Now()),
		Date:      timestamppb.New(time.Now()),
	}

	mockFin.EXPECT().
		DeleteOperation(gomock.Any(), gomock.Any()).
		Return(opResp, nil)

	mockKafka.EXPECT().
		WriteMessages(gomock.Any(), gomock.Any()).
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/operations/account/1/operation/1", nil)
	req = mux.SetURLVars(req, map[string]string{"acc_id": "1", "op_id": "1"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.DeleteOperation(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}


// утилиты
func ptrString(s string) *string    { return &s }
func ptrFloat64(f float64) *float64 { return &f }
