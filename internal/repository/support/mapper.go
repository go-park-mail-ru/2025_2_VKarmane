package support

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/dto"
)

func SupportDBToModel(db dto.SupportDB) models.Support {
	return models.Support{
		ID:              db.ID,
		UserID:          db.UserID,
		CategoryRequest: models.CategoryContacting(db.CategoryRequest),
		StatusRequest:   models.StatusContacting(db.StatusRequest),
		Message:         db.Message,
	}
}

func SupportModelToDB(m models.Support) dto.SupportDB {
	return dto.SupportDB{
		ID:              m.ID,
		UserID:          m.UserID,
		CategoryRequest: string(m.CategoryRequest),
		StatusRequest:   string(m.StatusRequest),
		Message:         m.Message,
	}
}
