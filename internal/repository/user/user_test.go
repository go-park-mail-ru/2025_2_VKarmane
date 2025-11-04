package user

import (
	"context"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/dto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_GetUserByLogin_Found(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	repo := NewRepository([]dto.UserDB{
		{ID: 1, FirstName: "A", LastName: "B", Email: "a@b.c", Login: "ab", Password: "hash"},
	}, fixedClock)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	u, err := repo.GetUserByLogin(ctx, "ab")
	assert.NoError(t, err)
	assert.Equal(t, "ab", u.Login)
	assert.Equal(t, 1, u.ID)
	assert.Equal(t, "A", u.FirstName)
	assert.Equal(t, "B", u.LastName)
	assert.Equal(t, "a@b.c", u.Email)
	assert.Equal(t, "hash", u.Password)
}

func TestUserRepository_GetUserByLogin_NotFound(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	repo := NewRepository([]dto.UserDB{
		{ID: 1, FirstName: "A", LastName: "B", Email: "a@b.c", Login: "ab", Password: "hash"},
	}, fixedClock)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	_, err := repo.GetUserByLogin(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
}

func TestUserRepository_GetUserByLogin_Empty(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	repo := NewRepository([]dto.UserDB{}, fixedClock)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	_, err := repo.GetUserByLogin(ctx, "any")
	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
}

func TestUserRepository_GetUserByID_Found(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	repo := NewRepository([]dto.UserDB{
		{ID: 1, FirstName: "A", LastName: "B", Email: "a@b.c", Login: "ab", Password: "hash"},
	}, fixedClock)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	u, err := repo.GetUserByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, u.ID)
	assert.Equal(t, "ab", u.Login)
	assert.Equal(t, "A", u.FirstName)
	assert.Equal(t, "B", u.LastName)
	assert.Equal(t, "a@b.c", u.Email)
	assert.Empty(t, u.Password)
}

func TestUserRepository_GetUserByID_NotFound(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	repo := NewRepository([]dto.UserDB{
		{ID: 1, FirstName: "A", LastName: "B", Email: "a@b.c", Login: "ab", Password: "hash"},
	}, fixedClock)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	_, err := repo.GetUserByID(ctx, 99)
	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
}

func TestUserRepository_GetUserByID_Empty(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	repo := NewRepository([]dto.UserDB{}, fixedClock)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	_, err := repo.GetUserByID(ctx, 1)
	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
}

func TestUserRepository_CreateUser_DuplicateLogin(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	repo := NewRepository([]dto.UserDB{
		{ID: 1, FirstName: "A", LastName: "B", Email: "a@b.c", Login: "ab", Password: "hash"},
	}, fixedClock)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	_, err := repo.CreateUser(ctx, models.User{Login: "ab", Email: "new@x.y"})
	assert.Error(t, err)
}

func TestUserRepository_CreateUser_DuplicateEmail(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	repo := NewRepository([]dto.UserDB{
		{ID: 1, FirstName: "A", LastName: "B", Email: "a@b.c", Login: "ab", Password: "hash"},
	}, fixedClock)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	_, err := repo.CreateUser(ctx, models.User{Login: "newlogin", Email: "a@b.c"})
	assert.Error(t, err)
}

func TestUserRepository_CreateUser_Success(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	repo := NewRepository([]dto.UserDB{
		{ID: 1, FirstName: "A", LastName: "B", Email: "a@b.c", Login: "ab", Password: "hash"},
	}, fixedClock)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	created, err := repo.CreateUser(ctx, models.User{FirstName: "N", LastName: "L", Email: "n@l.m", Login: "nl", Password: "p"})
	assert.NoError(t, err)
	assert.NotZero(t, created.ID)
	assert.Empty(t, created.Password)
}

func TestUserRepository_EditUserByID_Success(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	repo := NewRepository([]dto.UserDB{
		{ID: 1, FirstName: "Old", LastName: "Name", Email: "old@mail.com", Login: "user1"},
	}, fixedClock)

	req := models.UpdateProfileRequest{
		FirstName: "New",
		LastName:  "Name",
		Email:     "new@mail.com",
	}

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	updated, err := repo.EditUserByID(ctx, req, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, updated.ID)
	assert.Equal(t, "New", updated.FirstName)
	assert.Equal(t, "Name", updated.LastName)
	assert.Equal(t, "new@mail.com", updated.Email)
	assert.Empty(t, updated.Password)
}

func TestUserRepository_EditUserByID_EmailExists(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	repo := NewRepository([]dto.UserDB{
		{ID: 1, FirstName: "A", LastName: "B", Email: "a@b.c", Login: "ab"},
		{ID: 2, FirstName: "C", LastName: "D", Email: "x@y.z", Login: "cd"},
	}, fixedClock)

	req := models.UpdateProfileRequest{
		FirstName: "C",
		LastName:  "D",
		Email:     "a@b.c",
	}

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	_, err := repo.EditUserByID(ctx, req, 2)
	assert.Error(t, err)
	assert.Equal(t, ErrEmailExists, err)
}

func TestUserRepository_EditUserByID_NotFound(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	repo := NewRepository([]dto.UserDB{
		{ID: 1, FirstName: "A", LastName: "B", Email: "a@b.c", Login: "ab"},
	}, fixedClock)

	req := models.UpdateProfileRequest{
		FirstName: "Z",
		LastName:  "Y",
		Email:     "z@y.x",
	}

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	_, err := repo.EditUserByID(ctx, req, 99)
	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
}

func TestUserRepository_CreateUser_Success_EmptyRepo(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	repo := NewRepository([]dto.UserDB{}, fixedClock)

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	created, err := repo.CreateUser(ctx, models.User{
		FirstName: "New",
		LastName:  "User",
		Email:     "new@user.com",
		Login:     "newuser",
		Password:  "password",
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, created.ID)
	assert.Equal(t, "newuser", created.Login)
	assert.Equal(t, "New", created.FirstName)
	assert.Equal(t, "User", created.LastName)
	assert.Equal(t, "new@user.com", created.Email)
	assert.Empty(t, created.Password)
}

func TestUserRepository_GetAllUsers(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	users := []dto.UserDB{
		{ID: 1, FirstName: "A", LastName: "B", Email: "a@b.c", Login: "ab"},
		{ID: 2, FirstName: "C", LastName: "D", Email: "c@d.e", Login: "cd"},
	}
	repo := NewRepository(users, fixedClock)

	allUsers := repo.GetAllUsers()
	assert.Len(t, allUsers, 2)
	assert.Equal(t, 1, allUsers[0].ID)
	assert.Equal(t, 2, allUsers[1].ID)
}

func TestUserRepository_EditUserByID_SameEmail(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}
	repo := NewRepository([]dto.UserDB{
		{ID: 1, FirstName: "A", LastName: "B", Email: "a@b.c", Login: "ab"},
	}, fixedClock)

	req := models.UpdateProfileRequest{
		FirstName: "Updated",
		LastName:  "Name",
		Email:     "a@b.c",
	}

	ctx := logger.WithLogger(context.Background(), logger.NewSlogLogger())
	updated, err := repo.EditUserByID(ctx, req, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, updated.ID)
	assert.Equal(t, "Updated", updated.FirstName)
	assert.Equal(t, "Name", updated.LastName)
	assert.Equal(t, "a@b.c", updated.Email)
}
