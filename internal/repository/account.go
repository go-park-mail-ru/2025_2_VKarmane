package repository

import "github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"

func (s *Store) GetAccountsByUser(userID int) []models.Account {
	out := make([]models.Account, 0)
	accountsIDs := make(map[int]struct{})

	for _, ua := range s.UserAccounts {
		if ua.UserID == userID {
			accountsIDs[ua.AccountID] = struct{}{}
		}
	}

	for _, a := range s.Accounts {
		if _, ok := accountsIDs[a.ID]; ok {
			out = append(out, a)
		}
	}

	return out
}
