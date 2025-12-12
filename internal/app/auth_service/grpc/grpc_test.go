package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	svcerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/errors"
	authmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/models"
	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
)

func TestAuthServiceServer_RegisterSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockAuthUseCase(ctrl)
	server := NewAuthServer(uc)
	req := &authpb.RegisterRequest{Email: "user@example.com", Password: "pass", Login: "user"}
	expected := &authpb.AuthResponse{Token: "token"}

	uc.EXPECT().
		Register(gomock.Any(), gomock.AssignableToTypeOf(authmodels.RegisterRequest{})).
		DoAndReturn(func(_ context.Context, r authmodels.RegisterRequest) (*authpb.AuthResponse, error) {
			require.Equal(t, req.Email, r.Email)
			require.Equal(t, req.Login, r.Login)
			return expected, nil
		})

	resp, err := server.Register(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, expected, resp)
}

func TestAuthServiceServer_RegisterKnownError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockAuthUseCase(ctrl)
	server := NewAuthServer(uc)
	req := &authpb.RegisterRequest{}

	uc.EXPECT().
		Register(gomock.Any(), gomock.AssignableToTypeOf(authmodels.RegisterRequest{})).
		Return(nil, svcerrors.ErrUserExists)

	resp, err := server.Register(context.Background(), req)
	require.Nil(t, resp)
	st, _ := status.FromError(err)
	require.Equal(t, codes.AlreadyExists, st.Code())
}

func TestAuthServiceServer_LoginInternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockAuthUseCase(ctrl)
	server := NewAuthServer(uc)
	req := &authpb.LoginRequest{}

	uc.EXPECT().
		Login(gomock.Any(), gomock.AssignableToTypeOf(authmodels.LoginRequest{})).
		Return(nil, errors.New("boom"))

	resp, err := server.Login(context.Background(), req)
	require.Nil(t, resp)
	st, _ := status.FromError(err)
	require.Equal(t, codes.Internal, st.Code())
}

func TestAuthServiceServer_UpdateProfileErrorMap(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockAuthUseCase(ctrl)
	server := NewAuthServer(uc)
	req := &authpb.UpdateProfileRequest{}

	uc.EXPECT().
		UpdateProfile(gomock.Any(), gomock.AssignableToTypeOf(authmodels.UpdateProfileRequest{})).
		Return(nil, svcerrors.ErrForbidden)

	resp, err := server.UpdateProfile(context.Background(), req)
	require.Nil(t, resp)
	st, _ := status.FromError(err)
	require.Equal(t, codes.PermissionDenied, st.Code())
}

func TestAuthServiceServer_GetProfileSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockAuthUseCase(ctrl)
	server := NewAuthServer(uc)
	req := &authpb.UserID{UserID: 42}
	expected := &authpb.ProfileResponse{Id: 42}

	uc.EXPECT().GetProfile(gomock.Any(), 42).Return(expected, nil)

	resp, err := server.GetProfile(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, expected, resp)
}

func TestAuthServiceServer_GetCSRFError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockAuthUseCase(ctrl)
	server := NewAuthServer(uc)

	uc.EXPECT().GetCSRFToken(gomock.Any()).Return(nil, errors.New("boom"))

	resp, err := server.GetCSRF(context.Background(), nil)
	require.Nil(t, resp)
	st, _ := status.FromError(err)
	require.Equal(t, codes.Internal, st.Code())
}