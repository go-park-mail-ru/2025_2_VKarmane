package auth

import (
	"context"

	auth "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/models"
	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/proto"
)

type AuthUseCase interface {
    Register(context.Context, auth.RegisterRequest) (*authpb.AuthResponse, error)
    Login(context.Context, auth.LoginRequest) (*authpb.AuthResponse, error)
	GetProfile(context.Context, int) (*authpb.ProfileResponse, error)
	UpdateProfile(context.Context, auth.UpdateProfileRequest) (*authpb.ProfileResponse, error)
	GetCSRFToken(context.Context) (*authpb.CSRFTokenResponse, error)
}
