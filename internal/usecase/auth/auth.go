package auth

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type UseCase struct {
	authService AuthService
}

func NewUseCase(authService AuthService) *UseCase {
	return &UseCase{
		authService: authService,
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
