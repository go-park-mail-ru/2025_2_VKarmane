package budget

import (
	"fmt"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/account"
	budgetRepo "github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/budget"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/operation"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/budget"
)

type UseCase struct {
	budgetSvc BudgetService
}

func NewUseCase(store *repository.Store) *UseCase {
	accountRepo := account.NewRepository(store.Accounts, store.UserAccounts)
	budgetRepo := budgetRepo.NewRepository(store.Budget)
	operationRepo := operation.NewRepository(store.Operations)
	budgetService := budget.NewService(budgetRepo, accountRepo, operationRepo)

	return &UseCase{
		budgetSvc: budgetService,
	}
}

func (uc *UseCase) GetBudgetsForUser(userID int) ([]models.Budget, error) {
	budgetsData, err := uc.budgetSvc.GetBudgetsForUser(userID)
	if err != nil {
		return nil, fmt.Errorf("budget.GetBudgetsForUser: %w", err)
	}

	return budgetsData, nil
}

func (uc *UseCase) GetBudgetByID(userID, budgetID int) (models.Budget, error) {
	budgetsData, err := uc.budgetSvc.GetBudgetsForUser(userID)
	if err != nil {
		return models.Budget{}, fmt.Errorf("budget.GetBudgetByID: %w", err)
	}

	for _, budget := range budgetsData {
		if budget.ID == budgetID {
			return budget, nil
		}
	}

	return models.Budget{}, fmt.Errorf("budget.GetBudgetByID: %s", models.ErrCodeBudgetNotFound)
}
