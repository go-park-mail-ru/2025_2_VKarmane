package auth

import (
	authmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/models"
	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/proto"
)


func RegisterToRequest(req *authpb.RegisterRequest) authmodels.RegisterRequest {
	return authmodels.RegisterRequest {
		Email: req.Email,
		Password: req.Password,
		Login: req.Login,
	}
}

func LoginToRequest(req *authpb.LoginRequest) authmodels.LoginRequest {
	return authmodels.LoginRequest {
		Password: req.Password,
		Login: req.Login,
	}
}

func UpdateProfileToRequest(req *authpb.UpdateProfileRequest) authmodels.UpdateProfileRequest {
	return authmodels.UpdateProfileRequest{
		UserID: req.UserId,
		FirstName: req.FirstName,
		LastName: req.LastName,
		Email: req.Email,
		LogoHashedID: req.LogoHashedId,
	}
}

func ProtoIDtoInt(id *authpb.UserID) int {
	return int(id.UserID)
}

