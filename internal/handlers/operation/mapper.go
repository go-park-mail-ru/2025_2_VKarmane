package operation

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

func OperationToApi(op models.Operation) OperationAPI {
	return OperationAPI{
		ID:          op.ID,
		AccountID:   op.AccountID,
		CategoryID:  op.CategoryID,
		Sum:         op.Sum,
		Name:        op.Name,
		Status:      string(op.Status),
		Type:        string(op.Type),
		Description: op.Description,
		ReceiptURL:  op.ReceiptURL,
		Date:        op.CreatedAt,
	}
}

func OperationsToApi(userID int, ops []models.Operation) OperationsAPI {
	opDTOs := make([]OperationAPI, 0, len(ops))
	for _, op := range ops {
		opDTOs = append(opDTOs, OperationToApi(op))
	}
	return OperationsAPI{
		UserID:     userID,
		Operations: opDTOs,
	}
}
