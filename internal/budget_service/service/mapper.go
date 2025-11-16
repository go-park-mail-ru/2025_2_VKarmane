package budget

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	bdgmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/models"
	bdgpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/proto"
)

func CreateRequestToModel(req bdgmodels.CreateBudgetRequest, userID int) bdgmodels.Budget{
	return bdgmodels.Budget{
		UserID: userID,
		CategoryID: req.CategoryID,
		Amount: req.Amount,
		Actual: 0,
		CurrencyID: 1,
		CreatedAt: req.CreatedAt,
		PeriodStart: req.PeriodStart,
		PeriodEnd: req.PeriodEnd,
	}
}

func ModelBudgetToProto(bdg bdgmodels.Budget) *bdgpb.Budget {
	return &bdgpb.Budget{
		Id: int32(bdg.ID),
		UserId: int32(bdg.UserID),
		CategoryId: int32(bdg.CategoryID),
		Sum: bdg.Actual,
		Actual: bdg.Actual,
		CurrencyId: int32(bdg.CurrencyID),
		Description: bdg.Description,
		CreatedAt: timestamppb.New(bdg.CreatedAt),
		UpdatedAt: timestamppb.New(bdg.UpdatedAt),
		PeriodStart: timestamppb.New(bdg.PeriodStart),
		PeriodEnd: timestamppb.New(bdg.PeriodEnd),
	}
}

func ModelListToProto(bdg []bdgmodels.Budget) *bdgpb.ListBudgetsResponse {
	var budgets []*bdgpb.Budget
	for _, b := range bdg {
		budgets = append(budgets, ModelBudgetToProto(b))
	}
	return &bdgpb.ListBudgetsResponse {
		Budgets: budgets,
	}
}
