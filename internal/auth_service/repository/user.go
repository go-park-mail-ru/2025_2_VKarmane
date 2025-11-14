package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	authmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/models"
)

var ErrLoginExists = errors.New("login exists")
var ErrEmailExists = errors.New("email exists")
var ErrUserNotFound = errors.New("not Found")

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, user authmodels.User) (authmodels.User, error) {
	var id int

	if user.CreatedAt.IsZero() || user.UpdatedAt.IsZero() {
		query := `
			INSERT INTO "user" (user_name, surname, email, user_login, user_hashed_password, user_description, logo_hashed_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
			RETURNING _id, created_at, updated_at
		`

		err := r.db.QueryRowContext(ctx, query,
			user.FirstName,
			user.LastName,
			user.Email,
			user.Login,
			user.Password,
			user.Description,
			user.LogoHashedID,
		).Scan(&id, &user.CreatedAt, &user.UpdatedAt)

		if err != nil {
			return authmodels.User{}, fmt.Errorf("failed to create user: %w", err)
		}
	} else {
		query := `
			INSERT INTO "user" (user_name, surname, email, user_login, user_hashed_password, user_description, logo_hashed_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING _id
		`

		err := r.db.QueryRowContext(ctx, query,
			user.FirstName,
			user.LastName,
			user.Email,
			user.Login,
			user.Password,
			user.Description,
			user.LogoHashedID,
			user.CreatedAt,
			user.UpdatedAt,
		).Scan(&id)

		if err != nil {
			return authmodels.User{}, fmt.Errorf("failed to create user: %w", err)
		}
	}

	user.ID = id
	return user, nil
}


// func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (dto.UserDB, error) {
// 	query := `
// 		SELECT _id, user_name, surname, email, user_login, user_hashed_password, user_description, logo_hashed_id, created_at, updated_at
// 		FROM "user"
// 		WHERE email = $1
// 	`

// 	var user dto.UserDB
// 	err := r.db.QueryRowContext(ctx, query, email).Scan(
// 		&user.ID,
// 		&user.FirstName,
// 		&user.LastName,
// 		&user.Email,
// 		&user.Login,
// 		&user.Password,
// 		&user.Description,
// 		&user.LogoHashedID,
// 		&user.CreatedAt,
// 		&user.UpdatedAt,
// 	)

// 	if err != nil {
// 		return dto.UserDB{}, fmt.Errorf("failed to get user by email: %w", err)
// 	}

// 	return user, nil
// }

func (r *PostgresRepository) GetUserByLogin(ctx context.Context, login string) (authmodels.User, error) {
	query := `
		SELECT _id, user_name, surname, email, user_login, user_hashed_password, user_description, logo_hashed_id, created_at, updated_at
		FROM "user"
		WHERE user_login = $1
	`

	var user authmodels.User
	err := r.db.QueryRowContext(ctx, query, login).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Login,
		&user.Password,
		&user.Description,
		&user.LogoHashedID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return authmodels.User{}, fmt.Errorf("failed to get user by login: %w", err)
	}

	return user, nil
}

func (r *PostgresRepository) GetUserByID(ctx context.Context, id int) (authmodels.User, error) {
	query := `
		SELECT _id, user_name, surname, email, user_login, user_hashed_password, user_description, logo_hashed_id, created_at, updated_at
		FROM "user"
		WHERE _id = $1
	`

	var user authmodels.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Login,
		&user.Password,
		&user.Description,
		&user.LogoHashedID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return authmodels.User{}, ErrUserNotFound
		}
		return authmodels.User{}, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}


func (r *PostgresRepository) EditUserByID(ctx context.Context, user authmodels.UpdateProfileRequest) (authmodels.User, error) {
	query := `
		UPDATE "user"
		SET user_name = $1,
		    surname = $2,
		    email = $3,
		    logo_hashed_id = $4,
		    updated_at = NOW()
		WHERE _id = $5
		RETURNING _id, user_name, surname, email, user_login, user_hashed_password,
		          user_description, logo_hashed_id, created_at, updated_at
	`

	var updatedUser authmodels.User

	err := r.db.QueryRowContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.LogoHashedID,
		user.UserID, 
	).Scan(
		&updatedUser.ID,
		&updatedUser.FirstName,
		&updatedUser.LastName,
		&updatedUser.Email,
		&updatedUser.Login,
		&updatedUser.Password,
		&updatedUser.Description,
		&updatedUser.LogoHashedID,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)

	if err != nil {
		return authmodels.User{}, fmt.Errorf("failed to update user: %w", err)
	}

	return updatedUser, nil
}
