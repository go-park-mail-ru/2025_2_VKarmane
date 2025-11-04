package auth

import (
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestUserToApi(t *testing.T) {
	now := time.Now()
	user := models.User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Login:     "johndoe",
		Email:     "john@example.com",
		CreatedAt: now,
	}

	apiUser := UserToApi(user)

	assert.Equal(t, 1, apiUser.ID)
	assert.Equal(t, "John", apiUser.FirstName)
	assert.Equal(t, "Doe", apiUser.LastName)
	assert.Equal(t, "johndoe", apiUser.Login)
	assert.Equal(t, "john@example.com", apiUser.Email)
	assert.Equal(t, now, apiUser.CreatedAt)
}
