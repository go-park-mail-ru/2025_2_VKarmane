package user

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"

	serviceerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/errors"
	authmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/models"
)

func TestPostgresRepository_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	user := authmodels.User{
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "john@example.com",
		Login:        "johndoe",
		Password:     "hashedpassword",
		Description:  "desc",
		LogoHashedID: "logo123",
	}

	// мок для вставки с NOW()
	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO "user" (user_name, surname, email, user_login, user_hashed_password, user_description, logo_hashed_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING _id, created_at, updated_at
	`)).
		WithArgs(user.FirstName, user.LastName, user.Email, user.Login, user.Password, user.Description, user.LogoHashedID).
		WillReturnRows(sqlmock.NewRows([]string{"_id", "created_at", "updated_at"}).AddRow(1, time.Now(), time.Now()))

	createdUser, err := repo.CreateUser(context.Background(), user)
	require.NoError(t, err)
	require.Equal(t, 1, createdUser.ID)
}

func TestPostgresRepository_GetUserByLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	user := authmodels.User{
		ID:       1,
		FirstName: "John",
		LastName:  "Doe",
		Email:    "john@example.com",
		Login:    "johndoe",
		Password: "hashedpassword",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT _id, user_name, surname, email, user_login, user_hashed_password, user_description, logo_hashed_id, created_at, updated_at
		FROM "user"
		WHERE user_login = $1
	`)).
		WithArgs("johndoe").
		WillReturnRows(sqlmock.NewRows([]string{"_id", "user_name", "surname", "email", "user_login", "user_hashed_password", "user_description", "logo_hashed_id", "created_at", "updated_at"}).
			AddRow(user.ID, user.FirstName, user.LastName, user.Email, user.Login, user.Password, "", "", time.Now(), time.Now()))

	got, err := repo.GetUserByLogin(context.Background(), "johndoe")
	require.NoError(t, err)
	require.Equal(t, user.ID, got.ID)
	require.Equal(t, user.Login, got.Login)
}

func TestPostgresRepository_GetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	user := authmodels.User{
		ID:       1,
		FirstName: "John",
		LastName:  "Doe",
		Email:    "john@example.com",
		Login:    "johndoe",
		Password: "hashedpassword",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT _id, user_name, surname, email, user_login, user_hashed_password, user_description, logo_hashed_id, created_at, updated_at
		FROM "user"
		WHERE _id = $1
	`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"_id", "user_name", "surname", "email", "user_login", "user_hashed_password", "user_description", "logo_hashed_id", "created_at", "updated_at"}).
			AddRow(user.ID, user.FirstName, user.LastName, user.Email, user.Login, user.Password, "", "", time.Now(), time.Now()))

	got, err := repo.GetUserByID(context.Background(), 1)
	require.NoError(t, err)
	require.Equal(t, user.ID, got.ID)
}

func TestPostgresRepository_EditUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	req := authmodels.UpdateProfileRequest{
		UserID:    1,
		FirstName: "JohnUpdated",
		LastName:  "DoeUpdated",
		Email:     "johnupdated@example.com",
		LogoHashedID: "logo456",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		UPDATE "user"
		SET user_name = $1,
		    surname = $2,
		    email = $3,
		    logo_hashed_id = $4,
		    updated_at = NOW()
		WHERE _id = $5
		RETURNING _id, user_name, surname, email, user_login, user_hashed_password,
		          user_description, logo_hashed_id, created_at, updated_at
	`)).
		WithArgs(req.FirstName, req.LastName, req.Email, req.LogoHashedID, req.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"_id", "user_name", "surname", "email", "user_login", "user_hashed_password", "user_description", "logo_hashed_id", "created_at", "updated_at"}).
			AddRow(req.UserID, req.FirstName, req.LastName, req.Email, "johndoe", "hashedpassword", "", req.LogoHashedID, time.Now(), time.Now()))

	updated, err := repo.EditUserByID(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, req.FirstName, updated.FirstName)
	require.Equal(t, req.LastName, updated.LastName)
	require.Equal(t, req.Email, updated.Email)
	require.Equal(t, req.LogoHashedID, updated.LogoHashedID)
}

func TestPostgresRepository_CreateUser_UniqueViolation(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	user := authmodels.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Login:     "johndoe",
		Password:  "hashedpassword",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO "user" (user_name, surname, email, user_login, user_hashed_password, user_description, logo_hashed_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING _id, created_at, updated_at
	`)).
		WithArgs(user.FirstName, user.LastName, user.Email, user.Login, user.Password, "", "").
		WillReturnError(&pq.Error{Code: "23505", Constraint: "user_login_key"})

	_, err = repo.CreateUser(context.Background(), user)
	require.ErrorIs(t, err, serviceerrors.ErrLoginExists)
}

func TestPostgresRepository_GetUserByLogin_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT _id, user_name, surname, email, user_login, user_hashed_password, user_description, logo_hashed_id, created_at, updated_at
		FROM "user"
		WHERE user_login = $1
	`)).
		WithArgs("unknown").
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetUserByLogin(context.Background(), "unknown")
	require.ErrorIs(t, err, serviceerrors.ErrUserNotFound)
}

func TestPostgresRepository_GetUserByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT _id, user_name, surname, email, user_login, user_hashed_password, user_description, logo_hashed_id, created_at, updated_at
		FROM "user"
		WHERE _id = $1
	`)).
		WithArgs(99).
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetUserByID(context.Background(), 99)
	require.ErrorIs(t, err, serviceerrors.ErrUserNotFound)
}

func TestPostgresRepository_EditUserByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	req := authmodels.UpdateProfileRequest{
		UserID: 99,
		FirstName: "Test",
		LastName: "User",
		Email: "test@example.com",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		UPDATE "user"
		SET user_name = $1,
		    surname = $2,
		    email = $3,
		    logo_hashed_id = $4,
		    updated_at = NOW()
		WHERE _id = $5
		RETURNING _id, user_name, surname, email, user_login, user_hashed_password,
		          user_description, logo_hashed_id, created_at, updated_at
	`)).
		WithArgs(req.FirstName, req.LastName, req.Email, req.LogoHashedID, req.UserID).
		WillReturnError(sql.ErrNoRows)

	_, err = repo.EditUserByID(context.Background(), req)
	require.ErrorIs(t, err, serviceerrors.ErrUserNotFound)
}

func TestPostgresRepository_EditUserByID_EmailExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	req := authmodels.UpdateProfileRequest{
		UserID: 1,
		FirstName: "John",
		LastName: "Doe",
		Email: "existing@example.com",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		UPDATE "user"
		SET user_name = $1,
		    surname = $2,
		    email = $3,
		    logo_hashed_id = $4,
		    updated_at = NOW()
		WHERE _id = $5
		RETURNING _id, user_name, surname, email, user_login, user_hashed_password,
		          user_description, logo_hashed_id, created_at, updated_at
	`)).
		WithArgs(req.FirstName, req.LastName, req.Email, req.LogoHashedID, req.UserID).
		WillReturnError(&pq.Error{Code: "23505"})

	_, err = repo.EditUserByID(context.Background(), req)
	require.ErrorIs(t, err, serviceerrors.ErrEmailExists)
}
