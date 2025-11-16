package profile

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)


func IDtoProtoID(id int) *proto.UserID {
	return &proto.UserID{
		UserID: int32(id),
	}
}

func ProtoProfileToApiProfile(req *proto.ProfileResponse) *models.ProfileResponse {
	return &models.ProfileResponse{
		ID: int(req.Id),
		FirstName: req.FirstName,
		LastName: req.LastName,
		Login: req.Login,
		Email: req.Email,
		LogoHashedID: req.LogoHashedId,
		LogoURL: req.LogoUrl,
		CreatedAt: req.CreatedAt.AsTime(),
	}
}

func UpdateProfileApiToProto(req models.UpdateProfileRequest, id int) *proto.UpdateProfileRequest {
	return &proto.UpdateProfileRequest{
		UserId: int32(id),
		FirstName: req.FirstName,
		LastName: req.LastName,
		Email: req.Email,
		LogoHashedId: req.LogoHashedID,
	}
}
