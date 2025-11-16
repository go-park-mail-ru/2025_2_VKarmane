package budget

import (
	bdgpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func BudgetToAPI(bdg *bdgpb.Budget) models.Budget {
	return models.Budget{
		ID:          int(bdg.Id),
		UserID:      int(bdg.UserId),
		CategoryID:  int(bdg.CategoryId),
		CurrencyID:  int(bdg.CurrencyId),
		Amount:      bdg.Sum,
		Description: bdg.Description,
		CreatedAt:   bdg.CreatedAt.AsTime(),
		UpdatedAt:   bdg.UpdatedAt.AsTime(),
		PeriodStart: bdg.PeriodStart.AsTime(),
		PeriodEnd:   bdg.PeriodEnd.AsTime(),
	}
}


func BudgetsToAPI(userID int, bdgs *bdgpb.ListBudgetsResponse) []models.Budget {
	res := make([]models.Budget, 0, len(bdgs.Budgets))
	for _, b := range bdgs.Budgets {
		res = append(res, BudgetToAPI(b))
	}
	return res
}

func IDsToBudgetRequest(budgetID, userID int) *bdgpb.BudgetRequest {
	return &bdgpb.BudgetRequest{
		UserID: int32(userID),
		BudgetID: int32(budgetID),
	}
}

func ModelCreateReqtoProtoReq(req models.CreateBudgetRequest, userID int) *bdgpb.CreateBudgetRequest {
	return &bdgpb.CreateBudgetRequest{
		UserID: int32(userID),
		CategoryId: int32(req.CategoryID),
		Sum: req.Amount,
		Description: req.Description,
		CreatedAt: timestamppb.New(req.CreatedAt),
		PeriodStart: timestamppb.New(req.PeriodStart),
		PeriodEnd: timestamppb.New(req.PeriodEnd),
	}
}

func ModelUpdateReqtoProtoReq(req models.UpdatedBudgetRequest, budgetID, userID int) *bdgpb.UpdateBudgetRequest {
	return &bdgpb.UpdateBudgetRequest{
		UserID: int32(userID),
		BudgetID: int32(budgetID),
		Sum: req.Amount,
		Description: req.Description,
		PeriodStart: timestamppb.New(*req.PeriodStart),
		PeriodEnd: timestamppb.New(*req.PeriodEnd),
	}
}