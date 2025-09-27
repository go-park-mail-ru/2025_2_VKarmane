package auth

import "github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"

type AuthService interface {
	Register(req models.RegisterRequest) (models.AuthResponse, error)
	Login(req models.LoginRequest) (models.AuthResponse, error)
	GetUserByID(userID int) (models.User, error)
}
