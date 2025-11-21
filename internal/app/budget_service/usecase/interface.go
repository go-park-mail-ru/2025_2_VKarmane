package budget

import (
	"context"

	budg "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/models"
	budgetpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/proto"
)

type BudgetService interface {
	CreateBudget(context.Context, budg.CreateBudgetRequest, int) (*budgetpb.Budget, error)
	GetBudgetByID(ctx context.Context, budgetID, userID int) (*budgetpb.Budget, error)
	GetBudgets(context.Context, int) (*budgetpb.ListBudgetsResponse, error)
	UpdateBudget(context.Context, budg.UpdatedBudgetRequest) (*budgetpb.Budget, error)
	DeleteBudget(ctx context.Context, budgetID, userID int) (*budgetpb.Budget, error)
}
