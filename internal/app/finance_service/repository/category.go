package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	serviceerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/errors"
	finmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/models"
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
		return finmodels.Category{}, MapPgCategoryError(err)
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
			return nil, MapPgCategoryError(err)
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
		return nil, MapPgCategoryError(err)
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
			return nil, MapPgCategoryError(err)
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
		return finmodels.Category{}, MapPgCategoryError(err)
	}

	return categoryDBToModel(categoryDB), nil
}

func (r *PostgresRepository) GetCategoryByName(ctx context.Context, userID int, categoryName string) (finmodels.Category, error) {
	query := `
		SELECT _id, user_id, category_name, category_description, logo_hashed_id, created_at, updated_at
		FROM category
		WHERE user_id = $1 AND category_name = $2
	`

	var categoryDB CategoryDB
	err := r.db.QueryRowContext(ctx, query, userID, categoryName).Scan(
		&categoryDB.ID,
		&categoryDB.UserID,
		&categoryDB.Name,
		&categoryDB.Description,
		&categoryDB.LogoHashedID,
		&categoryDB.CreatedAt,
		&categoryDB.UpdatedAt,
	)

	if err != nil {
		return finmodels.Category{}, MapPgCategoryError(err)
	}

	return categoryDBToModel(categoryDB), nil
}

func (r *PostgresRepository) UpdateCategory(ctx context.Context, category finmodels.Category) error {

	log.Printf("hash %s", category.LogoHashedID)
	query := `
		UPDATE category 
		SET category_name = COALESCE($1,category_name), category_description = COALESCE($2,category_description), logo_hashed_id = COALESCE($3,logo_hashed_id), updated_at = NOW()
		WHERE _id = $4 AND user_id = $5
	`

	var description *string
	if category.Description != "" {
		description = &category.Description
	}

	var logoHash *string
	if category.LogoHashedID != "" {
		logoHash = &category.LogoHashedID
	}

	var ctgName *string
	if category.Name != "" {
		ctgName = &category.Name
	}
	res, err := r.db.ExecContext(ctx, query,
		ctgName,
		description,
		logoHash,
		category.ID,
		category.UserID,
	)

	if err != nil {
		return MapPgCategoryError(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return serviceerrors.ErrCategoryNotFound
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
		return MapPgCategoryError(err)
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
		return 0, MapPgCategoryError(err)
	}

	return count, nil
}

func (r *PostgresRepository) GetCategoriesReport(
	ctx context.Context,
	userID int,
	start, end time.Time,
) ([]finmodels.CategoryInReport, error) {

	query := `
        SELECT
            c._id AS category_id,
            c.category_name,
            COUNT(op._id) AS operations_count,
            COALESCE(SUM(op.sum), 0) AS total_sum
        FROM category AS c
        LEFT JOIN operation AS op
            ON op.category_id = c._id
            AND op.operation_status != 'reverted'
            AND op.operation_date >= $2
            AND op.operation_date <= $3
        LEFT JOIN account AS acc
            ON acc._id = op.account_from_id OR acc._id = op.account_to_id
        LEFT JOIN sharings AS sh
            ON sh.account_id = acc._id
        WHERE c.user_id = $1
            AND (sh.user_id = $1 OR sh.user_id IS NULL)
        GROUP BY c._id, c.category_name
        ORDER BY c.category_name;
    `

	rows, err := r.db.QueryContext(ctx, query, userID, start, end)
	if err != nil {
		return nil, MapPgCategoryError(err)
	}
	defer rows.Close()

	reports := []finmodels.CategoryInReport{}

	for rows.Next() {
		var rep finmodels.CategoryInReport

		err := rows.Scan(
			&rep.CategoryID,
			&rep.CategoryName,
			&rep.OperationCount,
			&rep.TotalSum,
		)
		if err != nil {
			return nil, MapPgCategoryError(err)
		}

		reports = append(reports, rep)
	}

	if err := rows.Err(); err != nil {
		return nil, MapPgCategoryError(err)
	}

	return reports, nil
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
