package user

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/dto"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user dto.UserDB) (int, error)
	GetUserByEmail(ctx context.Context, email string) (dto.UserDB, error)
	GetUserByLogin(ctx context.Context, login string) (dto.UserDB, error)
	GetUserByID(ctx context.Context, id int) (dto.UserDB, error)
	UpdateUser(ctx context.Context, user dto.UserDB) error
}
