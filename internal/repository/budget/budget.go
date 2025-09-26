package budget

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type Repository struct {
	budgets []BudgetDB
}

func NewRepository(budgets []BudgetDB) *Repository {
	return &Repository{
		budgets: budgets,
	}
}

func (r *Repository) GetBudgetsByUser(userID int) []models.Budget {
	out := make([]models.Budget, 0)

	for _, b := range r.budgets {
		if b.UserID == userID {
			out = append(out, BudgetDBToModel(b))
		}
	}

	return out
}
