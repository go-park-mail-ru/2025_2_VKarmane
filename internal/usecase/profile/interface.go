package profile

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type ProfileUseCase interface {
	GetProfile(ctx context.Context, userID int) (models.ProfileResponse, error)
	UpdateProfile(ctx context.Context, req models.UpdateProfileRequest, userID int) (models.ProfileResponse, error)
}
