package repository

import "github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"

func (s *Store) GetOperationsByAccount(accountID int) []models.Operation {
	out := make([]models.Operation, 0)

	for _, o := range s.Operations {
		if o.AccountID == accountID {
			out = append(out, o)
		}
	}

	return out
}
