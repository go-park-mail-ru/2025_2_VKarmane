package auth

import (
	"testing"
	"time"

	authmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/models"
	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/proto"

	"github.com/stretchr/testify/require"
)

func TestModelUserToProtoUser(t *testing.T) {
	created := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	updated := time.Date(2025, 1, 2, 15, 30, 0, 0, time.UTC)

	user := authmodels.User{
		ID:           10,
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "john@example.com",
		Description:  "test user",
		LogoHashedID: "abc123",
		CreatedAt:    created,
		UpdatedAt:    updated,
	}

	pb := ModelUserToProtoUser(user)

	require.IsType(t, &authpb.User{}, pb)

	require.Equal(t, int32(10), pb.Id)
	require.Equal(t, "John", pb.FirstName)
	require.Equal(t, "Doe", pb.LastName)
	require.Equal(t, "john@example.com", pb.Email)
	require.Equal(t, "test user", pb.Description)
	require.Equal(t, "abc123", pb.LogoHashedId)

	require.Equal(t, created, pb.CreatedAt.AsTime())
	require.Equal(t, updated, pb.UpdatedAt.AsTime())
}

func TestModelUserToProfile(t *testing.T) {
	user := authmodels.User{
		ID:           20,
		FirstName:    "Alice",
		LastName:     "Smith",
		Email:        "alice@example.com",
		Login:        "alice",
		LogoHashedID: "logo321",
	}

	pb := ModelUserToProfile(user)

	require.IsType(t, &authpb.ProfileResponse{}, pb)

	require.Equal(t, int32(20), pb.Id)
	require.Equal(t, "Alice", pb.FirstName)
	require.Equal(t, "Smith", pb.LastName)
	require.Equal(t, "alice@example.com", pb.Email)
	require.Equal(t, "alice", pb.Login)
	require.Equal(t, "logo321", pb.LogoHashedId)
}
