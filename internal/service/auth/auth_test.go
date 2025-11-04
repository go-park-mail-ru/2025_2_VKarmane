package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestService_Register(t *testing.T) {
	realClock := clock.RealClock{}
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
				Email:    "john@example.com",
				Login:    "johndoe",
				Password: "password123",
			},
			mockUser: models.User{
				ID:        1,
				FirstName: "",
				LastName:  "",
				Email:     "john@example.com",
				Login:     "johndoe",
				Password:  "hashed_password",
				CreatedAt: realClock.Now(),
			},
			mockError: nil,
			expectedResult: models.AuthResponse{
				Token: "jwt-token",
				User: models.User{
					ID:        1,
					FirstName: "",
					LastName:  "",
					Email:     "john@example.com",
					Login:     "johndoe",
					Password:  "hashed_password",
					CreatedAt: realClock.Now(),
				},
			},
			expectedError: "",
		},
		{
			name: "user creation error",
			request: models.RegisterRequest{
				Email:    "john@example.com",
				Login:    "johndoe",
				Password: "password123",
			},
			mockUser:       models.User{},
			mockError:      errors.New("user already exists"),
			expectedResult: models.AuthResponse{},
			expectedError:  "auth.Register: failed to create user: user already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mocks.NewMockUserRepository(ctrl)
			service := NewService(mockUserRepo, "test-secret", realClock)

			mockUserRepo.EXPECT().
				CreateUser(gomock.Any(), gomock.Any()).
				Return(tt.mockUser, tt.mockError)

			ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
			result, err := service.Register(ctx, tt.request)

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
		})
	}
}

func TestService_Login(t *testing.T) {
	realClock := clock.RealClock{}
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mocks.NewMockUserRepository(ctrl)
			service := NewService(mockUserRepo, "test-secret", realClock)

			mockUserRepo.EXPECT().
				GetUserByLogin(gomock.Any(), tt.request.Login).
				Return(tt.mockUser, tt.mockError)

			ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
			result, err := service.Login(ctx, tt.request)

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
		})
	}
}

func TestService_GetUserByID(t *testing.T) {
	realClock := clock.RealClock{}
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
				FirstName: "",
				LastName:  "",
				Email:     "john@example.com",
				Login:     "johndoe",
				CreatedAt: realClock.Now(),
			},
			mockError: nil,
			expectedUser: models.User{
				ID:        1,
				FirstName: "",
				LastName:  "",
				Email:     "john@example.com",
				Login:     "johndoe",
				CreatedAt: realClock.Now(),
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mocks.NewMockUserRepository(ctrl)
			service := NewService(mockUserRepo, "test-secret", realClock)

			mockUserRepo.EXPECT().
				GetUserByID(gomock.Any(), tt.userID).
				Return(tt.mockUser, tt.mockError)

			ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
			user, err := service.GetUserByID(ctx, tt.userID)

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
		})
	}
}

func TestService_EditUserByID(t *testing.T) {
	realClock := clock.RealClock{}
	tests := []struct {
		name          string
		req           models.UpdateProfileRequest
		userID        int
		mockUser      models.User
		mockError     error
		expectedUser  models.User
		expectedError string
	}{
		{
			name: "successful edit user",
			req: models.UpdateProfileRequest{
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
				CreatedAt: realClock.Now(),
			},
			mockError: nil,
			expectedUser: models.User{
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				Login:     "johndoe",
			},
			expectedError: "",
		},
		{
			name: "edit user conflict error",
			req: models.UpdateProfileRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
			},
			userID:        1,
			mockUser:      models.User{},
			mockError:     errors.New("email already exists"),
			expectedUser:  models.User{},
			expectedError: "auth.EditUserByID: email already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mocks.NewMockUserRepository(ctrl)
			service := NewService(mockUserRepo, "test-secret", realClock)

			mockUserRepo.EXPECT().
				EditUserByID(gomock.Any(), tt.req, tt.userID).
				Return(tt.mockUser, tt.mockError)

			ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
			user, err := service.EditUserByID(ctx, tt.req, tt.userID)

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
		})
	}
}
