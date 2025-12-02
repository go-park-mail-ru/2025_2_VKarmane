package repository

import (
	"context"

	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	serviceerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/errors"
	finmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/models"
	"github.com/stretchr/testify/require"
)

func setupPostgresDB(t *testing.T) (*PostgresRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewPostgresRepository(db)
	return repo, mock, func() { _ = db.Close() }
}

func TestCreateCategory(t *testing.T) {
	repo, mock, close := setupPostgresDB(t)
	defer close()

	category := finmodels.Category{
		UserID:       1,
		Name:         "Food",
		Description:  "desc",
		LogoHashedID: "logo123",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO category (user_id, category_name, category_description, logo_hashed_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING _id, created_at, updated_at
	`)).
		WithArgs(category.UserID, category.Name, &category.Description, category.LogoHashedID).
		WillReturnRows(sqlmock.NewRows([]string{"_id", "created_at", "updated_at"}).
			AddRow(10, time.Now(), time.Now()))

	out, err := repo.CreateCategory(context.Background(), category)
	require.NoError(t, err)
	require.Equal(t, 10, out.ID)
}

func TestGetCategoriesByUser(t *testing.T) {
	repo, mock, close := setupPostgresDB(t)
	defer close()

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT _id, user_id, category_name, category_description, logo_hashed_id, created_at, updated_at
		FROM category
		WHERE user_id = $1
		ORDER BY created_at DESC
	`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"_id", "user_id", "category_name", "category_description",
			"logo_hashed_id", "created_at", "updated_at",
		}).AddRow(10, 1, "Food", "desc", "logo123", time.Now(), time.Now()))

	cats, err := repo.GetCategoriesByUser(context.Background(), 1)
	require.NoError(t, err)
	require.Len(t, cats, 1)
	require.Equal(t, "Food", cats[0].Name)
}

func TestGetCategoriesWithStatsByUser(t *testing.T) {
	repo, mock, close := setupPostgresDB(t)
	defer close()

	mock.ExpectQuery(regexp.QuoteMeta(`
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
	`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"_id", "user_id", "category_name", "category_description",
			"logo_hashed_id", "created_at", "updated_at", "operations_count",
		}).AddRow(10, 1, "Food", "desc", "logo123", time.Now(), time.Now(), 3))

	out, err := repo.GetCategoriesWithStatsByUser(context.Background(), 1)
	require.NoError(t, err)
	require.Len(t, out, 1)
	require.Equal(t, 3, out[0].OperationsCount)
}

func TestGetCategoryByID(t *testing.T) {
	repo, mock, close := setupPostgresDB(t)
	defer close()

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT _id, user_id, category_name, category_description, logo_hashed_id, created_at, updated_at
		FROM category
		WHERE _id = $1 AND user_id = $2
	`)).
		WithArgs(10, 1).
		WillReturnRows(sqlmock.NewRows([]string{
			"_id", "user_id", "category_name", "category_description",
			"logo_hashed_id", "created_at", "updated_at",
		}).AddRow(10, 1, "Food", "desc", "logo123", time.Now(), time.Now()))

	cat, err := repo.GetCategoryByID(context.Background(), 1, 10)
	require.NoError(t, err)
	require.Equal(t, "Food", cat.Name)
}

func TestUpdateCategory_Success(t *testing.T) {
	repo, mock, close := setupPostgresDB(t)
	defer close()

	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE category 
		SET category_name = COALESCE($1,category_name), 
		    category_description = COALESCE($2,category_description), 
		    logo_hashed_id = COALESCE($3,logo_hashed_id), 
		    updated_at = NOW()
		WHERE _id = $4 AND user_id = $5
	`)).
		WithArgs("Food", "desc", "logo123", 10, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateCategory(context.Background(), finmodels.Category{
		ID:           10,
		UserID:       1,
		Name:         "Food",
		Description:  "desc",
		LogoHashedID: "logo123",
	})
	require.NoError(t, err)
}

func TestUpdateCategory_NotFound(t *testing.T) {
	repo, mock, close := setupPostgresDB(t)
	defer close()

	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE category 
		SET category_name = COALESCE($1,category_name), 
		    category_description = COALESCE($2,category_description), 
		    logo_hashed_id = COALESCE($3,logo_hashed_id), 
		    updated_at = NOW()
		WHERE _id = $4 AND user_id = $5
	`)).
		WithArgs("Food", "desc", "logo123", 10, 1).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.UpdateCategory(context.Background(), finmodels.Category{
		ID:           10,
		UserID:       1,
		Name:         "Food",
		Description:  "desc",
		LogoHashedID: "logo123",
	})
	require.Error(t, err)
	require.Equal(t, serviceerrors.ErrCategoryNotFound, err)
}

func TestDeleteCategory(t *testing.T) {
	repo, mock, close := setupPostgresDB(t)
	defer close()

	mock.ExpectExec(regexp.QuoteMeta(`
		DELETE FROM category
		WHERE _id = $1 AND user_id = $2
	`)).
		WithArgs(10, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteCategory(context.Background(), 1, 10)
	require.NoError(t, err)
}

func TestGetCategoryStats(t *testing.T) {
	repo, mock, close := setupPostgresDB(t)
	defer close()

	mock.ExpectQuery(regexp.QuoteMeta(`
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
	`)).
		WithArgs(5, 1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(7))

	count, err := repo.GetCategoryStats(context.Background(), 1, 5)
	require.NoError(t, err)
	require.Equal(t, 7, count)
}
