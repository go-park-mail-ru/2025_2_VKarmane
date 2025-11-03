package user

import (
	"context"
	"errors"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/dto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

var ErrLoginExists = errors.New("login exists")
var ErrEmailExists = errors.New("email exists")
var ErrUserNotFound = errors.New("not Found")

type Repository struct {
	users []dto.UserDB
	clock clock.Clock
}

func NewRepository(users []dto.UserDB, clck clock.Clock) *Repository {
	return &Repository{users: users, clock: clck}
}

func (r *Repository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	log := logger.FromContext(ctx)
	for _, u := range r.users {
		if u.Login == user.Login {
			log.Warn("User creation failed: login already exists", "login", user.Login)
			return models.User{}, ErrLoginExists
		}
		if u.Email == user.Email {
			log.Warn("User creation failed: email already exists", "email", user.Email)
			return models.User{}, ErrEmailExists
		}
	}

	newID := 1
	for _, u := range r.users {
		if u.ID >= newID {
			newID = u.ID + 1
		}
	}

	now := r.clock.Now()
	userDB := dto.UserDB{
		ID:        newID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Login:     user.Login,
		Password:  user.Password,
		CreatedAt: now,
		UpdatedAt: now,
	}

	r.users = append(r.users, userDB)

	return models.User{
		ID:        userDB.ID,
		FirstName: userDB.FirstName,
		LastName:  userDB.LastName,
		Email:     userDB.Email,
		Login:     userDB.Login,
		Password:  "",
		CreatedAt: userDB.CreatedAt,
		UpdatedAt: userDB.UpdatedAt,
	}, nil
}

func (r *Repository) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	log := logger.FromContext(ctx)
	for _, u := range r.users {
		if u.Login == login {
			return models.User{
				ID:        u.ID,
				FirstName: u.FirstName,
				LastName:  u.LastName,
				Email:     u.Email,
				Login:     u.Login,
				Password:  u.Password,
				CreatedAt: u.CreatedAt,
				UpdatedAt: u.UpdatedAt,
			}, nil
		}
	}

	log.Warn("User not found by login", "login", login)
	return models.User{}, ErrUserNotFound
}

func (r *Repository) GetUserByID(ctx context.Context, id int) (models.User, error) {
	log := logger.FromContext(ctx)
	for _, u := range r.users {
		if u.ID == id {
			return models.User{
				ID:        u.ID,
				FirstName: u.FirstName,
				LastName:  u.LastName,
				Email:     u.Email,
				Login:     u.Login,
				Password:  "",
				CreatedAt: u.CreatedAt,
				UpdatedAt: u.UpdatedAt,
			}, nil
		}
	}
	log.Warn("User not found by ID", "user_id", id)
	return models.User{}, ErrUserNotFound
}

func (r *Repository) GetAllUsers() []dto.UserDB {
	return r.users
}

func (r *Repository) EditUserByID(ctx context.Context, req models.UpdateProfileRequest, id int) (models.User, error) {
	log := logger.FromContext(ctx)
	now := time.Now()
	for i := range r.users {
		if r.users[i].ID != id && r.users[i].Email == req.Email {
			log.Warn("User update failed: email already exists", "email", req.Email)
			return models.User{}, ErrEmailExists
		}
		if r.users[i].ID == id {
			r.users[i].Email = req.Email
			r.users[i].FirstName = req.FirstName
			r.users[i].LastName = req.LastName

			return models.User{
				ID:        r.users[i].ID,
				FirstName: r.users[i].FirstName,
				LastName:  r.users[i].LastName,
				Email:     r.users[i].Email,
				Login:     r.users[i].Login,
				CreatedAt: r.users[i].CreatedAt,
				UpdatedAt: now,
			}, nil
		}
	}

	log.Warn("User not found by ID", "user_id", id)
	return models.User{}, ErrUserNotFound
}
