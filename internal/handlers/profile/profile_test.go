package profile

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
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetProfile_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileUC := mocks.NewMockProfileUseCase(ctrl)
	mockImageUC := mocks.NewMockImageUseCase(ctrl)
	handler := NewHandler(mockProfileUC, mockImageUC)

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	rr := httptest.NewRecorder()

	handler.GetProfile(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileUC := mocks.NewMockProfileUseCase(ctrl)
	mockImageUC := mocks.NewMockImageUseCase(ctrl)
	handler := NewHandler(mockProfileUC, mockImageUC)

	profile := models.ProfileResponse{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}

	mockProfileUC.EXPECT().GetProfile(gomock.Any(), 1).Return(profile, nil)

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.GetProfile(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileUC := mocks.NewMockProfileUseCase(ctrl)
	mockImageUC := mocks.NewMockImageUseCase(ctrl)
	handler := NewHandler(mockProfileUC, mockImageUC)

	profile := models.ProfileResponse{
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane@example.com",
	}

	mockProfileUC.EXPECT().UpdateProfile(gomock.Any(), gomock.Any(), 1).Return(profile, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("first_name", "Jane")
	_ = writer.WriteField("last_name", "Smith")
	_ = writer.WriteField("email", "jane@example.com")
	writer.Close()

	req := httptest.NewRequest(http.MethodPut, "/profile/edit", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.UpdateProfile(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateProfile_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileUC := mocks.NewMockProfileUseCase(ctrl)
	mockImageUC := mocks.NewMockImageUseCase(ctrl)
	handler := NewHandler(mockProfileUC, mockImageUC)

	req := httptest.NewRequest(http.MethodPut, "/profile/edit", nil)
	rr := httptest.NewRecorder()

	handler.UpdateProfile(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestUpdateProfile_WithImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileUC := mocks.NewMockProfileUseCase(ctrl)
	mockImageUC := mocks.NewMockImageUseCase(ctrl)
	handler := NewHandler(mockProfileUC, mockImageUC)

	profile := models.ProfileResponse{
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane@example.com",
	}

	mockImageUC.EXPECT().UploadImage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("img-123", nil)
	mockProfileUC.EXPECT().UpdateProfile(gomock.Any(), gomock.Any(), 1).Return(profile, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("first_name", "Jane")
	_ = writer.WriteField("last_name", "Smith")
	_ = writer.WriteField("email", "jane@example.com")
	part, _ := writer.CreateFormFile("avatar", "avatar.jpg")
	part.Write([]byte("fake image data"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPut, "/profile/edit", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	handler.UpdateProfile(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

