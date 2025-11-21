package auth

import (
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestLoginApiToProtoLogin(t *testing.T) {
	req := models.LoginRequest{
		Login:    "testuser",
		Password: "secret123",
	}

	protoReq := LoginApiToProtoLogin(req)

	assert.Equal(t, "testuser", protoReq.Login)
	assert.Equal(t, "secret123", protoReq.Password)
}

func TestRegisterApiToProtoRegister(t *testing.T) {
	req := models.RegisterRequest{
		Login:    "newuser",
		Email:    "new@example.com",
		Password: "pass1234",
	}

	protoReq := RegisterApiToProtoRegister(req)

	assert.Equal(t, "newuser", protoReq.Login)
	assert.Equal(t, "new@example.com", protoReq.Email)
	assert.Equal(t, "pass1234", protoReq.Password)
}
