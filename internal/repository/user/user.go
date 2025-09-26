package user

import (
	"errors"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/dto"
)

type Repository struct {
	users []dto.UserDB
}

func NewRepository(users []dto.UserDB) *Repository {
	return &Repository{users: users}
}

func (r *Repository) CreateUser(user models.User) (models.User, error) {
	for _, u := range r.users {
		if u.Login == user.Login {
			return models.User{}, errors.New("user with this login already exists")
		}
		if u.Email == user.Email {
			return models.User{}, errors.New("user with this email already exists")
		}
	}

	newID := 1
	for _, u := range r.users {
		if u.ID >= newID {
			newID = u.ID + 1
		}
	}

	now := time.Now()
	userDB := dto.UserDB{
		ID:        newID,
		Name:      user.Name,
		Surname:   user.Surname,
		Email:     user.Email,
		Login:     user.Login,
		Password:  user.Password,
		CreatedAt: now,
		UpdatedAt: now,
	}

	r.users = append(r.users, userDB)

	return models.User{
		ID:        userDB.ID,
		Name:      userDB.Name,
		Surname:   userDB.Surname,
		Email:     userDB.Email,
		Login:     userDB.Login,
		Password:  "",
		CreatedAt: userDB.CreatedAt,
		UpdatedAt: userDB.UpdatedAt,
	}, nil
}

func (r *Repository) GetUserByLogin(login string) (models.User, error) {
	for _, u := range r.users {
		if u.Login == login {
			return models.User{
				ID:        u.ID,
				Name:      u.Name,
				Surname:   u.Surname,
				Email:     u.Email,
				Login:     u.Login,
				Password:  u.Password,
				CreatedAt: u.CreatedAt,
				UpdatedAt: u.UpdatedAt,
			}, nil
		}
	}
	return models.User{}, errors.New("user not found")
}

func (r *Repository) GetUserByID(id int) (models.User, error) {
	for _, u := range r.users {
		if u.ID == id {
			return models.User{
				ID:        u.ID,
				Name:      u.Name,
				Surname:   u.Surname,
				Email:     u.Email,
				Login:     u.Login,
				Password:  "",
				CreatedAt: u.CreatedAt,
				UpdatedAt: u.UpdatedAt,
			}, nil
		}
	}
	return models.User{}, errors.New("user not found")
}

func (r *Repository) GetAllUsers() []dto.UserDB {
	return r.users
}
