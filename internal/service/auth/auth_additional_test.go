package auth

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/mocks"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestService_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	realClock := clock.RealClock{}
	hashed, err := utils.HashPassword("goodpass")
	if err != nil {
		t.Fatalf("hash err: %v", err)
	}

	repo := mocks.NewMockUserRepository(ctrl)
	svc := NewService(repo, "secret", realClock)

	user := models.User{ID: 10, Login: "john", Password: hashed}
	req := models.LoginRequest{Login: "john", Password: "goodpass"}

	repo.EXPECT().GetUserByLogin(gomock.Any(), "john").Return(user, nil)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	res, err := svc.Login(ctx, req)
	assert.NoError(t, err)
	assert.NotEmpty(t, res.Token)
	assert.Equal(t, 10, res.User.ID)
	assert.Equal(t, "john", res.User.Login)
	assert.Empty(t, res.User.Password)
}

func TestService_Login_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	realClock := clock.RealClock{}
	hashed, err := utils.HashPassword("goodpass")
	if err != nil {
		t.Fatalf("hash err: %v", err)
	}

	repo := mocks.NewMockUserRepository(ctrl)
	svc := NewService(repo, "secret", realClock)

	user := models.User{ID: 11, Login: "jane", Password: hashed}
	req := models.LoginRequest{Login: "jane", Password: "badpass"}

	repo.EXPECT().GetUserByLogin(gomock.Any(), "jane").Return(user, nil)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	_, err = svc.Login(ctx, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestService_Login_VerifyError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	realClock := clock.RealClock{}
	repo := mocks.NewMockUserRepository(ctrl)
	svc := NewService(repo, "secret", realClock)

	user := models.User{ID: 12, Login: "mike", Password: "bad-hash"}
	req := models.LoginRequest{Login: "mike", Password: "any"}

	repo.EXPECT().GetUserByLogin(gomock.Any(), "mike").Return(user, nil)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	_, err := svc.Login(ctx, req)
	assert.Error(t, err)
}
