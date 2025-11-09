package profile

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type ProfileService interface {
	GetProfile(ctx context.Context, userID int) (models.ProfileResponse, error)
	UpdateProfile(ctx context.Context, req models.UpdateProfileRequest, userID int) (models.ProfileResponse, error)
}

type ProfileRepository interface {
	GetUserByID(ctx context.Context, id int) (models.User, error)
	UpdateUser(ctx context.Context, user models.User) error
}
