package operation

import (
	"context"

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

func (r *Repository) GetOperationsByAccount(ctx context.Context, accountID int) []models.Operation {
	out := make([]models.Operation, 0)

	for _, o := range r.operations {
		if (o.AccountFromID != nil && *o.AccountFromID == accountID) ||
			(o.AccountToID != nil && *o.AccountToID == accountID) {
			out = append(out, OperationDBToModel(o))
		}
	}

	return out
}
