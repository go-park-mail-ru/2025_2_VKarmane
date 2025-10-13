package operation

import (
	"context"
	"time"
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

func (r *Repository) GetOperationsByAccount(ctx context.Context, accountID int) ([]models.Operation) {
	out := make([]models.Operation, 0)

	for _, o := range r.operations {
		if o.AccountID == accountID && o.Status != string(models.OperationReverted) {
			out = append(out, OperationDBToModel(o))
		}
	}

	return out
}

func (r *Repository) GetOperationByID(ctx context.Context, accountID int, opID int) models.Operation {
	for _, o := range r.operations {
		if o.AccountID == accountID && o.ID == opID && o.Status != string(models.OperationReverted){
			return OperationDBToModel(o)
		}
	}

	return models.Operation{}
}


func (r *Repository) CreateOperation(ctx context.Context, op models.Operation) models.Operation {
	opDB := OperationDB{
		ID:        len(r.operations) + 1,
		AccountID: op.AccountID,
		CategoryID:  op.CategoryID,
		Type:     string(op.Type),
		Status:     string(op.Status),
		Description:  op.Description,
		ReceiptURL: op.ReceiptURL,
		Name: op.Name,
		Sum: op.Sum,
		CurrencyID: op.CurrencyID,
		CreatedAt: time.Now(),
	}

	r.operations = append(r.operations, opDB)
	return OperationDBToModel(opDB)
}


func (r *Repository) UpdateOperation(ctx context.Context, req models.UpdateOperationRequest, accID, opID int) models.Operation {
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

            return OperationDBToModel(*op)
        }
    }
    return models.Operation{}
}



func (r *Repository) DeleteOperation(ctx context.Context, accID int, opID int) models.Operation {
	for i :=range  r.operations {
		if r.operations[i].AccountID == accID && r.operations[i].ID == opID {
			r.operations[i].Status = string(models.OperationReverted)
			return OperationDBToModel(r.operations[i])
		}
	}
	return models.Operation{}
}

