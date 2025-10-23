package service

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/auth"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/balance"
	budgetService "github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/budget"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/operation"
)

type Service struct {
	AuthUC    auth.AuthService
	BalanceUC balance.BalanceService
	BudgetUC  budgetService.BudgetService
	OpUC      operation.OperationService
}

func NewService(store repository.Repository, jwtSecret string) *Service {
	postgresStore := store.(*repository.PostgresStore)

	userRepoAdapter := auth.NewPostgresUserRepositoryAdapter(postgresStore.UserRepo)
	authService := auth.NewService(userRepoAdapter, jwtSecret)

	accountRepoAdapter := balance.NewPostgresAccountRepositoryAdapter(postgresStore)
	budgetRepoAdapter := budgetService.NewPostgresBudgetRepositoryAdapter(postgresStore)
	operationRepoAdapter := budgetService.NewPostgresOperationRepositoryAdapter(postgresStore)

	balanceService := balance.NewService(accountRepoAdapter)
	budgetService := budgetService.NewService(budgetRepoAdapter, accountRepoAdapter, operationRepoAdapter)
	opService := operation.NewService(accountRepoAdapter, operationRepoAdapter)

	return &Service{
		AuthUC:    authService,
		BalanceUC: balanceService,
		BudgetUC:  budgetService,
		OpUC:      opService,
	}
}
