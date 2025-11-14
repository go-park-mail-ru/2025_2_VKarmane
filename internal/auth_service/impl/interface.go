package auth

import (
	"context"

	auth "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/models"
	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthUseCase interface {
    Register(context.Context, auth.RegisterRequest) (*authpb.AuthResponse, error)
    Login(context.Context, auth.LoginRequest) (*authpb.AuthResponse, error)
    Logout(context.Context) (*emptypb.Empty, error)
	GetProfile(context.Context, int) (*authpb.ProfileResponse, error)
	UpdateProfile(context.Context, auth.UpdateProfileRequest) (*authpb.ProfileResponse, error)
	GetCSRFToken(ctx context.Context) (*authpb.CSRFTokenResponse, error)
}
