package operation

import (
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestOperationToApi(t *testing.T) {
	createdAt := time.Now()

	op := models.Operation{
		ID:          1,
		AccountID:   2,
		CategoryID:  3,
		Type:        models.OperationExpense,
		Status:      models.OperationFinished,
		Description: "Coffee at Starbucks",
		ReceiptURL:  "https://example.com/receipt.jpg",
		Name:        "Coffee",
		Sum:         5.50,
		CurrencyID:  1,
		CreatedAt:   createdAt,
	}

	apiOp := OperationToApi(op)

	assert.Equal(t, op.ID, apiOp.ID)
	assert.Equal(t, op.AccountID, apiOp.AccountID)
	assert.Equal(t, op.CategoryID, apiOp.CategoryID)
	assert.Equal(t, op.Sum, apiOp.Sum)
	assert.Equal(t, op.Name, apiOp.Name)
	assert.Equal(t, string(op.Type), apiOp.Type)
	assert.Equal(t, op.Description, apiOp.Description)
	assert.Equal(t, op.ReceiptURL, apiOp.ReceiptURL)
	assert.Equal(t, op.CreatedAt, apiOp.Date)
}

func TestOperationsToApi(t *testing.T) {
	userID := 42
	createdAt := time.Now()

	ops := []models.Operation{
		{
			ID:          1,
			AccountID:   2,
			CategoryID:  3,
			Type:        models.OperationIncome,
			Description: "Salary",
			ReceiptURL:  "",
			Name:        "Work salary",
			Sum:         3000,
			CreatedAt:   createdAt,
		},
		{
			ID:          2,
			AccountID:   2,
			CategoryID:  4,
			Type:        models.OperationExpense,
			Description: "Groceries",
			ReceiptURL:  "",
			Name:        "Food",
			Sum:         120,
			CreatedAt:   createdAt.Add(time.Hour),
		},
	}

	apiOps := OperationsToApi(userID, ops)

	assert.Equal(t, userID, apiOps.UserID)
	assert.Len(t, apiOps.Operations, len(ops))

	for i, op := range ops {
		apiOp := apiOps.Operations[i]
		assert.Equal(t, op.ID, apiOp.ID)
		assert.Equal(t, op.AccountID, apiOp.AccountID)
		assert.Equal(t, op.CategoryID, apiOp.CategoryID)
		assert.Equal(t, op.Sum, apiOp.Sum)
		assert.Equal(t, op.Name, apiOp.Name)
		assert.Equal(t, string(op.Type), apiOp.Type)
		assert.Equal(t, op.Description, apiOp.Description)
		assert.Equal(t, op.ReceiptURL, apiOp.ReceiptURL)
		assert.Equal(t, op.CreatedAt, apiOp.Date)
	}
}
