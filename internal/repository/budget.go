package repository

import "github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"

func (s *Store) GetBudgetsByUser(userID int) []models.Budget {
	out := make([]models.Budget, 0)

	for _, b := range s.Budget {
		if b.UserID == userID {
			out = append(out, b)
		}
	}

	return out
}
