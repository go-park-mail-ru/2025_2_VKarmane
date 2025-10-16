package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
)

var InvalidPassword = errors.New("invalid credentials")

type Service struct {
	userRepo  UserRepository
	jwtSecret string
}

func NewService(userRepo UserRepository, jwtSecret string) *Service {
	return &Service{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *Service) Register(ctx context.Context, req models.RegisterRequest) (models.AuthResponse, error) {
	log := logger.FromContext(ctx)
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		if log != nil {
			log.Error("Failed to hash password", "error", err)
		}

		return models.AuthResponse{}, fmt.Errorf("auth.Register: failed to hash password: %w", err)
	}

	user := models.User{
		FirstName: "", // Необязательное поле
		LastName:  "", // Необязательное поле
		Email:     req.Email,
		Login:     req.Login,
		Password:  hashedPassword,
	}

	createdUser, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		if log != nil {
			log.Error("Failed to create user", "error", err, "login", req.Login)
		}

		return models.AuthResponse{}, fmt.Errorf("auth.Register: failed to create user: %w", err)
	}

	token, err := utils.GenerateJWT(createdUser.ID, createdUser.Login, s.jwtSecret)
	if err != nil {
		return models.AuthResponse{}, fmt.Errorf("auth.Register: failed to generate token: %w", err)
	}

	return models.AuthResponse{
		Token: token,
		User:  createdUser,
	}, nil
}

func (s *Service) Login(ctx context.Context, req models.LoginRequest) (models.AuthResponse, error) {
	log := logger.FromContext(ctx)
	user, err := s.userRepo.GetUserByLogin(ctx, req.Login)
	if err != nil {
		if log != nil {
			log.Warn("Login attempt with invalid credentials", "login", req.Login, "error", err)
		}

		return models.AuthResponse{}, fmt.Errorf("auth.Login: invalid credentials")
	}

	valid, err := utils.VerifyPassword(req.Password, user.Password)
	if err != nil {
		if log != nil {
			log.Error("Failed to verify password", "error", err, "user_id", user.ID)
		}

		return models.AuthResponse{}, fmt.Errorf("auth.Login: failed to verify password: %w", err)
	}

	if !valid {
		if log != nil {
			log.Warn("Login attempt with invalid password", "login", req.Login, "user_id", user.ID)
		}

		return models.AuthResponse{}, InvalidPassword
	}

	token, err := utils.GenerateJWT(user.ID, user.Login, s.jwtSecret)
	if err != nil {
		if log != nil {
			log.Error("Failed to generate JWT token", "error", err, "user_id", user.ID)
		}

		return models.AuthResponse{}, fmt.Errorf("auth.Login: failed to generate token: %w", err)
	}

	user.Password = ""

	return models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *Service) GetUserByID(ctx context.Context, userID int) (models.User, error) {
	log := logger.FromContext(ctx)
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if log != nil {
			log.Error("Failed to get user by ID", "error", err, "user_id", userID)
		}

		return models.User{}, fmt.Errorf("auth.GetUserByID: %w", err)
	}

	return user, nil
}
