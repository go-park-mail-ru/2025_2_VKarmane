package user

import (
	"context"
	"database/sql"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/dto"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, user dto.UserDB) (int, error) {
	query := `
		INSERT INTO "user" (user_name, surname, email, user_login, user_hashed_password, user_description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING _id
	`

	var id int
	err := r.db.QueryRowContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Login,
		user.Password,
		user.Description,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&id)

	return id, err
}

func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (dto.UserDB, error) {
	query := `
		SELECT _id, user_name, surname, email, user_login, user_hashed_password, user_description, created_at, updated_at
		FROM "user"
		WHERE email = $1
	`

	var user dto.UserDB
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Login,
		&user.Password,
		&user.Description,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return user, err
}

func (r *PostgresRepository) GetUserByLogin(ctx context.Context, login string) (dto.UserDB, error) {
	query := `
		SELECT _id, user_name, surname, email, user_login, user_hashed_password, user_description, created_at, updated_at
		FROM "user"
		WHERE user_login = $1
	`

	var user dto.UserDB
	err := r.db.QueryRowContext(ctx, query, login).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Login,
		&user.Password,
		&user.Description,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return user, err
}

func (r *PostgresRepository) GetUserByID(ctx context.Context, id int) (dto.UserDB, error) {
	query := `
		SELECT _id, user_name, surname, email, user_login, user_hashed_password, user_description, created_at, updated_at
		FROM "user"
		WHERE _id = $1
	`

	var user dto.UserDB
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Login,
		&user.Password,
		&user.Description,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return user, err
}

func (r *PostgresRepository) CreateUserModel(ctx context.Context, user models.User) (models.User, error) {
	userDB := dto.UserDB{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Login:       user.Login,
		Password:    user.Password,
		Description: user.Description,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	id, err := r.CreateUser(ctx, userDB)
	if err != nil {
		return models.User{}, err
	}

	user.ID = id
	return user, nil
}

func (r *PostgresRepository) GetUserByLoginModel(ctx context.Context, login string) (models.User, error) {
	userDB, err := r.GetUserByLogin(ctx, login)
	if err != nil {
		return models.User{}, err
	}

	return dtoToModel(userDB), nil
}

func (r *PostgresRepository) GetUserByIDModel(ctx context.Context, id int) (models.User, error) {
	userDB, err := r.GetUserByID(ctx, id)
	if err != nil {
		return models.User{}, err
	}

	return dtoToModel(userDB), nil
}

func dtoToModel(userDB dto.UserDB) models.User {
	return models.User{
		ID:          userDB.ID,
		FirstName:   userDB.FirstName,
		LastName:    userDB.LastName,
		Email:       userDB.Email,
		Login:       userDB.Login,
		Password:    userDB.Password,
		Description: userDB.Description,
		CreatedAt:   userDB.CreatedAt,
		UpdatedAt:   userDB.UpdatedAt,
	}
}
