package operation

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/stretchr/testify/require"
)

func sampleOperations() []OperationDB {
	return []OperationDB{
		{ID: 1, AccountID: 1, Name: "Op1", Status: models.OperationFinished, Sum: 100},
		{ID: 2, AccountID: 1, Name: "Op2", Status: models.OperationFinished, Sum: 200},
		{ID: 3, AccountID: 2, Name: "OtherAccOp", Status: models.OperationFinished, Sum: 300},
		{ID: 4, AccountID: 1, Name: "Reverted", Status: models.OperationReverted, Sum: 400},
	}
}

func TestGetOperationsByAccount(t *testing.T) {
	repo := NewRepository(sampleOperations())

	ops, _ := repo.GetOperationsByAccount(context.Background(), 1)
	require.Len(t, ops, 2, "should exclude reverted ops and match accountID")
	require.Equal(t, "Op1", ops[0].Name)
	require.Equal(t, "Op2", ops[1].Name)
}

func TestGetOperationByID(t *testing.T) {
	repo := NewRepository(sampleOperations())

	op, _ := repo.GetOperationByID(context.Background(), 1, 2)
	require.Equal(t, 2, op.ID)
	require.Equal(t, "Op2", op.Name)
	op, _ = repo.GetOperationByID(context.Background(), 1, 999)
	require.Zero(t, op.ID, "should return empty struct when not found")

	op, _ = repo.GetOperationByID(context.Background(), 1, 4)
	require.Zero(t, op.ID)
}

func TestCreateOperation(t *testing.T) {
	repo := NewRepository([]OperationDB{})

	op := models.Operation{
		AccountID:  1,
		CategoryID: 2,
		Type:       models.OperationExpense,
		Status:     models.OperationFinished,
		Name:       "NewOp",
		Sum:        999,
	}

	created, _ := repo.CreateOperation(context.Background(), op)

	require.Equal(t, 1, created.ID)
	require.Equal(t, "NewOp", created.Name)
	require.Equal(t, models.OperationFinished, created.Status)
	require.NotZero(t, created.CreatedAt)
	require.Len(t, repo.operations, 1)
	require.Equal(t, "NewOp", repo.operations[0].Name)
}

func TestUpdateOperation(t *testing.T) {
	repo := NewRepository([]OperationDB{
		{ID: 1, AccountID: 10, Name: "OldName", Sum: 100},
	})

	newName := "Updated"
	newSum := float64(200)
	req := models.UpdateOperationRequest{
		Name: &newName,
		Sum:  &newSum,
	}

	updated, _ := repo.UpdateOperation(context.Background(), req, 10, 1)
	require.Equal(t, "Updated", updated.Name)
	require.Equal(t, float64(200), updated.Sum)

	empty, _ := repo.UpdateOperation(context.Background(), req, 999, 1)
	require.Zero(t, empty.ID)
}

func TestDeleteOperation(t *testing.T) {
	repo := NewRepository([]OperationDB{
		{ID: 1, AccountID: 7, Name: "Active", Status: models.OperationFinished},
	})

	deleted, _ := repo.DeleteOperation(context.Background(), 7, 1)
	require.Equal(t, 1, deleted.ID)
	require.Equal(t, models.OperationReverted, deleted.Status)
	require.Equal(t, models.OperationReverted, repo.operations[0].Status)
	empty, _ := repo.DeleteOperation(context.Background(), 7, 99)
	require.Zero(t, empty.ID)
}

func TestCreateOperationAssignsIncrementalIDs(t *testing.T) {
	repo := NewRepository([]OperationDB{
		{ID: 1, AccountID: 1, Name: "Old"},
	})

	newOp := models.Operation{
		AccountID: 1,
		Name:      "Second",
		Status:    models.OperationFinished,
	}
	created, _ := repo.CreateOperation(context.Background(), newOp)

	require.Equal(t, 2, created.ID)
	require.Len(t, repo.operations, 2)
	require.Equal(t, "Second", repo.operations[1].Name)
}
