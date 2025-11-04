package category

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresRepository_CreateCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	desc := "Food category"
	category := dto.CategoryDB{
		UserID:      1,
		Name:        "Food",
		Description: &desc,
		LogoHashedID: "logo123",
	}

	mock.ExpectQuery(`INSERT INTO category`).
		WithArgs(category.UserID, category.Name, category.Description, category.LogoHashedID, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"_id"}).AddRow(5))

	id, err := repo.CreateCategory(context.Background(), category)
	assert.NoError(t, err)
	assert.Equal(t, 5, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_CreateCategory_NilDescription(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	category := dto.CategoryDB{
		UserID:      1,
		Name:        "Food",
		Description: nil,
		LogoHashedID: "logo123",
	}

	mock.ExpectQuery(`INSERT INTO category`).
		WithArgs(category.UserID, category.Name, category.Description, category.LogoHashedID, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"_id"}).AddRow(5))

	id, err := repo.CreateCategory(context.Background(), category)
	assert.NoError(t, err)
	assert.Equal(t, 5, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_CreateCategory_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	desc := "Food category"
	category := dto.CategoryDB{
		UserID:      1,
		Name:        "Food",
		Description: &desc,
	}

	mock.ExpectQuery(`INSERT INTO category`).
		WithArgs(category.UserID, category.Name, category.Description, category.LogoHashedID, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	id, err := repo.CreateCategory(context.Background(), category)
	assert.Error(t, err)
	assert.Zero(t, id)
	assert.Contains(t, err.Error(), "failed to create category")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetCategoriesByUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	userID := 1
	now := time.Now()
	desc1 := "Food category"
	desc2 := "Transport category"
	rows := sqlmock.NewRows([]string{"_id", "user_id", "category_name", "category_description", "logo_hashed_id", "created_at", "updated_at"}).
		AddRow(1, 1, "Food", &desc1, "logo1", now, now).
		AddRow(2, 1, "Transport", &desc2, "", now, now)

	mock.ExpectQuery(`SELECT.*FROM category`).
		WithArgs(userID).
		WillReturnRows(rows)

	categories, err := repo.GetCategoriesByUser(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, categories, 2)
	assert.Equal(t, 1, categories[0].ID)
	assert.Equal(t, "Food", categories[0].Name)
	require.NotNil(t, categories[0].Description)
	assert.Equal(t, "Food category", *categories[0].Description)
	assert.Equal(t, 2, categories[1].ID)
	assert.Equal(t, "Transport", categories[1].Name)
	require.NotNil(t, categories[1].Description)
	assert.Equal(t, "Transport category", *categories[1].Description)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetCategoriesByUser_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	userID := 1
	rows := sqlmock.NewRows([]string{"_id", "user_id", "category_name", "category_description", "logo_hashed_id", "created_at", "updated_at"})

	mock.ExpectQuery(`SELECT.*FROM category`).
		WithArgs(userID).
		WillReturnRows(rows)

	categories, err := repo.GetCategoriesByUser(context.Background(), userID)
	assert.NoError(t, err)
	assert.Empty(t, categories)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetCategoriesByUser_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	userID := 1
	mock.ExpectQuery(`SELECT.*FROM category`).
		WithArgs(userID).
		WillReturnError(sql.ErrConnDone)

	categories, err := repo.GetCategoriesByUser(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, categories)
	assert.Contains(t, err.Error(), "failed to get categories by user")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetCategoryByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	userID := 1
	categoryID := 5
	now := time.Now()
	desc := "Food category"
	rows := sqlmock.NewRows([]string{"_id", "user_id", "category_name", "category_description", "logo_hashed_id", "created_at", "updated_at"}).
		AddRow(categoryID, userID, "Food", &desc, "logo1", now, now)

	mock.ExpectQuery(`SELECT.*FROM category`).
		WithArgs(categoryID, userID).
		WillReturnRows(rows)

	category, err := repo.GetCategoryByID(context.Background(), userID, categoryID)
	assert.NoError(t, err)
	assert.Equal(t, categoryID, category.ID)
	assert.Equal(t, userID, category.UserID)
	assert.Equal(t, "Food", category.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetCategoryByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	userID := 1
	categoryID := 999
	mock.ExpectQuery(`SELECT.*FROM category`).
		WithArgs(categoryID, userID).
		WillReturnError(sql.ErrNoRows)

	category, err := repo.GetCategoryByID(context.Background(), userID, categoryID)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.Zero(t, category.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetCategoryByID_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	userID := 1
	categoryID := 5
	mock.ExpectQuery(`SELECT.*FROM category`).
		WithArgs(categoryID, userID).
		WillReturnError(sql.ErrConnDone)

	category, err := repo.GetCategoryByID(context.Background(), userID, categoryID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get category by ID")
	assert.Zero(t, category.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_UpdateCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	desc := "Updated description"
	category := dto.CategoryDB{
		ID:          5,
		UserID:      1,
		Name:        "Updated Food",
		Description: &desc,
		LogoHashedID: "newlogo",
	}

	mock.ExpectExec(`UPDATE category`).
		WithArgs(category.Name, category.Description, category.LogoHashedID, sqlmock.AnyArg(), category.ID, category.UserID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateCategory(context.Background(), category)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_UpdateCategory_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	desc := "Updated description"
	category := dto.CategoryDB{
		ID:          5,
		UserID:      1,
		Name:        "Updated Food",
		Description: &desc,
	}

	mock.ExpectExec(`UPDATE category`).
		WithArgs(category.Name, category.Description, category.LogoHashedID, sqlmock.AnyArg(), category.ID, category.UserID).
		WillReturnError(sql.ErrConnDone)

	err = repo.UpdateCategory(context.Background(), category)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update category")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_DeleteCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	userID := 1
	categoryID := 5

	mock.ExpectExec(`DELETE FROM category`).
		WithArgs(categoryID, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeleteCategory(context.Background(), userID, categoryID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_DeleteCategory_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	userID := 1
	categoryID := 5

	mock.ExpectExec(`DELETE FROM category`).
		WithArgs(categoryID, userID).
		WillReturnError(sql.ErrConnDone)

	err = repo.DeleteCategory(context.Background(), userID, categoryID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete category")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetCategoryStats(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	userID := 1
	categoryID := 5
	expectedCount := 10

	rows := sqlmock.NewRows([]string{"count"}).AddRow(expectedCount)

	mock.ExpectQuery(`SELECT COUNT`).
		WithArgs(categoryID, userID).
		WillReturnRows(rows)

	count, err := repo.GetCategoryStats(context.Background(), userID, categoryID)
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetCategoryStats_Zero(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	userID := 1
	categoryID := 5

	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

	mock.ExpectQuery(`SELECT COUNT`).
		WithArgs(categoryID, userID).
		WillReturnRows(rows)

	count, err := repo.GetCategoryStats(context.Background(), userID, categoryID)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetCategoryStats_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepository(db)

	userID := 1
	categoryID := 5

	mock.ExpectQuery(`SELECT COUNT`).
		WithArgs(categoryID, userID).
		WillReturnError(sql.ErrConnDone)

	count, err := repo.GetCategoryStats(context.Background(), userID, categoryID)
	assert.Error(t, err)
	assert.Zero(t, count)
	assert.Contains(t, err.Error(), "failed to get category stats")
	assert.NoError(t, mock.ExpectationsWereMet())
}

