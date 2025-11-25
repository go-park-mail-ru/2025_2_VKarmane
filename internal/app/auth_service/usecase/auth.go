package auth

import (
	"context"

	pkgerrors "github.com/pkg/errors"

	authmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/models"
	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/proto"
	svcerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/errors"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

type UseCase struct {
	repo AuthRepository
	jwtSecret   string
	clck        clock.Clock
}

func NewAuthUseCase(repo AuthRepository, secret string, clck clock.Clock) *UseCase {
	return &UseCase{
		repo: repo,
		jwtSecret:   secret,
		clck:        clck,
	}
}

func (uc *UseCase) Register(ctx context.Context, req authmodels.RegisterRequest) (*authpb.AuthResponse, error) {
	log := logger.FromContext(ctx)
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		if log != nil {
			log.Error("Failed to hash password", "error", err)
		}
		return nil, pkgerrors.Wrap(err, "auth.Register: failed to hash password")
	}

	user := authmodels.User{
		FirstName: "User",
		LastName:  "User",
		Email:     req.Email,
		Login:     req.Login,
		Password:  hashedPassword,
	}

	createdUser, err := uc.repo.CreateUser(ctx, user)
	if err != nil {
		if log != nil {
			log.Error("Failed to create user", "error", err, "login", req.Login)
		}
		return nil, pkgerrors.Wrap(err, "auth.Register: failed to create user")
	}

	token, err := utils.GenerateJWT(createdUser.ID, createdUser.Login, uc.jwtSecret)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "auth.Register: failed to generate token")
	}

	return &authpb.AuthResponse{
		Token: token,
		User:  ModelUserToProtoUser(createdUser),
	}, nil
}

func (uc *UseCase) Login(ctx context.Context, req authmodels.LoginRequest) (*authpb.AuthResponse, error) {
	log := logger.FromContext(ctx)
	user, err := uc.repo.GetUserByLogin(ctx, req.Login)
	if err != nil {
		if log != nil {
			log.Warn("Login attempt with invalid credentials", "login", req.Login, "error", err)
		}
		return nil, pkgerrors.Wrap(err, "auth.Login: invalid credentials")
	}

	valid, err := utils.VerifyPassword(req.Password, user.Password)
	if err != nil {
		if log != nil {
			log.Error("Failed to verify password", "error", err, "user_id", user.ID)
		}
		return nil, pkgerrors.Wrap(err, "auth.Login: failed to verify password")
	}

	if !valid {
		if log != nil {
			log.Warn("Login attempt with invalid password", "login", req.Login, "user_id", user.ID)
		}
		return nil, svcerrors.ErrInvalidCredentials
	}

	token, err := utils.GenerateJWT(user.ID, user.Login, uc.jwtSecret)
	if err != nil {
		if log != nil {
			log.Error("Failed to generate JWT token", "error", err, "user_id", user.ID)
		}
		return nil, pkgerrors.Wrap(err, "auth.Login: failed to generate token")
	}

	user.Password = ""

	return &authpb.AuthResponse{
		Token: token,
		User:  ModelUserToProtoUser(user),
	}, nil
}

func (uc *UseCase) GetProfile(ctx context.Context, userID int) (*authpb.ProfileResponse, error) {
	log := logger.FromContext(ctx)
	user, err := uc.repo.GetUserByID(ctx, userID)
	if err != nil {
		if log != nil {
			log.Error("Failed to get user by ID", "error", err, "user_id", userID)
		}
		return nil, pkgerrors.Wrap(err, "auth.GetUserByID")
	}

	return ModelUserToProfile(user), nil
}

func (uc *UseCase) UpdateProfile(ctx context.Context, req authmodels.UpdateProfileRequest) (*authpb.ProfileResponse, error) {
	log := logger.FromContext(ctx)
	user, err := uc.repo.EditUserByID(ctx, req)
	if err != nil {
		if log != nil {
			log.Error("Failed to get update user by ID", "error", err, "user_id", req.UserID)
		}
		return nil, pkgerrors.Wrap(err, "auth.EditUserByID")
	}

	return ModelUserToProfile(user), nil
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


