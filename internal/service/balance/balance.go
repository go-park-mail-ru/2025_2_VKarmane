package balance

import (
	"context"
	"fmt"

	pkgerrors "github.com/pkg/errors"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	serviceerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/errors"
)

type Service struct {
	repo interface {
		GetAccountsByUser(ctx context.Context, userID int) ([]models.Account, error)
		CreateAccount(ctx context.Context, account models.Account, userID int) (models.Account, error)
		UpdateAccount(ctx context.Context, req models.UpdateAccountRequest, userID, accID int) (models.Account, error)
		DeleteAccount(ctx context.Context, userID, accID int) (models.Account, error)
	}
	clock clock.Clock
}

func NewService(repo interface {
	GetAccountsByUser(ctx context.Context, userID int) ([]models.Account, error)
	CreateAccount(ctx context.Context, account models.Account, userID int) (models.Account, error)
	UpdateAccount(ctx context.Context, req models.UpdateAccountRequest, userID, accID int) (models.Account, error)
	DeleteAccount(ctx context.Context, userID, accID int) (models.Account, error)
}, clck clock.Clock) *Service {
	return &Service{repo: repo, clock: clck}
}


func (s *Service) CheckAccountOwnership(ctx context.Context, accID int) bool {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok || userID == 0 {
		return false
	}
	accs, err := s.repo.GetAccountsByUser(ctx, userID)
	if err != nil {
		return false
	}

	if len(accs) == 0 {
		return false
	}
	for _, acc := range accs {
		if acc.ID == accID {
			return true
		}
	}
	return false
}

func (s *Service) GetBalanceForUser(ctx context.Context, userID int) ([]models.Account, error) {
	accounts, err := s.repo.GetAccountsByUser(ctx, userID)
	if err != nil {
		return []models.Account{}, pkgerrors.Wrap(err, "Failed to get balance for user")
	}

	return accounts, nil
}

func (s *Service) GetAccountByID(ctx context.Context, userID, accountID int) (models.Account, error) {
	accounts, err := s.repo.GetAccountsByUser(ctx, userID)
	if err != nil {
		return models.Account{}, pkgerrors.Wrap(err, "balance.GetAccountByID: failed to get accounts")
	}

	for _, account := range accounts {
		if account.ID == accountID {
			return account, nil
		}
	}

	return models.Account{}, fmt.Errorf("account not found: %s", models.ErrCodeAccountNotFound)
}

func (s *Service) CreateAccount(ctx context.Context, req models.CreateAccountRequest, userID int) (models.Account, error) {
	account := models.Account{
		Balance:    req.Balance,
		Type:       string(req.Type),
		CurrencyID: req.CurrencyID,
		CreatedAt:  s.clock.Now(),
		UpdatedAt:  s.clock.Now(),
	}


	createdAcc, err := s.repo.CreateAccount(ctx, account, userID)
	if err != nil {
		return models.Account{}, pkgerrors.Wrap(err, "failed to create account")
	}

	return createdAcc, nil
}


func (s *Service) UpdateAccount(ctx context.Context, req models.UpdateAccountRequest, userID, accID int) (models.Account, error) {
	if !s.CheckAccountOwnership(ctx, accID) {
		return models.Account{}, serviceerrors.ErrForbidden
	}

	updatedAcc, err := s.repo.UpdateAccount(ctx, req, userID, accID)
	if err != nil {
		return models.Account{}, pkgerrors.Wrap(err, "failed to update account")
	}

	return updatedAcc, nil
}

func (s *Service) DeleteAccount(ctx context.Context, userID, accID int) (models.Account, error) {
	if !s.CheckAccountOwnership(ctx, accID) {
		return models.Account{}, serviceerrors.ErrForbidden
	}

	deletedAcc, err := s.repo.DeleteAccount(ctx, userID, accID)
	if err != nil {
		return models.Account{}, pkgerrors.Wrap(err, "failed to delete account")
	}

	return deletedAcc, nil
}


	// userAcc := models.UserAccount{
	// 	UserID:    userID,
	// 	AccountID: createdAcc.ID,
	// 	CreatedAt: s.clock.Now(),
	// 	UpdatedAt: s.clock.Now(),
	// }
