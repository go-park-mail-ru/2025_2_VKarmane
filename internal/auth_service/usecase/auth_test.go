package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	authmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUseCase_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockAuthService(ctrl)

	fixed := clock.FixedClock{
		FixedTime: time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
	}

	uc := NewAuthUseCase(mockSvc, "secret", fixed)

	req := authmodels.RegisterRequest{Login: "aaa"}
	resp := &proto.AuthResponse{Token: "ok"}

	mockSvc.EXPECT().
		Register(gomock.Any(), req).
		Return(resp, nil)

	r, err := uc.Register(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, resp, r)
}

func TestUseCase_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockAuthService(ctrl)

	fixed := clock.FixedClock{FixedTime: time.Now()}

	uc := NewAuthUseCase(mockSvc, "secret", fixed)

	req := authmodels.LoginRequest{Login: "aaa"}
	resp := &proto.AuthResponse{Token: "token123"}

	mockSvc.EXPECT().
		Login(gomock.Any(), req).
		Return(resp, nil)

	r, err := uc.Login(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, resp, r)
}

func TestUseCase_GetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockAuthService(ctrl)

	fixed := clock.FixedClock{FixedTime: time.Now()}

	uc := NewAuthUseCase(mockSvc, "secret", fixed)

	resp := &proto.ProfileResponse{Id: 10}

	mockSvc.EXPECT().
		GetUserByID(gomock.Any(), 10).
		Return(resp, nil)

	r, err := uc.GetProfile(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, resp, r)
}

func TestUseCase_UpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockAuthService(ctrl)

	fixed := clock.FixedClock{FixedTime: time.Now()}

	uc := NewAuthUseCase(mockSvc, "secret", fixed)

	req := authmodels.UpdateProfileRequest{FirstName: "John"}
	resp := &proto.ProfileResponse{FirstName: "John"}

	mockSvc.EXPECT().
		EditUserByID(gomock.Any(), req).
		Return(resp, nil)

	r, err := uc.UpdateProfile(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, resp, r)
}

func TestUseCase_GetCSRFToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fixedTime := time.Date(2030, 5, 20, 10, 0, 0, 0, time.UTC)
	fixed := clock.FixedClock{FixedTime: fixedTime}

	mockSvc := mocks.NewMockAuthService(ctrl)

	uc := NewAuthUseCase(mockSvc, "secret", fixed)

	resp, err := uc.GetCSRFToken(context.Background())
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Token)
}

func TestUseCase_Register_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockAuthService(ctrl)
	fixed := clock.FixedClock{FixedTime: time.Now()}

	uc := NewAuthUseCase(mockSvc, "secret", fixed)

	mockSvc.EXPECT().
		Register(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("fail"))

	r, err := uc.Register(context.Background(), authmodels.RegisterRequest{})
	require.Nil(t, r)
	require.Error(t, err)
}
