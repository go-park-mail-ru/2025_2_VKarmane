package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	finmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/finance_service/models"
)

type CategoryDB struct {
	ID           int
	UserID       int
	Name         string
	Description  *string
	LogoHashedID string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (r *PostgresRepository) CreateCategory(ctx context.Context, category finmodels.Category) (finmodels.Category, error) {
	query := `
		INSERT INTO category (user_id, category_name, category_description, logo_hashed_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING _id, created_at, updated_at
	`

	var description *string
	if category.Description != "" {
		description = &category.Description
	}

	err := r.db.QueryRowContext(ctx, query,
		category.UserID,
		category.Name,
		description,
		category.LogoHashedID,
	).Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		return finmodels.Category{}, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

func (r *PostgresRepository) GetCategoriesByUser(ctx context.Context, userID int) ([]finmodels.Category, error) {
	query := `
		SELECT _id, user_id, category_name, category_description, logo_hashed_id, created_at, updated_at
		FROM category
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories by user: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var categories []finmodels.Category
	for rows.Next() {
		var categoryDB CategoryDB
		err := rows.Scan(
			&categoryDB.ID,
			&categoryDB.UserID,
			&categoryDB.Name,
			&categoryDB.Description,
			&categoryDB.LogoHashedID,
			&categoryDB.CreatedAt,
			&categoryDB.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, categoryDBToModel(categoryDB))
	}

	return categories, nil
}

func (r *PostgresRepository) GetCategoriesWithStatsByUser(ctx context.Context, userID int) ([]finmodels.CategoryWithStats, error) {
	query := `
		SELECT c._id, c.user_id, c.category_name, c.category_description, c.logo_hashed_id, 
		       c.created_at, c.updated_at,
		       COALESCE(COUNT(op._id), 0) as operations_count
		FROM category c
		LEFT JOIN operation op ON op.category_id = c._id 
			AND op.operation_status != 'reverted'
			AND (op.account_from_id IN (SELECT account_id FROM sharings WHERE user_id = $1)
			     OR op.account_to_id IN (SELECT account_id FROM sharings WHERE user_id = $1))
		WHERE c.user_id = $1
		GROUP BY c._id, c.user_id, c.category_name, c.category_description, c.logo_hashed_id, 
		         c.created_at, c.updated_at
		ORDER BY c.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories with stats by user: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var categories []finmodels.CategoryWithStats
	for rows.Next() {
		var categoryDB CategoryDB
		var operationsCount int
		err := rows.Scan(
			&categoryDB.ID,
			&categoryDB.UserID,
			&categoryDB.Name,
			&categoryDB.Description,
			&categoryDB.LogoHashedID,
			&categoryDB.CreatedAt,
			&categoryDB.UpdatedAt,
			&operationsCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, finmodels.CategoryWithStats{
			Category:        categoryDBToModel(categoryDB),
			OperationsCount: operationsCount,
		})
	}

	return categories, nil
}

func (r *PostgresRepository) GetCategoryByID(ctx context.Context, userID, categoryID int) (finmodels.Category, error) {
	query := `
		SELECT _id, user_id, category_name, category_description, logo_hashed_id, created_at, updated_at
		FROM category
		WHERE _id = $1 AND user_id = $2
	`

	var categoryDB CategoryDB
	err := r.db.QueryRowContext(ctx, query, categoryID, userID).Scan(
		&categoryDB.ID,
		&categoryDB.UserID,
		&categoryDB.Name,
		&categoryDB.Description,
		&categoryDB.LogoHashedID,
		&categoryDB.CreatedAt,
		&categoryDB.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return finmodels.Category{}, ErrCategoryNotFound
		}
		return finmodels.Category{}, fmt.Errorf("failed to get category by ID: %w", err)
	}

	return categoryDBToModel(categoryDB), nil
}

func (r *PostgresRepository) UpdateCategory(ctx context.Context, category finmodels.Category) error {
	query := `
		UPDATE category 
		SET category_name = $1, category_description = $2, logo_hashed_id = $3, updated_at = NOW()
		WHERE _id = $4 AND user_id = $5
	`

	var description *string
	if category.Description != "" {
		description = &category.Description
	}

	_, err := r.db.ExecContext(ctx, query,
		category.Name,
		description,
		category.LogoHashedID,
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

	_, err := r.db.ExecContext(ctx, query, categoryID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetCategoryStats(ctx context.Context, userID, categoryID int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM operation AS op
		JOIN account AS acc
		ON acc._id = op.account_from_id OR acc._id = op.account_to_id
		JOIN sharings AS sh
		ON sh.account_id = acc._id
		JOIN "user" AS u
		ON u._id = sh.user_id
		WHERE op.category_id = $1
		AND op.operation_status != 'reverted'
		AND u._id = $2;
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, categoryID, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get category stats: %w", err)
	}

	return count, nil
}

func categoryDBToModel(categoryDB CategoryDB) finmodels.Category {
	description := ""
	if categoryDB.Description != nil {
		description = *categoryDB.Description
	}

	return finmodels.Category{
		ID:           categoryDB.ID,
		UserID:       categoryDB.UserID,
		Name:         categoryDB.Name,
		Description:  description,
		LogoHashedID: categoryDB.LogoHashedID,
		CreatedAt:    categoryDB.CreatedAt,
		UpdatedAt:    categoryDB.UpdatedAt,
	}
}

var ErrCategoryNotFound = errors.New("category not found")
var ErrCategoryExists = errors.New("category already exists")

