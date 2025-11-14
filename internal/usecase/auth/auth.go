package auth

import (
	"context"
	"net/http"
	"os"

	pkgerrors "github.com/pkg/errors"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

type UseCase struct {
	authService AuthService
	clock       clock.Clock
	jwtSecret	string
}

func NewUseCase(authService AuthService, clck clock.Clock, secret string) *UseCase {
	return &UseCase{
		authService: authService,
		clock:       clck,
		jwtSecret: secret,
	}
}

func (uc *UseCase) Register(ctx context.Context, req models.RegisterRequest) (models.AuthResponse, error) {
	return uc.authService.Register(ctx, req)
}

func (uc *UseCase) Login(ctx context.Context, req models.LoginRequest) (models.AuthResponse, error) {
	return uc.authService.Login(ctx, req)
}

func (uc *UseCase) Logout(ctx context.Context, w http.ResponseWriter) {
	isProduction := os.Getenv("ENV") == "production"
	utils.ClearAuthCookie(w, isProduction)
	utils.ClearCSRFCookie(w, isProduction)
}

func (uc *UseCase) GetUserByID(ctx context.Context, userID int) (models.User, error) {
	return uc.authService.GetUserByID(ctx, userID)
}

func (uc *UseCase) EditUserByID(ctx context.Context, req models.UpdateProfileRequest, userID int) (models.User, error) {
	return uc.authService.EditUserByID(ctx, req, userID)
}

func (uc *UseCase) GetCSRFToken(ctx context.Context) (string, error) {
	clock := clock.RealClock{}

	token, err := utils.GenerateCSRF(clock.Now(), uc.jwtSecret)
	if err != nil {
		return "", pkgerrors.Wrap(err, "auth.GenerateCSRF")
	}
	return token, nil
}
