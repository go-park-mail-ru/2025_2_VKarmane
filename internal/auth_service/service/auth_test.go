package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	svcerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/errors"
	authmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/models"
	mock "github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	"github.com/stretchr/testify/require"

	"go.uber.org/mock/gomock"
)

func TestService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockAuthRepository(ctrl)
	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	s := NewService(repo, "secret", fixedClock)

	req := authmodels.RegisterRequest{
		Email:    "test@example.com",
		Login:    "testuser",
		Password: "password123",
	}

	createdUser := authmodels.User{
		ID:        1,
		Login:     req.Login,
		Email:     req.Email,
		Password:  "hashedpassword",
		FirstName: "User",
		LastName:  "User",
		CreatedAt: fixedClock.FixedTime,
		UpdatedAt: fixedClock.FixedTime,
	}

	repo.EXPECT().
		CreateUser(gomock.Any(), gomock.Any()).
		Return(createdUser, nil)

	resp, err := s.Register(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, createdUser.ID, int(resp.User.Id))
	require.Equal(t, createdUser.Login, resp.User.Login)
}

func TestService_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockAuthRepository(ctrl)
	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	s := NewService(repo, "secret", fixedClock)

	hashed, _ := utils.HashPassword("password123")
	user := authmodels.User{
		ID:       1,
		Login:    "testuser",
		Email:    "test@example.com",
		Password: hashed,
	}

	req := authmodels.LoginRequest{
		Login:    "testuser",
		Password: "password123",
	}

	repo.EXPECT().GetUserByLogin(gomock.Any(), "testuser").Return(user, nil)



	resp, err := s.Login(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, user.ID, int(resp.User.Id))
	require.Equal(t, user.Login, resp.User.Login)
}

func TestService_Login_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockAuthRepository(ctrl)
	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	s := NewService(repo, "secret", fixedClock)

	hashed, _ := utils.HashPassword("correctpassword")

	user := authmodels.User{
		ID:       1,
		Login:    "testuser",
		Email:    "test@example.com",
		Password: hashed, // валидный bcrypt-хеш
	}

	req := authmodels.LoginRequest{
		Login:    "testuser",
		Password: "wrongpassword", // неверный пароль
	}

	repo.EXPECT().GetUserByLogin(gomock.Any(), "testuser").Return(user, nil)

	// В этом тесте utils.VerifyPassword вернет false
	// Для юнит-теста можно временно подменить или мокировать utils.VerifyPassword

	_, err := s.Login(context.Background(), req)
	require.Error(t, err)
	require.True(t, errors.Is(err, svcerrors.ErrInvalidCredentials))
}

func TestService_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockAuthRepository(ctrl)
	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	s := NewService(repo, "secret", fixedClock)

	user := authmodels.User{
		ID:        1,
		Login:     "testuser",
		Email:     "test@example.com",
		FirstName: "User",
		LastName:  "User",
		CreatedAt: fixedClock.FixedTime,
		UpdatedAt: fixedClock.FixedTime,
	}

	repo.EXPECT().GetUserByID(gomock.Any(), 1).Return(user, nil)

	resp, err := s.GetUserByID(context.Background(), 1)
	require.NoError(t, err)
	require.Equal(t, int32(1), resp.Id)
	require.Equal(t, "testuser", resp.Login)
}

func TestService_EditUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockAuthRepository(ctrl)
	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	s := NewService(repo, "secret", fixedClock)

	req := authmodels.UpdateProfileRequest{
		UserID:    1,
		FirstName: "New",
		LastName:  "Name",
		Email:     "new@example.com",
	}

	updatedUser := authmodels.User{
		ID:        1,
		Login:     "testuser",
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		CreatedAt: fixedClock.FixedTime,
		UpdatedAt: fixedClock.FixedTime,
	}

	repo.EXPECT().EditUserByID(gomock.Any(), req).Return(updatedUser, nil)

	resp, err := s.EditUserByID(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, int32(1), resp.Id)
	require.Equal(t, "New", resp.FirstName)
	require.Equal(t, "Name", resp.LastName)
	require.Equal(t, "new@example.com", resp.Email)
}
