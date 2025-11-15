package support

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/dto"
	"github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, req models.Support) (models.Support, error) {
	dbReq := SupportModelToDB(req)

	query := `
		INSERT INTO support_requests (user_id, category, status, message)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, category, status, message
	`

	var out dto.SupportDB
	err := r.db.QueryRowContext(ctx, query,
		dbReq.UserID,
		dbReq.CategoryRequest,
		dbReq.StatusRequest,
		dbReq.Message,
	).Scan(
		&out.ID,
		&out.UserID,
		&out.CategoryRequest,
		&out.StatusRequest,
		&out.Message,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			// здесь можно добавить проверку на UniqueViolation / ForeignKeyViolation, если в будущем добавишь связи
		}
		return models.Support{}, fmt.Errorf("failed to create support request: %w", err)
	}

	return SupportDBToModel(out), nil
}

func (r *PostgresRepository) GetByUser(ctx context.Context, userID int) ([]models.Support, error) {
	query := `
		SELECT id, user_id, category, status, message
		FROM support_requests
		WHERE user_id = $1
		ORDER BY id DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user support requests: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var result []models.Support
	for rows.Next() {
		var dbReq dto.SupportDB
		err := rows.Scan(
			&dbReq.ID,
			&dbReq.UserID,
			&dbReq.CategoryRequest,
			&dbReq.StatusRequest,
			&dbReq.Message,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan support request: %w", err)
		}
		result = append(result, SupportDBToModel(dbReq))
	}

	return result, nil
}

func (r *PostgresRepository) UpdateStatus(ctx context.Context, id int, status models.StatusContacting) error {
	query := `
		UPDATE support_requests 
		SET status = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, string(status), id)
	if err != nil {
		return fmt.Errorf("failed to update support request status: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetStats(ctx context.Context) (map[models.StatusContacting]int, error) {
	query := `
		SELECT status, COUNT(*)
		FROM support_requests
		GROUP BY status
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get support stats: %w", err)
	}
	defer func() { _ = rows.Close() }()

	stats := make(map[models.StatusContacting]int)
	for rows.Next() {
		var status string
		var count int

		err := rows.Scan(&status, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stats row: %w", err)
		}

		stats[models.StatusContacting(status)] = count
	}

	return stats, nil
}
