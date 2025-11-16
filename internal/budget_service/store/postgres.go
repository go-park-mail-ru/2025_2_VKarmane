package store

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	bdgmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/models"
	// "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/proto"
	bdgrepo "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/repository"
)

type PostgresStore struct {
    DB       *sql.DB
    BudgetRepo *bdgrepo.PostgresRepository
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

    store.BudgetRepo = bdgrepo.NewPostgresRepository(db)

    return store, nil
}

func (s *PostgresStore) Close() error {
    if err := s.DB.Close(); err != nil {
        return fmt.Errorf("failed to close database: %w", err)
    }
    return nil
}



func (s *PostgresStore) CreateBudget(ctx context.Context, budget bdgmodels.Budget) (bdgmodels.Budget, error) {
    return s.BudgetRepo.CreateBudget(ctx, budget)
}


func (s *PostgresStore) GetBudgetsByUser(ctx context.Context, id int) ([]bdgmodels.Budget, error) {
    return s.BudgetRepo.GetBudgetsByUser(ctx, id)
}

func (s *PostgresStore) UpdateBudget(ctx context.Context, req bdgmodels.UpdatedBudgetRequest) (bdgmodels.Budget, error) {
    return s.BudgetRepo.UpdateBudget(ctx, req)
}

func (s *PostgresStore) DeleteBudget(ctx context.Context, budgetID int) (bdgmodels.Budget, error) {
	return s.BudgetRepo.DeleteBudget(ctx, budgetID)
}