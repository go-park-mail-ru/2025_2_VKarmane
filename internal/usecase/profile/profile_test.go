package profile

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUseCase_GetProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockProfileService(ctrl)
	uc := NewUseCase(mockSvc)

	expectedProfile := models.ProfileResponse{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}

	mockSvc.EXPECT().GetProfile(gomock.Any(), 1).Return(expectedProfile, nil)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	result, err := uc.GetProfile(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedProfile, result)
}

func TestUseCase_GetProfile_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockProfileService(ctrl)
	uc := NewUseCase(mockSvc)

	mockSvc.EXPECT().GetProfile(gomock.Any(), 1).Return(models.ProfileResponse{}, errors.New("db error"))

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	result, err := uc.GetProfile(ctx, 1)
	assert.Error(t, err)
	assert.Empty(t, result.Email)
}

func TestUseCase_UpdateProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockProfileService(ctrl)
	uc := NewUseCase(mockSvc)

	req := models.UpdateProfileRequest{
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane@example.com",
	}

	expectedProfile := models.ProfileResponse{
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane@example.com",
	}

	mockSvc.EXPECT().UpdateProfile(gomock.Any(), req, 1).Return(expectedProfile, nil)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	result, err := uc.UpdateProfile(ctx, req, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedProfile, result)
}
