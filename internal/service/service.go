package service

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
	accountRepo "github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/account"
	budgetRepo "github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/budget"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/operation"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/user"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/auth"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/balance"
	budgetService "github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/budget"
)

type Service struct {
	AuthUC    auth.AuthService
	BalanceUC balance.BalanceService
	BudgetUC  budgetService.BudgetService
}

func NewService(store *repository.Store) *Service {
	userRepo := user.NewRepository(store.Users)
	accountRepo := accountRepo.NewRepository(store.Accounts, store.UserAccounts)
	budgetRepo := budgetRepo.NewRepository(store.Budget)
	operationRepo := operation.NewRepository(store.Operations)

	authService := auth.NewService(userRepo, "your-secret-key")
	balanceService := balance.NewService(accountRepo)
	budgetService := budgetService.NewService(budgetRepo, accountRepo, operationRepo)

	return &Service{
		AuthUC:    authService,
		BalanceUC: balanceService,
		BudgetUC:  budgetService,
	}
}
