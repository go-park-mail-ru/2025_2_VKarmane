package profile

import (
	"context"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type Service struct {
	repo interface {
		GetUserByID(ctx context.Context, id int) (models.User, error)
		UpdateUser(ctx context.Context, user models.User) error
	}
}

func NewService(repo interface {
	GetUserByID(ctx context.Context, id int) (models.User, error)
	UpdateUser(ctx context.Context, user models.User) error
}) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetProfile(ctx context.Context, userID int) (models.ProfileResponse, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
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
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return models.ProfileResponse{}, fmt.Errorf("failed to get user: %w", err)
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Email = req.Email
	user.UpdatedAt = time.Now()

	err = s.repo.UpdateUser(ctx, user)
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
