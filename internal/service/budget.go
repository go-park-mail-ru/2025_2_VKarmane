package service

import "time"

type BudgetView struct {
	ID          int       `json:"budget_id"`
	UserID      int       `json:"user_id"`
	Amount      float64   `json:"amount"`
	Actual      float64   `json:"actual"`
	CurrencyID  int       `json:"currency_id"`
	Description string    `json:"description"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
}

func (s *Service) GetBudgetsForUser(userID int) ([]BudgetView, error) {
	budgets := s.store.GetBudgetsByUser(userID)
	accounts := s.store.GetAccountsByUser(userID)

	out := make([]BudgetView, 0, len(budgets))

	for _, budget := range budgets {
		var actual float64
		for _, account := range accounts {
			ops := s.store.GetOperationsByAccount(account.ID)
			for _, op := range ops {
				if op.CurrencyID != budget.CurrencyID {
					continue
				}
				if op.CreatedAt.Before(budget.PeriodStart) || op.CreatedAt.After(budget.PeriodEnd) {
					continue
				}

				if op.Type == "expense" {
					actual += op.Sum
				}
			}
		}
		out = append(out, BudgetView{
			ID:          budget.ID,
			UserID:      budget.UserID,
			Amount:      budget.Amount,
			Actual:      actual,
			CurrencyID:  budget.CurrencyID,
			Description: budget.Description,
			PeriodStart: budget.PeriodStart,
			PeriodEnd:   budget.PeriodEnd,
		})
	}

	return out, nil
}
