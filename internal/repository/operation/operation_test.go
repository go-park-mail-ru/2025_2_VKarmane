package operation

import (
	"context"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	"github.com/stretchr/testify/assert"
)

func TestRepository_GetOperationsByAccount(t *testing.T) {
	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	accID1 := 1
	accID2 := 2
	ops := []OperationDB{
		{ID: 1, AccountFromID: &accID1, Status: models.OperationFinished},
		{ID: 2, AccountToID: &accID1, Status: models.OperationFinished},
		{ID: 3, AccountFromID: &accID2, Status: models.OperationFinished},
		{ID: 4, AccountFromID: &accID1, Status: models.OperationReverted},
	}

	repo := NewRepository(ops, fixedClock)

	result, err := repo.GetOperationsByAccount(context.Background(), 1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestRepository_GetOperationByID(t *testing.T) {
	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	accID := 1
	ops := []OperationDB{
		{ID: 1, AccountFromID: &accID, Status: models.OperationFinished},
		{ID: 2, AccountFromID: &accID, Status: models.OperationFinished},
	}

	repo := NewRepository(ops, fixedClock)

	result, err := repo.GetOperationByID(context.Background(), 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, result.ID)
}

func TestRepository_CreateOperation(t *testing.T) {
	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	repo := NewRepository([]OperationDB{}, fixedClock)

	op := models.Operation{
		AccountID: 1,
		Name:      "Test",
		Sum:       100,
		Type:      models.OperationExpense,
	}

	result, err := repo.CreateOperation(context.Background(), op)
	assert.NoError(t, err)
	assert.NotZero(t, result.ID)
	assert.Equal(t, "Test", result.Name)
}

func TestRepository_UpdateOperation(t *testing.T) {
	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	accID := 1
	ops := []OperationDB{
		{ID: 1, AccountFromID: &accID, Name: "Old", Sum: 100, Status: models.OperationFinished},
	}

	repo := NewRepository(ops, fixedClock)

	newName := "Updated"
	newSum := 200.0
	req := models.UpdateOperationRequest{
		Name: &newName,
		Sum:  &newSum,
	}

	result, err := repo.UpdateOperation(context.Background(), req, 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", result.Name)
	assert.Equal(t, 200.0, result.Sum)
}

func TestRepository_DeleteOperation(t *testing.T) {
	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	accID := 1
	ops := []OperationDB{
		{ID: 1, AccountFromID: &accID, Status: models.OperationFinished},
	}

	repo := NewRepository(ops, fixedClock)

	result, err := repo.DeleteOperation(context.Background(), 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, models.OperationReverted, result.Status)
}

func TestRepository_GetOperationByID_NotFound(t *testing.T) {
	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	accID := 1
	ops := []OperationDB{
		{ID: 1, AccountFromID: &accID, Status: models.OperationFinished},
	}

	repo := NewRepository(ops, fixedClock)

	result, err := repo.GetOperationByID(context.Background(), 1, 99)
	assert.NoError(t, err)
	assert.Zero(t, result.ID)
}

func TestRepository_GetOperationByID_WrongAccount(t *testing.T) {
	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	accID1 := 1
	accID2 := 2
	ops := []OperationDB{
		{ID: 1, AccountFromID: &accID1, Status: models.OperationFinished},
	}

	repo := NewRepository(ops, fixedClock)

	result, err := repo.GetOperationByID(context.Background(), accID2, 1)
	assert.NoError(t, err)
	assert.Zero(t, result.ID)
}

func TestRepository_GetOperationByID_Reverted(t *testing.T) {
	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	accID := 1
	ops := []OperationDB{
		{ID: 1, AccountFromID: &accID, Status: models.OperationReverted},
	}

	repo := NewRepository(ops, fixedClock)

	result, err := repo.GetOperationByID(context.Background(), 1, 1)
	assert.NoError(t, err)
	assert.Zero(t, result.ID)
}

func TestRepository_GetOperationsByAccount_Empty(t *testing.T) {
	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	repo := NewRepository([]OperationDB{}, fixedClock)

	result, err := repo.GetOperationsByAccount(context.Background(), 1)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestRepository_GetOperationsByAccount_AccountToID(t *testing.T) {
	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	accID := 1
	ops := []OperationDB{
		{ID: 1, AccountToID: &accID, Status: models.OperationFinished},
	}

	repo := NewRepository(ops, fixedClock)

	result, err := repo.GetOperationsByAccount(context.Background(), 1)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 1, result[0].ID)
}

func TestRepository_GetOperationsByAccount_ExcludeReverted(t *testing.T) {
	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	accID := 1
	ops := []OperationDB{
		{ID: 1, AccountFromID: &accID, Status: models.OperationFinished},
		{ID: 2, AccountFromID: &accID, Status: models.OperationReverted},
		{ID: 3, AccountFromID: &accID, Status: models.OperationFinished},
	}

	repo := NewRepository(ops, fixedClock)

	result, err := repo.GetOperationsByAccount(context.Background(), 1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestRepository_UpdateOperation_NotFound(t *testing.T) {
	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	accID := 1
	ops := []OperationDB{
		{ID: 1, AccountFromID: &accID, Name: "Old", Sum: 100},
	}

	repo := NewRepository(ops, fixedClock)

	newName := "Updated"
	req := models.UpdateOperationRequest{
		Name: &newName,
	}

	result, err := repo.UpdateOperation(context.Background(), req, 1, 99)
	assert.NoError(t, err)
	assert.Zero(t, result.ID)
}

func TestRepository_UpdateOperation_AllFields(t *testing.T) {
	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	accID := 1
	catID := 5
	ops := []OperationDB{
		{ID: 1, AccountFromID: &accID, Name: "Old", Sum: 100, Description: "Old desc"},
	}

	repo := NewRepository(ops, fixedClock)

	newName := "Updated"
	newSum := 200.0
	newDesc := "New desc"
	newTime := time.Now()
	req := models.UpdateOperationRequest{
		Name:        &newName,
		Sum:         &newSum,
		Description: &newDesc,
		CategoryID:  &catID,
		CreatedAt:   &newTime,
	}

	result, err := repo.UpdateOperation(context.Background(), req, 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", result.Name)
	assert.Equal(t, 200.0, result.Sum)
}

func TestRepository_DeleteOperation_NotFound(t *testing.T) {
	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	accID := 1
	ops := []OperationDB{
		{ID: 1, AccountFromID: &accID, Status: models.OperationFinished},
	}

	repo := NewRepository(ops, fixedClock)

	result, err := repo.DeleteOperation(context.Background(), 1, 99)
	assert.NoError(t, err)
	assert.Zero(t, result.ID)
}

func TestRepository_CreateOperation_AllFields(t *testing.T) {
	fixedClock := clock.FixedClock{FixedTime: time.Now()}
	repo := NewRepository([]OperationDB{}, fixedClock)

	catID := 5
	currencyID := 1
	op := models.Operation{
		AccountID:   1,
		CategoryID:  catID,
		Name:        "Test Operation",
		Description: "Test Description",
		Sum:         150.75,
		Type:        models.OperationExpense,
		Status:      models.OperationFinished,
		CurrencyID:  currencyID,
		ReceiptURL:  "http://example.com/receipt",
	}

	result, err := repo.CreateOperation(context.Background(), op)
	assert.NoError(t, err)
	assert.NotZero(t, result.ID)
	assert.Equal(t, "Test Operation", result.Name)
	assert.Equal(t, "Test Description", result.Description)
	assert.Equal(t, 150.75, result.Sum)
	assert.Equal(t, models.OperationExpense, result.Type)
	assert.Equal(t, models.OperationFinished, result.Status)
}

