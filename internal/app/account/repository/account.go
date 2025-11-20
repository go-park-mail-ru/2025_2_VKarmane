package account

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

type Repository struct {
	accounts     []AccountDB
	userAccounts []UserAccountDB
	clock        clock.Clock
}

func NewRepository(accounts []AccountDB, userAccounts []UserAccountDB, clck clock.Clock) *Repository {
	return &Repository{
		accounts:     accounts,
		userAccounts: userAccounts,
		clock:        clck,
	}
}

func (r *Repository) GetAccountsByUser(ctx context.Context, userID int) ([]models.Account, error) {
	out := make([]models.Account, 0)
	accountsIDs := make(map[int]struct{})

	for _, ua := range r.userAccounts {
		if ua.UserID == userID {
			accountsIDs[ua.AccountID] = struct{}{}
		}
	}

	for _, a := range r.accounts {
		if _, ok := accountsIDs[a.ID]; ok {
			out = append(out, AccountDBToModel(a))
		}
	}

	return out, nil
}
