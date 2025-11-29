package category

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetCategories_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFin := mocks.NewMockFinanceServiceClient(ctrl)
	handler := NewHandler(mockFin, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	rr := httptest.NewRecorder()

	handler.GetCategories(rr, req)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetCategories_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFin := mocks.NewMockFinanceServiceClient(ctrl)
	handler := NewHandler(mockFin, nil, nil)

	mockFin.EXPECT().
		GetCategoriesWithStatsByUser(gomock.Any(), gomock.Any()).
		Return(&finpb.ListCategoriesWithStatsResponse{
			Categories: []*finpb.CategoryWithStats{
				{
					Category: &finpb.Category{
						Id:   1,
						Name: "Food",
					},
					OperationsCount: 5,
				},
			},
		}, nil)

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.GetCategories(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
}

func TestCreateCategory_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFin := mocks.NewMockFinanceServiceClient(ctrl)
	handler := NewHandler(mockFin, nil, nil)

	mockFin.EXPECT().
		CreateCategory(gomock.Any(), gomock.Any()).
		Return(&finpb.Category{
			Id:     1,
			UserId: 1,
			Name:   "Food",
		}, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("name", "Food")
	_ = writer.WriteField("description", "Food expenses")
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/categories", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.CreateCategory(rr, req)
	require.Equal(t, http.StatusCreated, rr.Code)
}

func TestGetCategoryByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFin := mocks.NewMockFinanceServiceClient(ctrl)
	handler := NewHandler(mockFin, nil, nil)

	mockFin.EXPECT().
		GetCategory(gomock.Any(), gomock.Any()).
		Return(&finpb.CategoryWithStats{
			Category: &finpb.Category{
				Id:   1,
				Name: "Food",
			},
			OperationsCount: 5,
		}, nil)

	req := httptest.NewRequest(http.MethodGet, "/categories/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.GetCategoryByID(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
}

func TestGetCategoryByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFin := mocks.NewMockFinanceServiceClient(ctrl)
	handler := NewHandler(mockFin, nil, nil)

	mockFin.EXPECT().
		GetCategory(gomock.Any(), gomock.Any()).
		Return(nil, status.Error(codes.NotFound, "not found"))

	req := httptest.NewRequest(http.MethodGet, "/categories/5", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "5"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.GetCategoryByID(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdateCategory_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFin := mocks.NewMockFinanceServiceClient(ctrl)
	mockKafka := mocks.NewMockKafkaProducer(ctrl)

	handler := NewHandler(mockFin, nil, mockKafka)

	mockFin.EXPECT().
		UpdateCategory(gomock.Any(), gomock.Any()).
		
		Return(&finpb.Category{
			Id:     1,
			UserId: 1,
			Name:   "Updated",
		}, nil)

	mockKafka.EXPECT().
		WriteMessages(gomock.Any(), gomock.Any()).
		Return(nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("name", "Updated")
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPut, "/categories/1", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.UpdateCategory(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteCategory_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFin := mocks.NewMockFinanceServiceClient(ctrl)
	handler := NewHandler(mockFin, nil, nil)

	// DeleteCategory возвращает *finpb.Category (пустой объект)
	mockFin.EXPECT().
		DeleteCategory(gomock.Any(), gomock.Any()).
		Return(&finpb.Category{}, nil)

	req := httptest.NewRequest(http.MethodDelete, "/categories/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.DeleteCategory(rr, req)
	require.Equal(t, http.StatusNoContent, rr.Code)
}

func TestCreateCategory_WithImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFin := mocks.NewMockFinanceServiceClient(ctrl)
	mockImgUC := mocks.NewMockImageUseCase(ctrl)
	handler := NewHandler(mockFin, mockImgUC, nil)

	// Мок на UploadImage
	mockImgUC.EXPECT().
		UploadImage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return("img-123", nil)

	// Мок на GetImageURL (вызывается в enrichCategoryWithLogoURL)
	mockImgUC.EXPECT().
		GetImageURL(gomock.Any(), "img-123").
		Return("https://example.com/img-123", nil)

	// Мок на CreateCategory
	mockFin.EXPECT().
		CreateCategory(gomock.Any(), gomock.Any()).
		Return(&finpb.Category{
			Id:           1,
			UserId:       1,
			Name:         "Food",
			LogoHashedId: "img-123",
		}, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("name", "Food")
	part, _ := writer.CreateFormFile("image", "category.jpg")
	_, _ = part.Write([]byte("fake image data"))
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/categories", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.CreateCategory(rr, req)
	require.Equal(t, http.StatusCreated, rr.Code)
}
