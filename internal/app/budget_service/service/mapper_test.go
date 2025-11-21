package budget

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	bdgmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreateRequestToModel(t *testing.T) {
	req := bdgmodels.CreateBudgetRequest{
		CategoryID:  2,
		Amount:      5000,
		CreatedAt:   time.Now(),
		PeriodStart: time.Now(),
		PeriodEnd:   time.Now().AddDate(0, 1, 0),
	}
	userID := 42

	budget := CreateRequestToModel(req, userID)

	assert.Equal(t, userID, budget.UserID)
	assert.Equal(t, req.CategoryID, budget.CategoryID)
	assert.Equal(t, req.Amount, budget.Amount)
	assert.Equal(t, 0.0, budget.Actual)
	assert.Equal(t, 1, budget.CurrencyID)
	assert.Equal(t, req.CreatedAt, budget.CreatedAt)
	assert.Equal(t, req.PeriodStart, budget.PeriodStart)
	assert.Equal(t, req.PeriodEnd, budget.PeriodEnd)
}

func TestModelBudgetToProto(t *testing.T) {
	now := time.Now()
	bdg := bdgmodels.Budget{
		ID:          1,
		UserID:      42,
		CategoryID:  2,
		Amount:      5000,
		Actual:      1000,
		CurrencyID:  1,
		Description: "Test budget",
		CreatedAt:   now,
		UpdatedAt:   now,
		PeriodStart: now,
		PeriodEnd:   now.AddDate(0, 1, 0),
	}

	pb := ModelBudgetToProto(bdg)

	assert.Equal(t, int32(bdg.ID), pb.Id)
	assert.Equal(t, int32(bdg.UserID), pb.UserId)
	assert.Equal(t, int32(bdg.CategoryID), pb.CategoryId)
	assert.Equal(t, bdg.Amount, pb.Sum)
	assert.Equal(t, bdg.Actual, pb.Actual)
	assert.Equal(t, int32(bdg.CurrencyID), pb.CurrencyId)
	assert.Equal(t, bdg.Description, pb.Description)
	assert.Equal(t, timestamppb.New(bdg.CreatedAt), pb.CreatedAt)
	assert.Equal(t, timestamppb.New(bdg.UpdatedAt), pb.UpdatedAt)
	assert.Equal(t, timestamppb.New(bdg.PeriodStart), pb.PeriodStart)
	assert.Equal(t, timestamppb.New(bdg.PeriodEnd), pb.PeriodEnd)
}

func TestModelListToProto(t *testing.T) {
	now := time.Now()
	budgets := []bdgmodels.Budget{
		{ID: 1, UserID: 1, CategoryID: 1, Amount: 100, Actual: 50, CurrencyID: 1, CreatedAt: now, UpdatedAt: now, PeriodStart: now, PeriodEnd: now},
		{ID: 2, UserID: 2, CategoryID: 2, Amount: 200, Actual: 150, CurrencyID: 1, CreatedAt: now, UpdatedAt: now, PeriodStart: now, PeriodEnd: now},
	}

	resp := ModelListToProto(budgets)

	assert.Len(t, resp.Budgets, 2)
	assert.Equal(t, int32(1), resp.Budgets[0].Id)
	assert.Equal(t, int32(2), resp.Budgets[1].Id)
	assert.Equal(t, budgets[0].Amount, resp.Budgets[0].Sum)
	assert.Equal(t, budgets[1].Amount, resp.Budgets[1].Sum)
}

func TestCreateRequestToModel_Empty(t *testing.T) {
	req := bdgmodels.CreateBudgetRequest{} // все поля нулевые
	userID := 0

	budget := CreateRequestToModel(req, userID)

	assert.Equal(t, 0, budget.UserID)
	assert.Equal(t, 0, budget.CategoryID)
	assert.Equal(t, 0.0, budget.Amount)
	assert.Equal(t, 0.0, budget.Actual)
	assert.Equal(t, 1, budget.CurrencyID) // по дефолту ставится 1
	assert.True(t, budget.CreatedAt.IsZero())
	assert.True(t, budget.PeriodStart.IsZero())
	assert.True(t, budget.PeriodEnd.IsZero())
}

func TestModelBudgetToProto_Empty(t *testing.T) {
	bdg := bdgmodels.Budget{} // нулевые значения

	pb := ModelBudgetToProto(bdg)

	assert.Equal(t, int32(0), pb.Id)
	assert.Equal(t, int32(0), pb.UserId)
	assert.Equal(t, int32(0), pb.CategoryId)
	assert.Equal(t, float64(0), pb.Sum)
	assert.Equal(t, float64(0), pb.Actual)
	assert.Equal(t, int32(0), pb.CurrencyId)
	assert.Equal(t, "", pb.Description)
	assert.True(t, pb.CreatedAt.AsTime().IsZero())
	assert.True(t, pb.UpdatedAt.AsTime().IsZero())
	assert.True(t, pb.PeriodStart.AsTime().IsZero())
	assert.True(t, pb.PeriodEnd.AsTime().IsZero())
}

func TestModelListToProto_Empty(t *testing.T) {
	resp := ModelListToProto([]bdgmodels.Budget{})

	assert.NotNil(t, resp)
	assert.Empty(t, resp.Budgets)
}
