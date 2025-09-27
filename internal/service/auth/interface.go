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

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
	GetUserByID(ctx context.Context, id int) (models.User, error)
}
