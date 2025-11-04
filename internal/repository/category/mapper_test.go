package category

import (
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/dto"
	"github.com/stretchr/testify/assert"
)

func TestCategoryDBToModel(t *testing.T) {
	now := time.Now()
	desc := "Food expenses"
	categoryDB := dto.CategoryDB{
		ID:           1,
		UserID:       10,
		Name:         "Food",
		Description:  &desc,
		LogoHashedID: "abc123",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	category := CategoryDBToModel(categoryDB)

	assert.Equal(t, 1, category.ID)
	assert.Equal(t, 10, category.UserID)
	assert.Equal(t, "Food", category.Name)
	assert.Equal(t, "Food expenses", category.Description)
	assert.Equal(t, "abc123", category.LogoHashedID)
	assert.Equal(t, now, category.CreatedAt)
	assert.Equal(t, now, category.UpdatedAt)
}

func TestCategoryModelToDB(t *testing.T) {
	now := time.Now()
	category := models.Category{
		ID:           2,
		UserID:       20,
		Name:         "Transport",
		Description:  "Transport expenses",
		LogoHashedID: "xyz789",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	categoryDB := CategoryModelToDB(category)

	assert.Equal(t, 2, categoryDB.ID)
	assert.Equal(t, 20, categoryDB.UserID)
	assert.Equal(t, "Transport", categoryDB.Name)
	assert.NotNil(t, categoryDB.Description)
	assert.Equal(t, "Transport expenses", *categoryDB.Description)
	assert.Equal(t, "xyz789", categoryDB.LogoHashedID)
	assert.Equal(t, now, categoryDB.CreatedAt)
	assert.Equal(t, now, categoryDB.UpdatedAt)
}

