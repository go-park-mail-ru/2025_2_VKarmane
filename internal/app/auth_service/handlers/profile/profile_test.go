package profile

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetProfile_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mocks.NewMockAuthServiceClient(ctrl)
	mockImage := mocks.NewMockImageUseCase(ctrl)
	h := NewHandler(mockImage, mockAuth)

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	rr := httptest.NewRecorder()

	h.GetProfile(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mocks.NewMockAuthServiceClient(ctrl)
	mockImage := mocks.NewMockImageUseCase(ctrl)
	h := NewHandler(mockImage, mockAuth)

	profile := &authpb.ProfileResponse{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}

	mockAuth.EXPECT().
		GetProfile(gomock.Any(), gomock.Any()).
		Return(profile, nil)

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	h.GetProfile(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateProfile_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mocks.NewMockAuthServiceClient(ctrl)
	mockImage := mocks.NewMockImageUseCase(ctrl)
	h := NewHandler(mockImage, mockAuth)

	req := httptest.NewRequest(http.MethodPut, "/profile/edit", nil)
	rr := httptest.NewRecorder()

	h.UpdateProfile(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestUpdateProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mocks.NewMockAuthServiceClient(ctrl)
	mockImage := mocks.NewMockImageUseCase(ctrl)
	h := NewHandler(mockImage, mockAuth)

	profile := &authpb.ProfileResponse{
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane@example.com",
	}

	mockAuth.EXPECT().
		UpdateProfile(gomock.Any(), gomock.Any()).
		Return(profile, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("first_name", "Jane")
	writer.WriteField("last_name", "Smith")
	writer.WriteField("email", "jane@example.com")
	writer.Close()

	req := httptest.NewRequest(http.MethodPut, "/profile/edit", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))

	rr := httptest.NewRecorder()

	h.UpdateProfile(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateProfile_WithImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mocks.NewMockAuthServiceClient(ctrl)
	mockImage := mocks.NewMockImageUseCase(ctrl)
	h := NewHandler(mockImage, mockAuth)

	profile := &authpb.ProfileResponse{
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane@example.com",
	}

	// Мок загрузки аватарки
	mockImage.EXPECT().
		UploadImage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return("img-123", nil)

	// Мок gRPC UpdateProfile
	mockAuth.EXPECT().
		UpdateProfile(gomock.Any(), gomock.Any()).
		Return(profile, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("first_name", "Jane")
	writer.WriteField("last_name", "Smith")
	writer.WriteField("email", "jane@example.com")

	part, _ := writer.CreateFormFile("avatar", "avatar.jpg")
	part.Write([]byte("fake image data"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPut, "/profile/edit", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	h.UpdateProfile(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}
