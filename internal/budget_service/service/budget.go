package budget

import (
	"context"

	pkgerrors "github.com/pkg/errors"

	bdgerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/errors"
	bdgmodels "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/models"
	bdgpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

type Service struct {
	repo BudgetRepository
	clock clock.Clock
}

func NewService(repo BudgetRepository, clck clock.Clock) *Service {
	return &Service{
		repo:  repo,
		clock: clck,
	}
}

func (s *Service) CheckBudgetOwnership(ctx context.Context, budgetID, userID int) bool {
	budgets, err := s.repo.GetBudgetsByUser(ctx, userID)
	if err != nil {
		return false
	}

	if budgets == nil {
		return false
	}
	for _, budget := range budgets {
		if budget.ID == budgetID {
			return true
		}
	}
	return false
}

func (s *Service) GetBudgets(ctx context.Context, userID int) (*bdgpb.ListBudgetsResponse, error) {
	budgets, err := s.repo.GetBudgetsByUser(ctx, userID)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "Failed to get budgets for user")
	}

	// accounts, err := s.repo.GetAccountsByUser(ctx, userID)
	// if err != nil {
	// 	return []models.Budget{}, pkgerrors.Wrap(err, "Failed to get accounts for user")
	// }

	// if accounts == nil {
	// 	accounts = []models.Account{}
	// }

	for i := range budgets {
		// var actual float64
		// for _, account := range accounts {
		// 	ops, err := s.repo.GetOperationsByAccount(ctx, account.ID)
		// 	if err != nil {
		// 		return []models.Budget{}, pkgerrors.Wrap(err, "Failed to get budgets for user")
		// 	}
		// 	for _, op := range ops {
		// 		if op.CurrencyID != budgets[i].CurrencyID {
		// 			continue
		// 		}
		// 		if op.CreatedAt.Before(budgets[i].PeriodStart) || op.CreatedAt.After(budgets[i].PeriodEnd) {
		// 			continue
		// 		}

		// 		if op.Type == "expense" {
		// 			actual += op.Sum
		// 		}
		// 	}
		// } // пока нет сервиса
		budgets[i].Actual = 123
	}

	return ModelListToProto(budgets), nil
}

func (s *Service) GetBudgetByID(ctx context.Context, budgetID, userID int) (*bdgpb.Budget, error) {
	if !s.CheckBudgetOwnership(ctx, budgetID, userID) {
		return nil, bdgerrors.ErrForbidden
	}
	budgets, err := s.repo.GetBudgetsByUser(ctx, userID)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "budget.GetBudgetByID: failed to get budgets")
	}

	for _, budget := range budgets {
		if budget.ID == budgetID {
			return ModelBudgetToProto(budget), nil
		}
	}

	return nil, bdgerrors.ErrBudgetNotFound
}

func (s *Service) CreateBudget(ctx context.Context, req bdgmodels.CreateBudgetRequest, userID int) (*bdgpb.Budget, error) {
	budget := CreateRequestToModel(req, userID)
	createdBgt, err := s.repo.CreateBudget(ctx, budget)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "Failed to create budget")
	}

	return ModelBudgetToProto(createdBgt), nil
}

func (s *Service) UpdateBudget(ctx context.Context, req bdgmodels.UpdatedBudgetRequest) (*bdgpb.Budget, error) {
	if !s.CheckBudgetOwnership(ctx, req.BudgetID, req.UserID) {
		return nil, bdgerrors.ErrForbidden
	}

	updatedBgt, err := s.repo.UpdateBudget(ctx, req)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "Failed to update budget")
	}

	return ModelBudgetToProto(updatedBgt), nil
}

func (s *Service) DeleteBudget(ctx context.Context, userID, budgetID int) (*bdgpb.Budget, error) {
	if !s.CheckBudgetOwnership(ctx, budgetID, userID) {
		return nil, bdgerrors.ErrForbidden
	}

	deletedBgt, err := s.repo.DeleteBudget(ctx, budgetID)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "Failed to delete budget")
	}

	return ModelBudgetToProto(deletedBgt), nil
}
