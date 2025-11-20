package category

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/dto"
)

func CategoryDBToModel(categoryDB dto.CategoryDB) models.Category {
	var description string
	if categoryDB.Description != nil {
		description = *categoryDB.Description
	}

	return models.Category{
		ID:           categoryDB.ID,
		UserID:       categoryDB.UserID,
		Name:         categoryDB.Name,
		Description:  description,
		LogoHashedID: categoryDB.LogoHashedID,
		CreatedAt:    categoryDB.CreatedAt,
		UpdatedAt:    categoryDB.UpdatedAt,
	}
}

func CategoryModelToDB(category models.Category) dto.CategoryDB {
	var description *string
	if category.Description != "" {
		description = &category.Description
	}

	return dto.CategoryDB{
		ID:           category.ID,
		UserID:       category.UserID,
		Name:         category.Name,
		Description:  description,
		LogoHashedID: category.LogoHashedID,
		CreatedAt:    category.CreatedAt,
		UpdatedAt:    category.UpdatedAt,
	}
}
