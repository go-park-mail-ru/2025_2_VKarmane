package auth

import (
	"context"

	authmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/models"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user authmodels.User) (authmodels.User, error)
	GetUserByLogin(ctx context.Context, login string) (authmodels.User, error)
	GetUserByID(ctx context.Context, id int) (authmodels.User, error)
	EditUserByID(ctx context.Context, req authmodels.UpdateProfileRequest) (authmodels.User, error)
}
