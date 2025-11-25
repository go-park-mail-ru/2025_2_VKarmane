package auth

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

func LoginApiToProtoLogin(req models.LoginRequest) *proto.LoginRequest {
	return &proto.LoginRequest{
		Login:    req.Login,
		Password: req.Password,
	}
}

func RegisterApiToProtoRegister(req models.RegisterRequest) *proto.RegisterRequest {
	return &proto.RegisterRequest{
		Login:    req.Login,
		Email:    req.Email,
		Password: req.Password,
	}
}
