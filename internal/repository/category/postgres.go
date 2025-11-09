package category

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

func (r *PostgresRepository) CreateCategory(ctx context.Context, category dto.CategoryDB) (int, error) {
	query := `
		INSERT INTO category (user_id, category_name, category_description, logo_hashed_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING _id
	`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare stmt: %w", err)
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRowContext(ctx,
		category.UserID,
		category.Name,
		category.Description,
		category.LogoHashedID,
		time.Now(),
		time.Now(),
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create category: %w", err)
	}

	return id, nil
}

func (r *PostgresRepository) GetCategoriesByUser(ctx context.Context, userID int) ([]dto.CategoryDB, error) {
	query := `
		SELECT _id, user_id, category_name, category_description, logo_hashed_id, created_at, updated_at
		FROM category
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return []dto.CategoryDB{},fmt.Errorf("failed to prepare stmt: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories by user: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var categories []dto.CategoryDB
	for rows.Next() {
		var category dto.CategoryDB
		err := rows.Scan(
			&category.ID,
			&category.UserID,
			&category.Name,
			&category.Description,
			&category.LogoHashedID,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *PostgresRepository) GetCategoryByID(ctx context.Context, userID, categoryID int) (dto.CategoryDB, error) {
	query := `
		SELECT _id, user_id, category_name, category_description, logo_hashed_id, created_at, updated_at
		FROM category
		WHERE _id = $1 AND user_id = $2
	`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return dto.CategoryDB{},fmt.Errorf("failed to prepare stmt: %w", err)
	}
	defer stmt.Close()
	

	var category dto.CategoryDB
	err = stmt.QueryRowContext(ctx, categoryID, userID).Scan(
		&category.ID,
		&category.UserID,
		&category.Name,
		&category.Description,
		&category.LogoHashedID,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return dto.CategoryDB{}, sql.ErrNoRows
		}
		return dto.CategoryDB{}, fmt.Errorf("failed to get category by ID: %w", err)
	}

	return category, nil
}

func (r *PostgresRepository) UpdateCategory(ctx context.Context, category dto.CategoryDB) error {
	query := `
		UPDATE category 
		SET category_name = $1, category_description = $2, logo_hashed_id = $3, updated_at = $4
		WHERE _id = $5 AND user_id = $6
	`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare stmt: %w", err)
	}
	defer stmt.Close()
	

	_, err = stmt.ExecContext(ctx,
		category.Name,
		category.Description,
		category.LogoHashedID,
		time.Now(),
		category.ID,
		category.UserID,
	)

	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}

	return nil
}

func (r *PostgresRepository) DeleteCategory(ctx context.Context, userID, categoryID int) error {
	query := `
		DELETE FROM category
		WHERE _id = $1 AND user_id = $2
	`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare stmt: %w", err)
	}
	defer stmt.Close()
	
	_, err = stmt.ExecContext(ctx, categoryID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetCategoryStats(ctx context.Context, userID, categoryID int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM operation
		WHERE category_id = $1 AND (account_from_id IN (
			SELECT _id FROM account WHERE user_id = $2
		) OR account_to_id IN (
			SELECT _id FROM account WHERE user_id = $2
		)) AND operation_status != 'reverted'
	`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare stmt: %w", err)
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRowContext(ctx, categoryID, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get category stats: %w", err)
	}

	return count, nil
}
