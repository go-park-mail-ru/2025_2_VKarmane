package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRegisterSuccess(t *testing.T) {
	m := mocks.NewAuthUseCase(t)
	h := NewHandler(m, logger.NewSlogLogger())

	reqBody := models.RegisterRequest{Email: "u@e.co", Login: "user1", Password: "password"}
	b, _ := json.Marshal(reqBody)

	m.On("Register", mock.Anything, reqBody).Return(models.AuthResponse{Token: "tok", User: models.User{ID: 1, Login: "user1"}}, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(b))
	rr := httptest.NewRecorder()
	h.Register(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
}

func TestRegisterValidationError(t *testing.T) {
	m := mocks.NewAuthUseCase(t)
	h := NewHandler(m, logger.NewSlogLogger())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString("{"))
	rr := httptest.NewRecorder()
	h.Register(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestLoginUnauthorizedPaths(t *testing.T) {
	m := mocks.NewAuthUseCase(t)
	h := NewHandler(m, logger.NewSlogLogger())

	body := models.LoginRequest{Login: "user", Password: "pass123"}
	b, _ := json.Marshal(body)

	m.On("Login", mock.Anything, body).Return(models.AuthResponse{}, errors.New("user not found"))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(b))
	rr := httptest.NewRecorder()
	h.Login(rr, req)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetProfileUnauthorized(t *testing.T) {
	m := mocks.NewAuthUseCase(t)
	h := NewHandler(m, logger.NewSlogLogger())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
	rr := httptest.NewRecorder()
	h.GetProfile(rr, req)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestEditUserSuccess(t *testing.T) {
	m := mocks.NewAuthUseCase(t)
	h := NewHandler(m, logger.NewSlogLogger())

	reqBody := models.UpdateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}
	b, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/profile/edit", bytes.NewReader(b))
	req = req.WithContext(context.WithValue(context.Background(), middleware.UserIDKey, 1))
	rr := httptest.NewRecorder()

	expectedUser := models.User{ID: 1, FirstName: "John", LastName: "Doe", Email: "john@example.com"}
	m.On("EditUserByID", mock.Anything, reqBody, 1).Return(expectedUser, nil)

	h.EditUser(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
}

func TestEditUserUnauthorized(t *testing.T) {
	m := mocks.NewAuthUseCase(t)
	h := NewHandler(m, logger.NewSlogLogger())

	reqBody := models.UpdateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}
	b, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/profile/edit", bytes.NewReader(b))
	rr := httptest.NewRecorder()

	h.EditUser(rr, req)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
}
