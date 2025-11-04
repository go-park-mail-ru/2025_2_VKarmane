package operation

import (
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestOperationToAPI(t *testing.T) {
	now := time.Now()
	op := models.Operation{
		ID:           1,
		AccountID:    10,
		CategoryID:   5,
		CategoryName: "Food",
		Type:         "expense",
		Status:       "finished",
		Name:         "Lunch",
		Sum:          250.50,
		CurrencyID:   1,
		CreatedAt:    now,
		Date:         now,
	}

	apiOp := OperationToAPI(op)

	assert.Equal(t, 1, apiOp.ID)
	assert.Equal(t, 10, apiOp.AccountID)
	assert.Equal(t, 5, apiOp.CategoryID)
	assert.Equal(t, "expense", apiOp.Type)
	assert.Equal(t, "finished", apiOp.Status)
	assert.Equal(t, "Lunch", apiOp.Name)
	assert.Equal(t, 250.50, apiOp.Sum)
}

func TestOperationsToAPI(t *testing.T) {
	now := time.Now()
	ops := []models.Operation{
		{ID: 1, Name: "Op1", Sum: 100, CreatedAt: now},
		{ID: 2, Name: "Op2", Sum: 200, CreatedAt: now},
	}

	apiOps := OperationsToAPI(ops)

	assert.Len(t, apiOps.Operations, 2)
	assert.Equal(t, 1, apiOps.Operations[0].ID)
	assert.Equal(t, 2, apiOps.Operations[1].ID)
	assert.Equal(t, 100.0, apiOps.Operations[0].Sum)
	assert.Equal(t, 200.0, apiOps.Operations[1].Sum)
}

