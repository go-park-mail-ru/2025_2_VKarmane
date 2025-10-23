package auth

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/dto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/user"
)

type PostgresUserRepositoryAdapter struct {
	userRepo user.UserRepository
}

func NewPostgresUserRepositoryAdapter(userRepo user.UserRepository) *PostgresUserRepositoryAdapter {
	return &PostgresUserRepositoryAdapter{
		userRepo: userRepo,
	}
}

func (a *PostgresUserRepositoryAdapter) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	userDB := dto.UserDB{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Login:       user.Login,
		Password:    user.Password,
		Description: user.Description,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	id, err := a.userRepo.CreateUser(ctx, userDB)
	if err != nil {
		return models.User{}, err
	}

	user.ID = id
	return user, nil
}

func (a *PostgresUserRepositoryAdapter) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	userDB, err := a.userRepo.GetUserByLogin(ctx, login)
	if err != nil {
		return models.User{}, err
	}

	return dtoToModel(userDB), nil
}

func (a *PostgresUserRepositoryAdapter) GetUserByID(ctx context.Context, id int) (models.User, error) {
	userDB, err := a.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return models.User{}, err
	}

	return dtoToModel(userDB), nil
}

func dtoToModel(userDB dto.UserDB) models.User {
	return models.User{
		ID:          userDB.ID,
		FirstName:   userDB.FirstName,
		LastName:    userDB.LastName,
		Email:       userDB.Email,
		Login:       userDB.Login,
		Password:    userDB.Password,
		Description: userDB.Description,
		CreatedAt:   userDB.CreatedAt,
		UpdatedAt:   userDB.UpdatedAt,
	}
}
