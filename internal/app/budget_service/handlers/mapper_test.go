package budget

import (
	"testing"
	"time"

	bdgpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestBudgetToAPI(t *testing.T) {
	fixedTime := time.Date(2030, 5, 20, 10, 0, 0, 0, time.UTC)
	fixed := clock.FixedClock{FixedTime: fixedTime}
	createdAt := fixed.FixedTime
	updatedAt := createdAt.Add(time.Hour)
	periodStart := createdAt.Add(-time.Hour * 24)
	periodEnd := createdAt.Add(time.Hour * 24)

	protoBdg := &bdgpb.Budget{
		Id:          1,
		UserId:      10,
		CurrencyId:  30,
		Sum:         100.50,
		Description: "Test budget",
		CreatedAt:   timestamppb.New(createdAt),
		UpdatedAt:   timestamppb.New(updatedAt),
		PeriodStart: timestamppb.New(periodStart),
		PeriodEnd:   timestamppb.New(periodEnd),
	}

	apiBdg := BudgetToAPI(protoBdg)

	assert.Equal(t, 1, apiBdg.ID)
	assert.Equal(t, 10, apiBdg.UserID)
	assert.Equal(t, 30, apiBdg.CurrencyID)
	assert.Equal(t, 100.50, apiBdg.Amount)
	assert.Equal(t, "Test budget", apiBdg.Description)
	assert.Equal(t, createdAt, apiBdg.CreatedAt)
	assert.Equal(t, updatedAt, apiBdg.UpdatedAt)
	assert.Equal(t, periodStart, apiBdg.PeriodStart)
	assert.Equal(t, periodEnd, apiBdg.PeriodEnd)
}

func TestBudgetsToAPI(t *testing.T) {
	budgetsProto := &bdgpb.ListBudgetsResponse{
		Budgets: []*bdgpb.Budget{
			{Id: 1, UserId: 10, Sum: 50},
			{Id: 2, UserId: 10, Sum: 75},
		},
	}

	result := BudgetsToAPI(10, budgetsProto)

	assert.Len(t, result, 2)
	assert.Equal(t, 1, result[0].ID)
	assert.Equal(t, 50.0, result[0].Amount)
	assert.Equal(t, 2, result[1].ID)
	assert.Equal(t, 75.0, result[1].Amount)
}

func TestIDsToBudgetRequest(t *testing.T) {
	req := IDsToBudgetRequest(5, 10)
	assert.Equal(t, int32(5), req.BudgetID)
	assert.Equal(t, int32(10), req.UserID)
}

func TestModelCreateReqtoProtoReq(t *testing.T) {
	fixedTime := time.Date(2030, 5, 20, 10, 0, 0, 0, time.UTC)
	fixed := clock.FixedClock{FixedTime: fixedTime}
	now := fixed.FixedTime
	req := models.CreateBudgetRequest{
		Amount:      100,
		Description: "Test",
		CreatedAt:   now,
		PeriodStart: now,
		PeriodEnd:   now.Add(time.Hour),
	}

	protoReq := ModelCreateReqtoProtoReq(req, 42)
	assert.Equal(t, int32(42), protoReq.UserID)
	assert.Equal(t, 100.0, protoReq.Sum)
	assert.Equal(t, "Test", protoReq.Description)
	assert.Equal(t, now, protoReq.CreatedAt.AsTime())
	assert.Equal(t, now, protoReq.PeriodStart.AsTime())
	assert.Equal(t, now.Add(time.Hour), protoReq.PeriodEnd.AsTime())
}

func TestModelUpdateReqtoProtoReq(t *testing.T) {
	fixedTime := time.Date(2030, 5, 20, 10, 0, 0, 0, time.UTC)
	fixed := clock.FixedClock{FixedTime: fixedTime}
	now := fixed.FixedTime
	periodStart := now
	periodEnd := now.Add(time.Hour)
	amount := 200.0

	descr := "test"
	req := models.UpdatedBudgetRequest{
		Amount:      &amount,
		Description: &descr,
		PeriodStart: &periodStart,
		PeriodEnd:   &periodEnd,
	}

	protoReq := ModelUpdateReqtoProtoReq(req, 7, 42)
	assert.Equal(t, int32(42), protoReq.UserID)
	assert.Equal(t, int32(7), protoReq.BudgetID)
	assert.Equal(t, 200.0, *protoReq.Sum)
	assert.Equal(t, "test", *protoReq.Description)
	assert.Equal(t, periodStart, protoReq.PeriodStart.AsTime())
	assert.Equal(t, periodEnd, protoReq.PeriodEnd.AsTime())
}
