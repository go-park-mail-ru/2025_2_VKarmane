package auth

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

type UseCase struct {
	authService AuthService
	clock       clock.Clock
}

func NewUseCase(authService AuthService, clck clock.Clock) *UseCase {
	return &UseCase{
		authService: authService,
		clock:       clck,
	}
}

func (uc *UseCase) Register(ctx context.Context, req models.RegisterRequest) (models.AuthResponse, error) {
	return uc.authService.Register(ctx, req)
}

func (uc *UseCase) Login(ctx context.Context, req models.LoginRequest) (models.AuthResponse, error) {
	return uc.authService.Login(ctx, req)
}

func (uc *UseCase) GetUserByID(ctx context.Context, userID int) (models.User, error) {
	return uc.authService.GetUserByID(ctx, userID)
}

func (uc *UseCase) EditUserByID(ctx context.Context, req models.UpdateProfileRequest, userID int) (models.User, error) {
	return uc.authService.EditUserByID(ctx, req, userID)
}
