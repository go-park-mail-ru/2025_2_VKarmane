package service

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
	accountRepo "github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/account"
	budgetRepo "github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/budget"
	opRepo "github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/operation"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/user"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/auth"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/balance"
	budgetService "github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/budget"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/operation"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

type Service struct {
	AuthUC    auth.AuthService
	BalanceUC balance.BalanceService
	BudgetUC  budgetService.BudgetService
	OpUC      operation.OperationService
}

func NewService(store *repository.Store) *Service {
	realClock := clock.RealClock{}
	userRepo := user.NewRepository(store.Users, realClock)
	accountRepo := accountRepo.NewRepository(store.Accounts, store.UserAccounts, realClock)
	budgetRepo := budgetRepo.NewRepository(store.Budget, realClock)
	operationRepo := opRepo.NewRepository(store.Operations, realClock)

	authService := auth.NewService(userRepo, "your-secret-key", realClock)
	balanceService := balance.NewService(accountRepo, realClock)
	budgetService := budgetService.NewService(budgetRepo, accountRepo, operationRepo, realClock)
	opService := operation.NewService(accountRepo, operationRepo, realClock)

	return &Service{
		AuthUC:    authService,
		BalanceUC: balanceService,
		BudgetUC:  budgetService,
		OpUC:      opService,
	}
}
