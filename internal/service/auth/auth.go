package auth

import (
	"errors"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/user"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
)

type Service struct {
	userRepo  *user.Repository
	jwtSecret string
}

func NewService(userRepo *user.Repository, jwtSecret string) *Service {
	return &Service{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *Service) Register(req models.RegisterRequest) (models.AuthResponse, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return models.AuthResponse{}, err
	}

	user := models.User{
		Name:     "",
		Surname:  "",
		Email:    req.Email,
		Login:    req.Login,
		Password: hashedPassword,
	}

	createdUser, err := s.userRepo.CreateUser(user)
	if err != nil {
		return models.AuthResponse{}, err
	}

	token, err := utils.GenerateJWT(createdUser.ID, createdUser.Login, s.jwtSecret)
	if err != nil {
		return models.AuthResponse{}, err
	}

	return models.AuthResponse{
		Token: token,
		User:  createdUser,
	}, nil
}

func (s *Service) Login(req models.LoginRequest) (models.AuthResponse, error) {
	user, err := s.userRepo.GetUserByLogin(req.Login)
	if err != nil {
		return models.AuthResponse{}, errors.New("invalid credentials")
	}

	valid, err := utils.VerifyPassword(req.Password, user.Password)
	if err != nil {
		return models.AuthResponse{}, err
	}

	if !valid {
		return models.AuthResponse{}, errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID, user.Login, s.jwtSecret)
	if err != nil {
		return models.AuthResponse{}, err
	}

	user.Password = ""

	return models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *Service) GetUserByID(userID int) (models.User, error) {
	return s.userRepo.GetUserByID(userID)
}
