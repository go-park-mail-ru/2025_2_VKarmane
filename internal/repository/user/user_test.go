package user

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/dto"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_GetUserByLogin_Found(t *testing.T) {
	repo := NewRepository([]dto.UserDB{
		{ID: 1, FirstName: "A", LastName: "B", Email: "a@b.c", Login: "ab", Password: "hash"},
	})

	u, err := repo.GetUserByLogin(context.Background(), "ab")
	assert.NoError(t, err)
	assert.Equal(t, "ab", u.Login)
}

func TestUserRepository_GetUserByID_NotFound(t *testing.T) {
	repo := NewRepository([]dto.UserDB{
		{ID: 1, FirstName: "A", LastName: "B", Email: "a@b.c", Login: "ab", Password: "hash"},
	})

	_, err := repo.GetUserByID(context.Background(), 99)
	assert.Error(t, err)
}

func TestUserRepository_CreateUser_DuplicateLogin(t *testing.T) {
	repo := NewRepository([]dto.UserDB{
		{ID: 1, FirstName: "A", LastName: "B", Email: "a@b.c", Login: "ab", Password: "hash"},
	})

	_, err := repo.CreateUser(context.Background(), models.User{Login: "ab", Email: "new@x.y"})
	assert.Error(t, err)
}

func TestUserRepository_CreateUser_DuplicateEmail(t *testing.T) {
	repo := NewRepository([]dto.UserDB{
		{ID: 1, FirstName: "A", LastName: "B", Email: "a@b.c", Login: "ab", Password: "hash"},
	})

	_, err := repo.CreateUser(context.Background(), models.User{Login: "newlogin", Email: "a@b.c"})
	assert.Error(t, err)
}

func TestUserRepository_CreateUser_Success(t *testing.T) {
	repo := NewRepository([]dto.UserDB{
		{ID: 1, FirstName: "A", LastName: "B", Email: "a@b.c", Login: "ab", Password: "hash"},
	})

	created, err := repo.CreateUser(context.Background(), models.User{FirstName: "N", LastName: "L", Email: "n@l.m", Login: "nl", Password: "p"})
	assert.NoError(t, err)
	assert.NotZero(t, created.ID)
	assert.Empty(t, created.Password)
}

func TestUserRepository_EditUserByID_Success(t *testing.T) {
	repo := NewRepository([]dto.UserDB{
		{ID: 1, FirstName: "Old", LastName: "Name", Email: "old@mail.com", Login: "user1"},
	})

	req := models.UpdateUserRequest{
		FirstName: "New",
		LastName:  "Name",
		Email:     "new@mail.com",
	}

	updated, err := repo.EditUserByID(context.Background(), req, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, updated.ID)
	assert.Equal(t, "New", updated.FirstName)
	assert.Equal(t, "Name", updated.LastName)
	assert.Equal(t, "new@mail.com", updated.Email)
	assert.Empty(t, updated.Password)
}

func TestUserRepository_EditUserByID_EmailExists(t *testing.T) {
	repo := NewRepository([]dto.UserDB{
		{ID: 1, FirstName: "A", LastName: "B", Email: "a@b.c", Login: "ab"},
		{ID: 2, FirstName: "C", LastName: "D", Email: "x@y.z", Login: "cd"},
	})

	req := models.UpdateUserRequest{
		FirstName: "C",
		LastName:  "D",
		Email:     "a@b.c",
	}

	_, err := repo.EditUserByID(context.Background(), req, 2)
	assert.Error(t, err)
	assert.Equal(t, EmailExistsErr, err)
}

func TestUserRepository_EditUserByID_NotFound(t *testing.T) {
	repo := NewRepository([]dto.UserDB{
		{ID: 1, FirstName: "A", LastName: "B", Email: "a@b.c", Login: "ab"},
	})

	req := models.UpdateUserRequest{
		FirstName: "Z",
		LastName:  "Y",
		Email:     "z@y.x",
	}

	_, err := repo.EditUserByID(context.Background(), req, 99)
	assert.Error(t, err)
	assert.Equal(t, UserNotFound, err)
}
