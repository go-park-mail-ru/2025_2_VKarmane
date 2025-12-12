package grpc

import (
	"testing"

	"github.com/stretchr/testify/require"

	authmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/models"
	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/proto"
)

func TestRegisterToRequest(t *testing.T) {
	req := &authpb.RegisterRequest{
		Email:    "user@example.com",
		Password: "secret",
		Login:    "user",
	}
	result := RegisterToRequest(req)

	require.Equal(t, req.Email, result.Email)
	require.Equal(t, req.Password, result.Password)
	require.Equal(t, req.Login, result.Login)
}

func TestLoginToRequest(t *testing.T) {
	req := &authpb.LoginRequest{Login: "user", Password: "secret"}
	result := LoginToRequest(req)

	require.Equal(t, req.Login, result.Login)
	require.Equal(t, req.Password, result.Password)
}

func TestUpdateProfileToRequest(t *testing.T) {
	req := &authpb.UpdateProfileRequest{
		UserId:       42,
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "john@example.com",
		LogoHashedId: "hash",
	}
	result := UpdateProfileToRequest(req)

	require.Equal(t, authmodels.UpdateProfileRequest{
		UserID:       req.UserId,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		LogoHashedID: req.LogoHashedId,
	}, result)
}

func TestProtoIDtoInt(t *testing.T) {
	require.Equal(t, 7, ProtoIDtoInt(&authpb.UserID{UserID: 7}))
}