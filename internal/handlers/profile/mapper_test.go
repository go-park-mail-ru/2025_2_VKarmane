package profile

import (
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestIDtoProtoID(t *testing.T) {
	res := IDtoProtoID(42)

	require.NotNil(t, res)
	require.Equal(t, int32(42), res.UserID)
}

func TestProtoProfileToApiProfile(t *testing.T) {
	now := time.Now()

	protoProfile := &proto.ProfileResponse{
		Id:            10,
		FirstName:     "John",
		LastName:      "Doe",
		Login:         "jdoe",
		Email:         "john@example.com",
		LogoHashedId:  "hash123",
		LogoUrl:       "http://example.com/logo.jpg",
		CreatedAt:     timestamppb.New(now),
	}

	api := ProtoProfileToApiProfile(protoProfile)

	require.NotNil(t, api)
	require.Equal(t, 10, api.ID)
	require.Equal(t, "John", api.FirstName)
	require.Equal(t, "Doe", api.LastName)
	require.Equal(t, "jdoe", api.Login)
	require.Equal(t, "john@example.com", api.Email)
	require.Equal(t, "hash123", api.LogoHashedID)
	require.Equal(t, "http://example.com/logo.jpg", api.LogoURL)
	require.True(t, api.CreatedAt.Equal(now))
}

func TestUpdateProfileApiToProto(t *testing.T) {
	req := models.UpdateProfileRequest{
		FirstName:    "Jane",
		LastName:     "Smith",
		Email:        "jane@example.com",
		LogoHashedID: "abc123",
	}

	res := UpdateProfileApiToProto(req, 7)

	require.NotNil(t, res)
	require.Equal(t, int32(7), res.UserId)
	require.Equal(t, "Jane", res.FirstName)
	require.Equal(t, "Smith", res.LastName)
	require.Equal(t, "jane@example.com", res.Email)
	require.Equal(t, "abc123", res.LogoHashedId)
}
