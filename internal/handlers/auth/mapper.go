package auth

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

func UserToApi(op models.User) UserAPI {
	return UserAPI{
		ID:        op.ID,
		FirstName: op.FirstName,
		LastName:  op.LastName,
		Email:     op.Email,
		Login:     op.Login,
		Password:  op.Password,
		CreatedAt: op.CreatedAt,
		UpdatedAt: op.UpdatedAt,
	}
}
