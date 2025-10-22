package operation

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

type Repository struct {
	operations []OperationDB
	clock      clock.Clock
}

func NewRepository(operations []OperationDB, clck clock.Clock) *Repository {
	return &Repository{
		operations: operations,
		clock:      clck,
	}
}

func (r *Repository) GetOperationsByAccount(ctx context.Context, accountID int) ([]models.Operation, error) {
	out := make([]models.Operation, 0)

	for _, o := range r.operations {
		if o.AccountID == accountID && o.Status != models.OperationReverted {
			out = append(out, OperationDBToModel(o))
		}
	}

	return out, nil
}

func (r *Repository) GetOperationByID(ctx context.Context, accountID int, opID int) (models.Operation, error) {
	for _, o := range r.operations {
		if o.AccountID == accountID && o.ID == opID && o.Status != models.OperationReverted {
			return OperationDBToModel(o), nil
		}
	}

	return models.Operation{}, nil
}

func (r *Repository) CreateOperation(ctx context.Context, op models.Operation) (models.Operation, error) {
	opDB := OperationDB{
		ID:          len(r.operations) + 1,
		AccountID:   op.AccountID,
		CategoryID:  op.CategoryID,
		Type:        op.Type,
		Status:      op.Status,
		Description: op.Description,
		ReceiptURL:  op.ReceiptURL,
		Name:        op.Name,
		Sum:         op.Sum,
		CurrencyID:  op.CurrencyID,
		CreatedAt:   op.CreatedAt,
		ReceiverID:  op.ReceiverID,
	}

	r.operations = append(r.operations, opDB)
	return OperationDBToModel(opDB), nil
}

func (r *Repository) UpdateOperation(ctx context.Context, req models.UpdateOperationRequest, accID, opID int) (models.Operation, error) {
	for i := range r.operations {
		op := &r.operations[i]
		if op.AccountID == accID && op.ID == opID {

			if req.CategoryID != nil {
				op.CategoryID = *req.CategoryID
			}
			if req.Sum != nil {
				op.Sum = *req.Sum
			}
			if req.Name != nil {
				op.Name = *req.Name
			}
			if req.Description != nil {
				op.Description = *req.Description
			}
			if req.CategoryID != nil {
				op.CreatedAt = *req.CreatedAt
			}

			return OperationDBToModel(*op), nil
		}
	}
	return models.Operation{}, nil
}

func (r *Repository) DeleteOperation(ctx context.Context, accID int, opID int) (models.Operation, error) {
	for i := range r.operations {
		if r.operations[i].AccountID == accID && r.operations[i].ID == opID {
			r.operations[i].Status = models.OperationReverted
			return OperationDBToModel(r.operations[i]), nil
		}
	}
	return models.Operation{}, nil
}
