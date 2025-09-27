package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_Register(t *testing.T) {
	tests := []struct {
		name           string
		request        models.RegisterRequest
		mockUser       models.User
		mockError      error
		expectedResult models.AuthResponse
		expectedError  string
	}{
		{
			name: "successful registration",
			request: models.RegisterRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				Login:     "johndoe",
				Password:  "password123",
			},
			mockUser: models.User{
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				Login:     "johndoe",
				Password:  "hashed_password",
				CreatedAt: time.Now(),
			},
			mockError: nil,
			expectedResult: models.AuthResponse{
				Token: "jwt-token",
				User: models.User{
					ID:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@example.com",
					Login:     "johndoe",
					Password:  "hashed_password",
					CreatedAt: time.Now(),
				},
			},
			expectedError: "",
		},
		{
			name: "user creation error",
			request: models.RegisterRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				Login:     "johndoe",
				Password:  "password123",
			},
			mockUser:       models.User{},
			mockError:      errors.New("user already exists"),
			expectedResult: models.AuthResponse{},
			expectedError:  "auth.Register: failed to create user: user already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := &mocks.UserRepository{}
			service := NewService(mockUserRepo, "test-secret")

			mockUserRepo.On("CreateUser", mock.Anything, mock.Anything).Return(tt.mockUser, tt.mockError)

			result, err := service.Register(context.Background(), tt.request)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result.Token)
				assert.Equal(t, tt.expectedResult.User.ID, result.User.ID)
				assert.Equal(t, tt.expectedResult.User.FirstName, result.User.FirstName)
				assert.Equal(t, tt.expectedResult.User.LastName, result.User.LastName)
				assert.Equal(t, tt.expectedResult.User.Email, result.User.Email)
				assert.Equal(t, tt.expectedResult.User.Login, result.User.Login)
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestService_Login(t *testing.T) {
	tests := []struct {
		name          string
		request       models.LoginRequest
		mockUser      models.User
		mockError     error
		expectedError string
	}{
		{
			name: "user not found",
			request: models.LoginRequest{
				Login:    "nonexistent",
				Password: "password123",
			},
			mockUser:      models.User{},
			mockError:     errors.New("user not found"),
			expectedError: "auth.Login: invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := &mocks.UserRepository{}
			service := NewService(mockUserRepo, "test-secret")

			mockUserRepo.On("GetUserByLogin", mock.Anything, tt.request.Login).Return(tt.mockUser, tt.mockError)

			result, err := service.Login(context.Background(), tt.request)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result.Token)
				assert.Equal(t, tt.mockUser.ID, result.User.ID)
				assert.Equal(t, tt.mockUser.FirstName, result.User.FirstName)
				assert.Equal(t, tt.mockUser.LastName, result.User.LastName)
				assert.Equal(t, tt.mockUser.Email, result.User.Email)
				assert.Equal(t, tt.mockUser.Login, result.User.Login)
				assert.Empty(t, result.User.Password)
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetUserByID(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		mockUser      models.User
		mockError     error
		expectedUser  models.User
		expectedError string
	}{
		{
			name:   "successful get user by id",
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
			expectedError: "",
		},
		{
			name:          "user not found",
			userID:        999,
			mockUser:      models.User{},
			mockError:     errors.New("user not found"),
			expectedUser:  models.User{},
			expectedError: "auth.GetUserByID: user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := &mocks.UserRepository{}
			service := NewService(mockUserRepo, "test-secret")

			mockUserRepo.On("GetUserByID", mock.Anything, tt.userID).Return(tt.mockUser, tt.mockError)

			user, err := service.GetUserByID(context.Background(), tt.userID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
				assert.Equal(t, tt.expectedUser.FirstName, user.FirstName)
				assert.Equal(t, tt.expectedUser.LastName, user.LastName)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.Equal(t, tt.expectedUser.Login, user.Login)
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}
