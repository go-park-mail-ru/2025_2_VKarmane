package auth

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	authmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/models"
	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/proto"
)


func ModelUserToProtoUser(user authmodels.User) *authpb.User {
	return &authpb.User{
		Id: int32(user.ID),
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email: user.Email,
		Description: user.Description,
		LogoHashedId: user.LogoHashedID,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func ModelUserToProfile(user authmodels.User) *authpb.ProfileResponse {
	return &authpb.ProfileResponse {
		Id: int32(user.ID),
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email: user.Email,
		Login: user.Login,
		LogoHashedId: user.LogoHashedID,
	}
}