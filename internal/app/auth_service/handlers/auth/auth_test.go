package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRegister_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockAuthServiceClient(ctrl)
	handler := NewHandler(clock.RealClock{}, logger.NewSlogLogger(), mockClient)

	registerReq := models.RegisterRequest{
		Email:    "test@example.com",
		Login:    "testuser",
		Password: "password123",
	}

	mockClient.
		EXPECT().
		Register(gomock.Any(), gomock.Any()).
		Return(&authpb.AuthResponse{
			Token: "jwt-token",
			User: &authpb.User{
				Id:    1,
				Login: "testuser",
				Email: "test@example.com",
			},
		}, nil)

	body, _ := json.Marshal(registerReq)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
}

func TestLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockAuthServiceClient(ctrl)
	handler := NewHandler(clock.RealClock{}, logger.NewSlogLogger(), mockClient)

	loginReq := models.LoginRequest{
		Login:    "testuser",
		Password: "password123",
	}

	mockClient.
		EXPECT().
		Login(gomock.Any(), gomock.Any()).
		Return(&authpb.AuthResponse{
			Token: "jwt-token",
			User: &authpb.User{
				Id:    1,
				Login: "testuser",
			},
		}, nil)

	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestLogout_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockAuthServiceClient(ctrl)
	handler := NewHandler(clock.RealClock{}, logger.NewSlogLogger(), mockClient)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	rr := httptest.NewRecorder()

	handler.Logout(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestRegister_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockAuthServiceClient(ctrl)
	handler := NewHandler(clock.RealClock{}, logger.NewSlogLogger(), mockClient)

	invalidReq := map[string]string{"email": "invalid"}
	body, _ := json.Marshal(invalidReq)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestLogin_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockAuthServiceClient(ctrl)
	handler := NewHandler(clock.RealClock{}, logger.NewSlogLogger(), mockClient)

	invalidReq := map[string]string{"login": ""}

	body, _ := json.Marshal(invalidReq)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestRegister_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockAuthServiceClient(ctrl)
	handler := NewHandler(clock.RealClock{}, logger.NewSlogLogger(), mockClient)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString("invalid json"))
	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestLogin_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockAuthServiceClient(ctrl)
	handler := NewHandler(clock.RealClock{}, logger.NewSlogLogger(), mockClient)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString("invalid json"))
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetCSRFToken_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockAuthServiceClient(ctrl)
	handler := NewHandler(clock.RealClock{}, logger.NewSlogLogger(), mockClient)

	mockClient.
		EXPECT().
		GetCSRF(gomock.Any(), gomock.Any()).
		Return(&authpb.CSRFTokenResponse{Token: "csrf123"}, nil)

	req := httptest.NewRequest(http.MethodGet, "/auth/csrf", nil)
	rr := httptest.NewRecorder()

	handler.GetCSRFToken(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}
