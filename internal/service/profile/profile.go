package profile

import (
	"context"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type Service struct {
	profileRepo ProfileRepository
}

func NewService(profileRepo ProfileRepository) *Service {
	return &Service{
		profileRepo: profileRepo,
	}
}

func (s *Service) GetProfile(ctx context.Context, userID int) (models.ProfileResponse, error) {
	user, err := s.profileRepo.GetUserByID(ctx, userID)
	if err != nil {
		return models.ProfileResponse{}, fmt.Errorf("failed to get user: %w", err)
	}

	return models.ProfileResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Login:     user.Login,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *Service) UpdateProfile(ctx context.Context, req models.UpdateProfileRequest, userID int) (models.ProfileResponse, error) {
	user, err := s.profileRepo.GetUserByID(ctx, userID)
	if err != nil {
		return models.ProfileResponse{}, fmt.Errorf("failed to get user: %w", err)
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Email = req.Email
	user.UpdatedAt = time.Now()

	err = s.profileRepo.UpdateUser(ctx, user)
	if err != nil {
		return models.ProfileResponse{}, fmt.Errorf("failed to update user: %w", err)
	}

	return models.ProfileResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Login:     user.Login,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}
