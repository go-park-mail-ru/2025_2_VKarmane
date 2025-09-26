package operation

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type Repository struct {
	operations []OperationDB
}

func NewRepository(operations []OperationDB) *Repository {
	return &Repository{
		operations: operations,
	}
}

func (r *Repository) GetOperationsByAccount(accountID int) []models.Operation {
	out := make([]models.Operation, 0)

	for _, o := range r.operations {
		if o.AccountID == accountID {
			out = append(out, OperationDBToModel(o))
		}
	}

	return out
}
