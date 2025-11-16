package auth

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

// func UserToApi(op models.User) UserAPI {
// 	return UserAPI{
// 		ID:        op.ID,
// 		FirstName: op.FirstName,
// 		LastName:  op.LastName,
// 		Email:     op.Email,
// 		Login:     op.Login,
// 		Password:  op.Password,
// 		CreatedAt: op.CreatedAt,
// 		UpdatedAt: op.UpdatedAt,
// 	}
// }

func LoginApiToProtoLogin(req models.LoginRequest) *proto.LoginRequest {
	return &proto.LoginRequest{
		Login: req.Login,
		Password: req.Password,
	}
}

func RegisterApiToProtoRegister(req models.RegisterRequest) *proto.RegisterRequest {
	return &proto.RegisterRequest{
		Login: req.Login,
		Email: req.Email,
		Password: req.Password,
	}
}
