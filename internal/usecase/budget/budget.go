package budget

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/budget"
)

type UseCase struct {
	budgetSvc *budget.Service
}

func NewUseCase(store *repository.Store) *UseCase {
	return &UseCase{
		budgetSvc: budget.NewService(store),
	}
}

func (uc *UseCase) GetBudgetsForUser(userID int) ([]models.Budget, error) {
	budgetsData, err := uc.budgetSvc.GetBudgetsForUser(userID)
	if err != nil {
		return nil, err
	}

	return budgetsData, nil
}

func (uc *UseCase) GetBudgetByID(userID, budgetID int) (models.Budget, error) {
	budgetsData, err := uc.budgetSvc.GetBudgetsForUser(userID)
	if err != nil {
		return models.Budget{}, err
	}

	for _, budget := range budgetsData {
		if budget.ID == budgetID {
			return budget, nil
		}
	}

	return models.Budget{}, models.ErrBudgetNotFound
}
