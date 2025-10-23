package auth

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
	"github.com/go-park-mail-ru/2025_2_VKarmane/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_Login_Success(t *testing.T) {
	hashed, err := utils.HashPassword("goodpass")
	if err != nil {
		t.Fatalf("hash err: %v", err)
	}

	repo := &mocks.UserRepository{}
	svc := NewService(repo, "secret")

	user := models.User{ID: 10, Login: "john", Password: hashed}
	req := models.LoginRequest{Login: "john", Password: "goodpass"}

	repo.On("GetUserByLogin", mock.Anything, "john").Return(user, nil)

	res, err := svc.Login(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEmpty(t, res.Token)
	assert.Equal(t, 10, res.User.ID)
	assert.Equal(t, "john", res.User.Login)
	assert.Empty(t, res.User.Password)
}

func TestService_Login_InvalidPassword(t *testing.T) {
	hashed, err := utils.HashPassword("goodpass")
	if err != nil {
		t.Fatalf("hash err: %v", err)
	}

	repo := &mocks.UserRepository{}
	svc := NewService(repo, "secret")

	user := models.User{ID: 11, Login: "jane", Password: hashed}
	req := models.LoginRequest{Login: "jane", Password: "badpass"}

	repo.On("GetUserByLogin", mock.Anything, "jane").Return(user, nil)

	_, err = svc.Login(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestService_Login_VerifyError(t *testing.T) {
	repo := &mocks.UserRepository{}
	svc := NewService(repo, "secret")

	user := models.User{ID: 12, Login: "mike", Password: "bad-hash"}
	req := models.LoginRequest{Login: "mike", Password: "any"}

	repo.On("GetUserByLogin", mock.Anything, "mike").Return(user, nil)

	_, err := svc.Login(context.Background(), req)
	assert.Error(t, err)
}
