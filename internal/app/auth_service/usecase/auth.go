package auth

import (
	"context"

	pkgerrors "github.com/pkg/errors"

	authmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/models"
	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

type UseCase struct {
	authService AuthService
	jwtSecret   string
	clck        clock.Clock
}

func NewAuthUseCase(svc AuthService, secret string, clck clock.Clock) *UseCase {
	return &UseCase{
		authService: svc,
		jwtSecret:   secret,
		clck:        clck,
	}
}

func (uc *UseCase) Register(ctx context.Context, req authmodels.RegisterRequest) (*authpb.AuthResponse, error) {
	log := logger.FromContext(ctx)
	user, err := uc.authService.Register(ctx, req)
	if err != nil {
		if log != nil {
			log.Error("Failed to register user", "error", err)
		}
		return nil, err
	}

	return user, nil
}

func (uc *UseCase) Login(ctx context.Context, req authmodels.LoginRequest) (*authpb.AuthResponse, error) {
	log := logger.FromContext(ctx)
	user, err := uc.authService.Login(ctx, req)
	if err != nil {
		if log != nil {
			log.Error("Failed to login user", "error", err)
		}
		return nil, err
	}

	return user, nil
}

func (uc *UseCase) GetProfile(ctx context.Context, userID int) (*authpb.ProfileResponse, error) {
	log := logger.FromContext(ctx)
	profile, err := uc.authService.GetUserByID(ctx, userID)
	if err != nil {
		if log != nil {
			log.Error("Failed to get profile", "error", err)
		}
		return nil, err
	}
	return profile, nil
}

func (uc *UseCase) UpdateProfile(ctx context.Context, req authmodels.UpdateProfileRequest) (*authpb.ProfileResponse, error) {
	log := logger.FromContext(ctx)
	profile, err := uc.authService.EditUserByID(ctx, req)
	if err != nil {
		if log != nil {
			log.Error("Failed to update profile", "error", err)
		}
		return nil, err
	}
	return profile, nil
}

func (uc *UseCase) GetCSRFToken(ctx context.Context) (*authpb.CSRFTokenResponse, error) {
	clock := clock.RealClock{}
	log := logger.FromContext(ctx)

	token, err := utils.GenerateCSRF(clock.Now(), uc.jwtSecret)
	if err != nil {
		if log != nil {
			log.Error("Failed to get CSRF-Token", "error", err)
		}
		return nil, pkgerrors.Wrap(err, "auth.GenerateCSRF")
	}
	return &authpb.CSRFTokenResponse{Token: token}, nil
}
