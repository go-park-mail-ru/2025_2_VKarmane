package store

import (
    "context"
    "database/sql"
    "fmt"

    userrepo "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/repository"
    authmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/models"

    _ "github.com/lib/pq"
)

type PostgresStore struct {
    DB       *sql.DB
    UserRepo *userrepo.PostgresRepository
}

func NewPostgresStore(dsn string) (*PostgresStore, error) {
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    store := &PostgresStore{
        DB: db,
    }

    store.UserRepo = userrepo.NewPostgresRepository(db)

    return store, nil
}

func (s *PostgresStore) Close() error {
    if err := s.DB.Close(); err != nil {
        return fmt.Errorf("failed to close database: %w", err)
    }
    return nil
}

//
// USER STORE METHODS — чистые прокси на UserRepository
//

func (s *PostgresStore) CreateUser(ctx context.Context, user authmodels.User) (authmodels.User, error) {
    return s.UserRepo.CreateUser(ctx, user)
}


func (s *PostgresStore) GetUserByID(ctx context.Context, id int) (authmodels.User, error) {
    return s.UserRepo.GetUserByID(ctx, id)
}

func (s *PostgresStore) EditUserByID(ctx context.Context, req authmodels.UpdateProfileRequest) (authmodels.User, error) {
    return s.UserRepo.EditUserByID(ctx, req)
}

func (s *PostgresStore) GetUserByLogin(ctx context.Context, login string) (authmodels.User, error) {
	return s.UserRepo.GetUserByLogin(ctx, login)
}