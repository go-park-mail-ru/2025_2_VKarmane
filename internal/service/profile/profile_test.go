package profile

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type combinedRepo struct {
	userRepo *mocks.MockUserRepository
}

func (c *combinedRepo) GetUserByID(ctx context.Context, id int) (models.User, error) {
	return c.userRepo.GetUserByID(ctx, id)
}

func (c *combinedRepo) UpdateUser(ctx context.Context, user models.User) error {
	return c.userRepo.UpdateUser(ctx, user)
}

func TestService_GetProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	repo := &combinedRepo{userRepo: mockUserRepo}
	svc := NewService(repo)

	user := models.User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Login:     "johndoe",
		CreatedAt: time.Now(),
	}

	mockUserRepo.EXPECT().GetUserByID(gomock.Any(), 1).Return(user, nil)

	result, err := svc.GetProfile(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "John", result.FirstName)
	assert.Equal(t, "Doe", result.LastName)
	assert.Equal(t, "john@example.com", result.Email)
}

func TestService_GetProfile_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	repo := &combinedRepo{userRepo: mockUserRepo}
	svc := NewService(repo)

	mockUserRepo.EXPECT().GetUserByID(gomock.Any(), 1).Return(models.User{}, errors.New("user not found"))

	result, err := svc.GetProfile(context.Background(), 1)
	assert.Error(t, err)
	assert.Empty(t, result.Email)
}

func TestService_UpdateProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	repo := &combinedRepo{userRepo: mockUserRepo}
	svc := NewService(repo)

	user := models.User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Login:     "johndoe",
	}

	req := models.UpdateProfileRequest{
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane@example.com",
	}

	mockUserRepo.EXPECT().GetUserByID(gomock.Any(), 1).Return(user, nil)
	mockUserRepo.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(nil)

	result, err := svc.UpdateProfile(context.Background(), req, 1)
	assert.NoError(t, err)
	assert.Equal(t, "Jane", result.FirstName)
	assert.Equal(t, "Smith", result.LastName)
}

