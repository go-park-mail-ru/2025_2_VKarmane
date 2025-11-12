package budget

import (
	"context"
	"fmt"

	pkgerrors "github.com/pkg/errors"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	serviceerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/errors"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

var ErrForbidden = pkgerrors.New("forbidden")

type Service struct {
	repo interface {
		GetBudgetsByUser(ctx context.Context, userID int) ([]models.Budget, error)
		GetAccountsByUser(ctx context.Context, userID int) ([]models.Account, error)
		GetOperationsByAccount(ctx context.Context, accountID int) ([]models.OperationInList, error)
		CreateBudget(ctx context.Context, budget models.Budget) (models.Budget, error)
		UpdateBudget(ctx context.Context, req models.UpdatedBudgetRequest, userID, budgetID int) (models.Budget, error)
		DeleteBudget(ctx context.Context, budgetID int) (models.Budget, error)
	}
	clock clock.Clock
}

func NewService(repo interface {
	GetBudgetsByUser(ctx context.Context, userID int) ([]models.Budget, error)
	GetAccountsByUser(ctx context.Context, userID int) ([]models.Account, error)
	GetOperationsByAccount(ctx context.Context, accountID int) ([]models.OperationInList, error)
	CreateBudget(ctx context.Context, budget models.Budget) (models.Budget, error)
	UpdateBudget(ctx context.Context, req models.UpdatedBudgetRequest, userID, budgetID int) (models.Budget, error)
	DeleteBudget(ctx context.Context, budgetID int) (models.Budget, error)
}, clck clock.Clock) *Service {
	return &Service{
		repo:  repo,
		clock: clck,
	}
}

func (s *Service) CheckBudgetOwnership(ctx context.Context, budgetID int) bool {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok || userID == 0 {
		return false
	}
	budgets, err := s.repo.GetBudgetsByUser(ctx, userID)
	if err != nil {
		return false
	}

	if len(budgets) == 0 {
		return false
	}
	for _, budget := range budgets {
		if budget.ID == budgetID {
			return true
		}
	}
	return false
}

func (s *Service) GetBudgetsForUser(ctx context.Context, userID int) ([]models.Budget, error) {
	budgets, err := s.repo.GetBudgetsByUser(ctx, userID)
	if err != nil {
		return []models.Budget{}, pkgerrors.Wrap(err, "Failed to get budgets for user")
	}
	if budgets == nil {
		budgets = []models.Budget{}
	}
	accounts, err := s.repo.GetAccountsByUser(ctx, userID)
	if err != nil {
		return []models.Budget{}, pkgerrors.Wrap(err, "Failed to get accounts for user")
	}

	if accounts == nil {
		accounts = []models.Account{}
	}

	for i := range budgets {
		var actual float64
		for _, account := range accounts {
			ops, err := s.repo.GetOperationsByAccount(ctx, account.ID)
			if err != nil {
				return []models.Budget{}, pkgerrors.Wrap(err, "Failed to get budgets for user")
			}
			for _, op := range ops {
				if op.CurrencyID != budgets[i].CurrencyID {
					continue
				}
				if op.CreatedAt.Before(budgets[i].PeriodStart) || op.CreatedAt.After(budgets[i].PeriodEnd) {
					continue
				}

				if op.Type == "expense" {
					actual += op.Sum
				}
			}
		}
		budgets[i].Actual = actual
	}

	return budgets, nil
}

func (s *Service) GetBudgetByID(ctx context.Context, userID, budgetID int) (models.Budget, error) {
	if !s.CheckBudgetOwnership(ctx, budgetID) {
		return models.Budget{}, serviceerrors.ErrForbidden
	}
	budgets, err := s.repo.GetBudgetsByUser(ctx, userID)
	if err != nil {
		return models.Budget{}, pkgerrors.Wrap(err, "budget.GetBudgetByID: failed to get budgets")
	}

	for _, budget := range budgets {
		if budget.ID == budgetID {
			return budget, nil
		}
	}

	return models.Budget{}, fmt.Errorf("budget not found: %s", models.ErrCodeBudgetNotFound)
}

func (s *Service) CreateBudget(ctx context.Context, req models.CreateBudgetRequest, userID int) (models.Budget, error) {
	budget := models.Budget{
		UserID:      userID,
		CategoryID:  req.CategoryID,
		Amount:      req.Amount,
		Actual:      0,
		CurrencyID:  1,
		Description: req.Description,
		CreatedAt:   s.clock.Now(),
		PeriodStart: req.PeriodStart,
		PeriodEnd:   req.PeriodEnd,
	}
	createdBgt, err := s.repo.CreateBudget(ctx, budget)
	if err != nil {
		return models.Budget{}, pkgerrors.Wrap(err, "Failed to create budget")
	}

	return createdBgt, nil
}

func (s *Service) UpdateBudget(ctx context.Context, req models.UpdatedBudgetRequest, userID, budgetID int) (models.Budget, error) {
	if !s.CheckBudgetOwnership(ctx, budgetID) {
		return models.Budget{}, serviceerrors.ErrForbidden
	}

	updatedBgt, err := s.repo.UpdateBudget(ctx, req, userID, budgetID)
	if err != nil {
		return models.Budget{}, pkgerrors.Wrap(err, "Failed to update budget")
	}

	return updatedBgt, nil
}

func (s *Service) DeleteBudget(ctx context.Context, userID, budgetID int) (models.Budget, error) {
	if !s.CheckBudgetOwnership(ctx, budgetID) {
		return models.Budget{}, serviceerrors.ErrForbidden
	}

	deletedBgt, err := s.repo.DeleteBudget(ctx, budgetID)
	if err != nil {
		return models.Budget{}, pkgerrors.Wrap(err, "Failed to delete budget")
	}

	return deletedBgt, nil
}
