package category

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetCategories_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCatUC := mocks.NewMockCategoryUseCase(ctrl)
	mockImgUC := mocks.NewMockImageUseCase(ctrl)
	handler := NewHandler(mockCatUC, mockImgUC)

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	rr := httptest.NewRecorder()

	handler.GetCategories(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetCategories_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCatUC := mocks.NewMockCategoryUseCase(ctrl)
	mockImgUC := mocks.NewMockImageUseCase(ctrl)
	handler := NewHandler(mockCatUC, mockImgUC)

	categories := []models.CategoryWithStats{
		{Category: models.Category{ID: 1, Name: "Food"}, OperationsCount: 5},
	}

	mockCatUC.EXPECT().GetCategoriesByUser(gomock.Any(), 1).Return(categories, nil)

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.GetCategories(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestCreateCategory_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCatUC := mocks.NewMockCategoryUseCase(ctrl)
	mockImgUC := mocks.NewMockImageUseCase(ctrl)
	handler := NewHandler(mockCatUC, mockImgUC)

	category := models.Category{
		ID:   1,
		Name: "Food",
	}

	mockCatUC.EXPECT().CreateCategory(gomock.Any(), gomock.Any(), 1).Return(category, nil)

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

	mockCatUC := mocks.NewMockCategoryUseCase(ctrl)
	mockImgUC := mocks.NewMockImageUseCase(ctrl)
	handler := NewHandler(mockCatUC, mockImgUC)

	category := models.CategoryWithStats{
		Category:        models.Category{ID: 1, Name: "Food"},
		OperationsCount: 5,
	}

	mockCatUC.EXPECT().GetCategoryByID(gomock.Any(), 1, 1).Return(category, nil)

	req := httptest.NewRequest(http.MethodGet, "/categories/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.GetCategoryByID(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteCategory_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCatUC := mocks.NewMockCategoryUseCase(ctrl)
	mockImgUC := mocks.NewMockImageUseCase(ctrl)
	handler := NewHandler(mockCatUC, mockImgUC)

	mockCatUC.EXPECT().DeleteCategory(gomock.Any(), 1, 1).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/categories/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.DeleteCategory(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)
}

func TestGetCategoryByID_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCatUC := mocks.NewMockCategoryUseCase(ctrl)
	mockImgUC := mocks.NewMockImageUseCase(ctrl)
	handler := NewHandler(mockCatUC, mockImgUC)

	req := httptest.NewRequest(http.MethodGet, "/categories/invalid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.GetCategoryByID(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDeleteCategory_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCatUC := mocks.NewMockCategoryUseCase(ctrl)
	mockImgUC := mocks.NewMockImageUseCase(ctrl)
	handler := NewHandler(mockCatUC, mockImgUC)

	req := httptest.NewRequest(http.MethodDelete, "/categories/invalid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.DeleteCategory(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateCategory_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCatUC := mocks.NewMockCategoryUseCase(ctrl)
	mockImgUC := mocks.NewMockImageUseCase(ctrl)
	handler := NewHandler(mockCatUC, mockImgUC)

	req := httptest.NewRequest(http.MethodPost, "/categories", nil)
	rr := httptest.NewRecorder()

	handler.CreateCategory(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetCategoryByID_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCatUC := mocks.NewMockCategoryUseCase(ctrl)
	mockImgUC := mocks.NewMockImageUseCase(ctrl)
	handler := NewHandler(mockCatUC, mockImgUC)

	req := httptest.NewRequest(http.MethodGet, "/categories/1", nil)
	rr := httptest.NewRecorder()

	handler.GetCategoryByID(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestDeleteCategory_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCatUC := mocks.NewMockCategoryUseCase(ctrl)
	mockImgUC := mocks.NewMockImageUseCase(ctrl)
	handler := NewHandler(mockCatUC, mockImgUC)

	req := httptest.NewRequest(http.MethodDelete, "/categories/1", nil)
	rr := httptest.NewRecorder()

	handler.DeleteCategory(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestUpdateCategory_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCatUC := mocks.NewMockCategoryUseCase(ctrl)
	mockImgUC := mocks.NewMockImageUseCase(ctrl)
	handler := NewHandler(mockCatUC, mockImgUC)

	category := models.Category{
		ID:   1,
		Name: "Updated",
	}

	mockCatUC.EXPECT().UpdateCategory(gomock.Any(), gomock.Any(), 1, 1).Return(category, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("name", "Updated")
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPut, "/categories/1", body)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.UpdateCategory(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateCategory_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCatUC := mocks.NewMockCategoryUseCase(ctrl)
	mockImgUC := mocks.NewMockImageUseCase(ctrl)
	handler := NewHandler(mockCatUC, mockImgUC)

	req := httptest.NewRequest(http.MethodPut, "/categories/1", nil)
	rr := httptest.NewRecorder()

	handler.UpdateCategory(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestCreateCategory_WithImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCatUC := mocks.NewMockCategoryUseCase(ctrl)
	mockImgUC := mocks.NewMockImageUseCase(ctrl)
	handler := NewHandler(mockCatUC, mockImgUC)

	category := models.Category{
		ID:   1,
		Name: "Food",
	}

	mockImgUC.EXPECT().UploadImage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("img-123", nil)
	mockCatUC.EXPECT().CreateCategory(gomock.Any(), gomock.Any(), 1).Return(category, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("name", "Food")
	_ = writer.WriteField("description", "Food expenses")
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

func TestUpdateCategory_WithImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCatUC := mocks.NewMockCategoryUseCase(ctrl)
	mockImgUC := mocks.NewMockImageUseCase(ctrl)
	handler := NewHandler(mockCatUC, mockImgUC)

	category := models.Category{
		ID:   1,
		Name: "Updated",
	}

	mockImgUC.EXPECT().UploadImage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("img-123", nil)
	mockCatUC.EXPECT().UpdateCategory(gomock.Any(), gomock.Any(), 1, 1).Return(category, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("name", "Updated")
	part, _ := writer.CreateFormFile("image", "category.jpg")
	_, _ = part.Write([]byte("fake image data"))
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPut, "/categories/1", body)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.UpdateCategory(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}
