package auth

import (
	"context"

	authmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/models"
	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	
	pkgerrors "github.com/pkg/errors"
)

type UseCase struct {
	authService AuthService
	jwtSecret string
	clck clock.Clock
}

func NewAuthUseCase(svc AuthService, secret string, clck clock.Clock) *UseCase {
	return &UseCase{
		authService: svc,
		jwtSecret: secret,
		clck: clck,
	}
}


func (uc *UseCase) Register(ctx context.Context, req authmodels.RegisterRequest) (*authpb.AuthResponse, error) {
	user, err := uc.authService.Register(ctx, req)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UseCase) Login(ctx context.Context, req authmodels.LoginRequest) (*authpb.AuthResponse, error) {
	user, err := uc.authService.Login(ctx, req)
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (uc *UseCase) GetProfile(ctx context.Context, userID int) (*authpb.ProfileResponse, error) {
	return uc.authService.GetUserByID(ctx, userID)
}

func (uc *UseCase) UpdateProfile(ctx context.Context, req authmodels.UpdateProfileRequest) (*authpb.ProfileResponse, error) {
	return uc.authService.EditUserByID(ctx, req)
}

func (uc *UseCase) GetCSRFToken(ctx context.Context) (*authpb.CSRFTokenResponse, error) {
	clock := clock.RealClock{}

	token, err := utils.GenerateCSRF(clock.Now(), uc.jwtSecret)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "auth.GenerateCSRF")
	}
	return &authpb.CSRFTokenResponse{Token: token}, nil
}
