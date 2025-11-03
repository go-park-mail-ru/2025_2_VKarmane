package profile

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/profile"
)

type UseCase struct {
	profileSvc profile.ProfileService
}

func NewUseCase(profileSvc profile.ProfileService) *UseCase {
	return &UseCase{
		profileSvc: profileSvc,
	}
}

func (uc *UseCase) GetProfile(ctx context.Context, userID int) (models.ProfileResponse, error) {
	log := logger.FromContext(ctx)
	profile, err := uc.profileSvc.GetProfile(ctx, userID)
	if err != nil {
		log.Error("Failed to get profile", "error", err, "user_id", userID)
		return models.ProfileResponse{}, fmt.Errorf("profile.GetProfile: %w", err)
	}

	return profile, nil
}

func (uc *UseCase) UpdateProfile(ctx context.Context, req models.UpdateProfileRequest, userID int) (models.ProfileResponse, error) {
	log := logger.FromContext(ctx)
	profile, err := uc.profileSvc.UpdateProfile(ctx, req, userID)
	if err != nil {
		log.Error("Failed to update profile", "error", err, "user_id", userID)
		return models.ProfileResponse{}, fmt.Errorf("profile.UpdateProfile: %w", err)
	}

	return profile, nil
}
