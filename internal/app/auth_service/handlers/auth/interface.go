package auth

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type AuthUseCase interface {
	Register(ctx context.Context, req models.RegisterRequest) (models.AuthResponse, error)
	Login(ctx context.Context, req models.LoginRequest) (models.AuthResponse, error)
	Logout(ctx context.Context, w http.ResponseWriter)
	GetUserByID(ctx context.Context, userID int) (models.User, error)
	EditUserByID(ctx context.Context, req models.UpdateProfileRequest, userID int) (models.User, error)
	GetCSRFToken(ctx context.Context) (string, error)
}
