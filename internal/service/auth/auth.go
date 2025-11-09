package auth

import (
	"context"
	"errors"

	pkgerrors "github.com/pkg/errors"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

var ErrInvalidPassword = errors.New("invalid credentials")

type Service struct {
	repo interface {
		CreateUser(ctx context.Context, user models.User) (models.User, error)
		GetUserByLogin(ctx context.Context, login string) (models.User, error)
		GetUserByID(ctx context.Context, id int) (models.User, error)
		EditUserByID(ctx context.Context, req models.UpdateProfileRequest, userID int) (models.User, error)
	}
	jwtSecret string
	clock     clock.Clock
}

func NewService(repo interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
	GetUserByID(ctx context.Context, id int) (models.User, error)
	EditUserByID(ctx context.Context, req models.UpdateProfileRequest, userID int) (models.User, error)
}, jwtSecret string, clck clock.Clock) *Service {
	return &Service{
		repo:      repo,
		jwtSecret: jwtSecret,
		clock:     clck,
	}
}

func (s *Service) Register(ctx context.Context, req models.RegisterRequest) (models.AuthResponse, error) {
	log := logger.FromContext(ctx)
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Error("Failed to hash password", "error", err)
		return models.AuthResponse{}, pkgerrors.Wrap(err, "auth.Register: failed to hash password")
	}

	user := models.User{
		FirstName: "User",
		LastName:  "User",
		Email:     req.Email,
		Login:     req.Login,
		Password:  hashedPassword,
	}

	createdUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		log.Error("Failed to create user", "error", err, "login", req.Login)
		return models.AuthResponse{}, pkgerrors.Wrap(err, "auth.Register: failed to create user")
	}

	token, err := utils.GenerateJWT(createdUser.ID, createdUser.Login, s.jwtSecret)
	if err != nil {
		return models.AuthResponse{}, pkgerrors.Wrap(err, "auth.Register: failed to generate token")
	}

	return models.AuthResponse{
		Token: token,
		User:  createdUser,
	}, nil
}

func (s *Service) Login(ctx context.Context, req models.LoginRequest) (models.AuthResponse, error) {
	log := logger.FromContext(ctx)
	user, err := s.repo.GetUserByLogin(ctx, req.Login)
	if err != nil {
		log.Warn("Login attempt with invalid credentials", "login", req.Login, "error", err)
		return models.AuthResponse{}, pkgerrors.Wrap(err, "auth.Login: invalid credentials")
	}

	valid, err := utils.VerifyPassword(req.Password, user.Password)
	if err != nil {
		log.Error("Failed to verify password", "error", err, "user_id", user.ID)
		return models.AuthResponse{}, pkgerrors.Wrap(err, "auth.Login: failed to verify password")
	}

	if !valid {
		log.Warn("Login attempt with invalid password", "login", req.Login, "user_id", user.ID)
		return models.AuthResponse{}, ErrInvalidPassword
	}

	token, err := utils.GenerateJWT(user.ID, user.Login, s.jwtSecret)
	if err != nil {
		log.Error("Failed to generate JWT token", "error", err, "user_id", user.ID)
		return models.AuthResponse{}, pkgerrors.Wrap(err, "auth.Login: failed to generate token")
	}

	user.Password = ""

	return models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *Service) GetUserByID(ctx context.Context, userID int) (models.User, error) {
	log := logger.FromContext(ctx)
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		log.Error("Failed to get user by ID", "error", err, "user_id", userID)
		return models.User{}, pkgerrors.Wrap(err, "auth.GetUserByID")
	}

	return user, nil
}

func (s *Service) EditUserByID(ctx context.Context, req models.UpdateProfileRequest, userID int) (models.User, error) {
	log := logger.FromContext(ctx)
	user, err := s.repo.EditUserByID(ctx, req, userID)
	if err != nil {
		log.Error("Failed to get update user by ID", "error", err, "user_id", userID)
		return models.User{}, pkgerrors.Wrap(err, "auth.EditUserByID")
	}

	return user, nil
}
