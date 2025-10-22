package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	"github.com/go-park-mail-ru/2025_2_VKarmane/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseCase_Register(t *testing.T) {
	tests := []struct {
		name           string
		request        models.RegisterRequest
		mockResponse   models.AuthResponse
		mockError      error
		expectedResult models.AuthResponse
		expectedError  error
	}{
		{
			name: "successful registration",
			request: models.RegisterRequest{
				Email:    "john@example.com",
				Login:    "johndoe",
				Password: "password123",
			},
			mockResponse: models.AuthResponse{
				Token: "jwt-token-123",
				User: models.User{
					ID:        1,
					FirstName: "",
					LastName:  "",
					Email:     "john@example.com",
					Login:     "johndoe",
					CreatedAt: time.Now(),
				},
			},
			mockError: nil,
			expectedResult: models.AuthResponse{
				Token: "jwt-token-123",
				User: models.User{
					ID:        1,
					FirstName: "",
					LastName:  "",
					Email:     "john@example.com",
					Login:     "johndoe",
					CreatedAt: time.Now(),
				},
			},
			expectedError: nil,
		},
		{
			name: "registration with service error",
			request: models.RegisterRequest{
				Email:    "john@example.com",
				Login:    "johndoe",
				Password: "password123",
			},
			mockResponse:   models.AuthResponse{},
			mockError:      errors.New("user already exists"),
			expectedResult: models.AuthResponse{},
			expectedError:  errors.New("user already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthService := &mocks.AuthService{}
			uc := NewUseCase(mockAuthService, clock.RealClock{})

			mockAuthService.On("Register", mock.Anything, tt.request).Return(tt.mockResponse, tt.mockError)

			result, err := uc.Register(context.Background(), tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult.Token, result.Token)
				assert.Equal(t, tt.expectedResult.User.ID, result.User.ID)
				assert.Equal(t, tt.expectedResult.User.FirstName, result.User.FirstName)
				assert.Equal(t, tt.expectedResult.User.LastName, result.User.LastName)
				assert.Equal(t, tt.expectedResult.User.Email, result.User.Email)
				assert.Equal(t, tt.expectedResult.User.Login, result.User.Login)
			}

			mockAuthService.AssertExpectations(t)
		})
	}
}

func TestUseCase_Login(t *testing.T) {
	tests := []struct {
		name           string
		request        models.LoginRequest
		mockResponse   models.AuthResponse
		mockError      error
		expectedResult models.AuthResponse
		expectedError  error
	}{
		{
			name: "successful login",
			request: models.LoginRequest{
				Login:    "johndoe",
				Password: "password123",
			},
			mockResponse: models.AuthResponse{
				Token: "jwt-token-123",
				User: models.User{
					ID:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@example.com",
					Login:     "johndoe",
					CreatedAt: time.Now(),
				},
			},
			mockError: nil,
			expectedResult: models.AuthResponse{
				Token: "jwt-token-123",
				User: models.User{
					ID:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@example.com",
					Login:     "johndoe",
					CreatedAt: time.Now(),
				},
			},
			expectedError: nil,
		},
		{
			name: "login with invalid credentials",
			request: models.LoginRequest{
				Login:    "johndoe",
				Password: "wrongpassword",
			},
			mockResponse:   models.AuthResponse{},
			mockError:      errors.New("invalid credentials"),
			expectedResult: models.AuthResponse{},
			expectedError:  errors.New("invalid credentials"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthService := &mocks.AuthService{}
			uc := NewUseCase(mockAuthService, clock.RealClock{})

			mockAuthService.On("Login", mock.Anything, tt.request).Return(tt.mockResponse, tt.mockError)

			result, err := uc.Login(context.Background(), tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult.Token, result.Token)
				assert.Equal(t, tt.expectedResult.User.ID, result.User.ID)
				assert.Equal(t, tt.expectedResult.User.FirstName, result.User.FirstName)
				assert.Equal(t, tt.expectedResult.User.LastName, result.User.LastName)
				assert.Equal(t, tt.expectedResult.User.Email, result.User.Email)
				assert.Equal(t, tt.expectedResult.User.Login, result.User.Login)
			}

			mockAuthService.AssertExpectations(t)
		})
	}
}

func TestUseCase_GetUserByID(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		mockUser      models.User
		mockError     error
		expectedUser  models.User
		expectedError error
	}{
		{
			name:   "successful get user by id",
			userID: 1,
			mockUser: models.User{
				ID:        1,
				FirstName: "",
				LastName:  "",
				Email:     "john@example.com",
				Login:     "johndoe",
				CreatedAt: time.Now(),
			},
			mockError: nil,
			expectedUser: models.User{
				ID:        1,
				FirstName: "",
				LastName:  "",
				Email:     "john@example.com",
				Login:     "johndoe",
				CreatedAt: time.Now(),
			},
			expectedError: nil,
		},
		{
			name:          "user not found",
			userID:        999,
			mockUser:      models.User{},
			mockError:     errors.New("user not found"),
			expectedUser:  models.User{},
			expectedError: errors.New("user not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthService := &mocks.AuthService{}
			uc := NewUseCase(mockAuthService, clock.RealClock{})

			mockAuthService.On("GetUserByID", mock.Anything, tt.userID).Return(tt.mockUser, tt.mockError)

			user, err := uc.GetUserByID(context.Background(), tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
				assert.Equal(t, tt.expectedUser.FirstName, user.FirstName)
				assert.Equal(t, tt.expectedUser.LastName, user.LastName)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.Equal(t, tt.expectedUser.Login, user.Login)
			}

			mockAuthService.AssertExpectations(t)
		})
	}
}

func TestUseCase_EditUserByID(t *testing.T) {
	tests := []struct {
		name          string
		req           models.UpdateUserRequest
		userID        int
		mockUser      models.User
		mockError     error
		expectedUser  models.User
		expectedError error
	}{
		{
			name: "successful edit user",
			req: models.UpdateUserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
			},
			userID: 1,
			mockUser: models.User{
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				Login:     "johndoe",
				CreatedAt: time.Now(),
			},
			mockError: nil,
			expectedUser: models.User{
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				Login:     "johndoe",
				CreatedAt: time.Now(),
			},
			expectedError: nil,
		},
		{
			name: "edit user email conflict",
			req: models.UpdateUserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
			},
			userID:        1,
			mockUser:      models.User{},
			mockError:     errors.New("email already exists"),
			expectedUser:  models.User{},
			expectedError: errors.New("email already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthService := &mocks.AuthService{}
			uc := NewUseCase(mockAuthService, clock.RealClock{})

			mockAuthService.On("EditUserByID", mock.Anything, tt.req, tt.userID).
				Return(tt.mockUser, tt.mockError)

			user, err := uc.EditUserByID(context.Background(), tt.req, tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
				assert.Equal(t, tt.expectedUser.FirstName, user.FirstName)
				assert.Equal(t, tt.expectedUser.LastName, user.LastName)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.Equal(t, tt.expectedUser.Login, user.Login)
			}

			mockAuthService.AssertExpectations(t)
		})
	}
}
