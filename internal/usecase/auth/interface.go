package auth

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type AuthService interface {
	Register(ctx context.Context, req models.RegisterRequest) (models.AuthResponse, error)
	Login(ctx context.Context, req models.LoginRequest) (models.AuthResponse, error)
	GetUserByID(ctx context.Context, userID int) (models.User, error)
}
