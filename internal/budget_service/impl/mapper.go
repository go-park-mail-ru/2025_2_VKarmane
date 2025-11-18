package budget

import (
	"time"

	budgmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/models"
	budgpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/proto"
)

func ProtoBudgetToApi(bdg *budgpb.Budget) budgmodels.Budget {
	return budgmodels.Budget{
		ID:          int(bdg.Id),
		UserID:      int(bdg.UserId),
		CategoryID:  int(bdg.CategoryId),
		Amount:      bdg.Sum,
		Actual:      bdg.Actual,
		CurrencyID:  int(bdg.CurrencyId),
		Description: bdg.Description,
		CreatedAt:   bdg.CreatedAt.AsTime(),
		UpdatedAt:   bdg.UpdatedAt.AsTime(),
		PeriodStart: bdg.PeriodStart.AsTime(),
		PeriodEnd:   bdg.PeriodEnd.AsTime(),
	}
}

func ProtoCreateRequestToModel(req *budgpb.CreateBudgetRequest) (budgmodels.CreateBudgetRequest, int) {
	return budgmodels.CreateBudgetRequest{
		CategoryID:  int(req.CategoryId),
		Amount:      req.Sum,
		Description: req.Description,
		CreatedAt:   req.CreatedAt.AsTime(),
		PeriodStart: req.PeriodEnd.AsTime(),
		PeriodEnd:   req.PeriodEnd.AsTime(),
	}, int(req.UserID)
}

func ProtoUpdateRequestToModel(req *budgpb.UpdateBudgetRequest) budgmodels.UpdatedBudgetRequest {
	var periodStart *time.Time
	if req.PeriodStart != nil {
		t := req.PeriodStart.AsTime()
		periodStart = &t
	}
	var periodEnd *time.Time
	if req.PeriodEnd != nil {
		t := req.PeriodEnd.AsTime()
		periodEnd = &t
	}
	return budgmodels.UpdatedBudgetRequest{
		UserID:      int(req.UserID),
		BudgetID:    int(req.BudgetID),
		Amount:      req.Sum,
		Description: req.Description,
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
	}
}

func ProtoIDToInt(id *budgpb.UserID) int {
	return int(id.UserID)
}

func ProtoBudgetReqToInts(req *budgpb.BudgetRequest) (int, int) {
	return int(req.UserID), int(req.BudgetID)
}
