package auth

import (
	"fmt"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
)

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

func (s *Service) Register(req models.RegisterRequest) (models.AuthResponse, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return models.AuthResponse{}, fmt.Errorf("auth.Register: failed to hash password: %w", err)
	}

	user := models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Login:     req.Login,
		Password:  hashedPassword,
	}

	createdUser, err := s.userRepo.CreateUser(user)
	if err != nil {
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

func (s *Service) Login(req models.LoginRequest) (models.AuthResponse, error) {
	user, err := s.userRepo.GetUserByLogin(req.Login)
	if err != nil {
		return models.AuthResponse{}, fmt.Errorf("auth.Login: invalid credentials")
	}

	valid, err := utils.VerifyPassword(req.Password, user.Password)
	if err != nil {
		return models.AuthResponse{}, fmt.Errorf("auth.Login: failed to verify password: %w", err)
	}

	if !valid {
		return models.AuthResponse{}, fmt.Errorf("auth.Login: invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID, user.Login, s.jwtSecret)
	if err != nil {
		return models.AuthResponse{}, fmt.Errorf("auth.Login: failed to generate token: %w", err)
	}

	user.Password = ""

	return models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *Service) GetUserByID(userID int) (models.User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return models.User{}, fmt.Errorf("auth.GetUserByID: %w", err)
	}
	return user, nil
}
