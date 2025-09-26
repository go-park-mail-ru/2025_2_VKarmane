package auth

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/auth"
)

type UseCase struct {
	authService *auth.Service
}

func NewUseCase(authService *auth.Service) *UseCase {
	return &UseCase{
		authService: authService,
	}
}

func (uc *UseCase) Register(req models.RegisterRequest) (models.AuthResponse, error) {
	return uc.authService.Register(req)
}

func (uc *UseCase) Login(req models.LoginRequest) (models.AuthResponse, error) {
	return uc.authService.Login(req)
}

func (uc *UseCase) GetUserByID(userID int) (models.User, error) {
	return uc.authService.GetUserByID(userID)
}
