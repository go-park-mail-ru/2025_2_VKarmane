package user

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresRepository_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	now := time.Now()
	user := dto.UserDB{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Login:     "johndoe",
		Password:  "hashed_password",
		CreatedAt: now,
		UpdatedAt: now,
	}

	mock.ExpectQuery(`INSERT INTO "user"`).
		WithArgs(user.FirstName, user.LastName, user.Email, user.Login, user.Password, user.Description, user.LogoHashedID, user.CreatedAt, user.UpdatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"_id"}).AddRow(5))

	id, err := repo.CreateUser(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, 5, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_CreateUser_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	now := time.Now()
	user := dto.UserDB{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Login:     "johndoe",
		Password:  "hashed_password",
		CreatedAt: now,
		UpdatedAt: now,
	}

	mock.ExpectQuery(`INSERT INTO "user"`).
		WithArgs(user.FirstName, user.LastName, user.Email, user.Login, user.Password, user.Description, user.LogoHashedID, user.CreatedAt, user.UpdatedAt).
		WillReturnError(sql.ErrConnDone)

	id, err := repo.CreateUser(context.Background(), user)
	assert.Error(t, err)
	assert.Zero(t, id)
	assert.Contains(t, err.Error(), "failed to create user")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	now := time.Now()
	email := "john@example.com"
	rows := sqlmock.NewRows([]string{"_id", "user_name", "surname", "email", "user_login", "user_hashed_password", "user_description", "logo_hashed_id", "created_at", "updated_at"}).
		AddRow(1, "John", "Doe", email, "johndoe", "hashed", "", "", now, now)

	mock.ExpectQuery(`SELECT.*FROM "user"`).
		WithArgs(email).
		WillReturnRows(rows)

	user, err := repo.GetUserByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, "johndoe", user.Login)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetUserByEmail_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	email := "nonexistent@example.com"
	mock.ExpectQuery(`SELECT.*FROM "user"`).
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetUserByEmail(context.Background(), email)
	assert.Error(t, err)
	assert.Zero(t, user.ID)
	assert.Contains(t, err.Error(), "failed to get user by email")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetUserByLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	now := time.Now()
	login := "johndoe"
	rows := sqlmock.NewRows([]string{"_id", "user_name", "surname", "email", "user_login", "user_hashed_password", "user_description", "logo_hashed_id", "created_at", "updated_at"}).
		AddRow(1, "John", "Doe", "john@example.com", login, "hashed", "", "", now, now)

	mock.ExpectQuery(`SELECT.*FROM "user"`).
		WithArgs(login).
		WillReturnRows(rows)

	user, err := repo.GetUserByLogin(context.Background(), login)
	assert.NoError(t, err)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, login, user.Login)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetUserByLogin_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	login := "nonexistent"
	mock.ExpectQuery(`SELECT.*FROM "user"`).
		WithArgs(login).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetUserByLogin(context.Background(), login)
	assert.Error(t, err)
	assert.Zero(t, user.ID)
	assert.Contains(t, err.Error(), "failed to get user by login")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	now := time.Now()
	userID := 1
	rows := sqlmock.NewRows([]string{"_id", "user_name", "surname", "email", "user_login", "user_hashed_password", "user_description", "logo_hashed_id", "created_at", "updated_at"}).
		AddRow(userID, "John", "Doe", "john@example.com", "johndoe", "hashed", "", "", now, now)

	mock.ExpectQuery(`SELECT.*FROM "user"`).
		WithArgs(userID).
		WillReturnRows(rows)

	user, err := repo.GetUserByID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetUserByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	userID := 999
	mock.ExpectQuery(`SELECT.*FROM "user"`).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetUserByID(context.Background(), userID)
	assert.Error(t, err)
	assert.Zero(t, user.ID)
	assert.Contains(t, err.Error(), "failed to get user by ID")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_UpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	now := time.Now()
	user := dto.UserDB{
		ID:        1,
		FirstName: "Updated",
		LastName:  "Name",
		Email:     "updated@example.com",
		Login:     "updated",
		Password:  "new_hashed",
		UpdatedAt: now,
	}

	mock.ExpectExec(`UPDATE "user"`).
		WithArgs(user.FirstName, user.LastName, user.Email, user.LogoHashedID, user.UpdatedAt, user.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateUser(context.Background(), user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_UpdateUser_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	now := time.Now()
	user := dto.UserDB{
		ID:        1,
		FirstName: "Updated",
		LastName:  "Name",
		Email:     "updated@example.com",
		Login:     "updated",
		Password:  "new_hashed",
		UpdatedAt: now,
	}

	mock.ExpectExec(`UPDATE "user"`).
		WithArgs(user.FirstName, user.LastName, user.Email, user.LogoHashedID, user.UpdatedAt, user.ID).
		WillReturnError(sql.ErrConnDone)

	err = repo.UpdateUser(context.Background(), user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update user")
	assert.NoError(t, mock.ExpectationsWereMet())
}

