package auth

import (
	"context"

	authmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/models"
	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/proto"
)

type AuthService interface {
	Register(ctx context.Context, req authmodels.RegisterRequest) (*authpb.AuthResponse, error)
	Login(ctx context.Context, req authmodels.LoginRequest) (*authpb.AuthResponse, error)
	GetUserByID(ctx context.Context, userID int) (*authpb.ProfileResponse, error)
	EditUserByID(ctx context.Context, req authmodels.UpdateProfileRequest) (*authpb.ProfileResponse, error)
}
