package service

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/auth"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/balance"
	budgetService "github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/budget"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/category"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/operation"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/profile"
)

type Service struct {
	AuthUC     auth.AuthService
	BalanceUC  balance.BalanceService
	BudgetUC   budgetService.BudgetService
	OpUC       operation.OperationService
	CategoryUC category.CategoryService
	ProfileUC  profile.ProfileService
}

func NewService(store repository.Repository, jwtSecret string) *Service {
	postgresStore := store.(*repository.PostgresStore)

	userRepoAdapter := auth.NewPostgresUserRepositoryAdapter(postgresStore.UserRepo)
	authService := auth.NewService(userRepoAdapter, jwtSecret)

	accountRepoAdapter := balance.NewPostgresAccountRepositoryAdapter(postgresStore)
	budgetRepoAdapter := budgetService.NewPostgresBudgetRepositoryAdapter(postgresStore)
	operationRepoAdapter := budgetService.NewPostgresOperationRepositoryAdapter(postgresStore)
	operationAccountRepoAdapter := operation.NewPostgresAccountRepositoryAdapter(postgresStore)
	operationOperationRepoAdapter := operation.NewPostgresOperationRepositoryAdapter(postgresStore)
	categoryRepoAdapter := category.NewPostgresCategoryRepositoryAdapter(store)
	profileRepoAdapter := profile.NewPostgresProfileRepositoryAdapter(store)

	balanceService := balance.NewService(accountRepoAdapter)
	budgetService := budgetService.NewService(budgetRepoAdapter, accountRepoAdapter, operationRepoAdapter)
	opService := operation.NewService(operationAccountRepoAdapter, operationOperationRepoAdapter)
	categoryService := category.NewService(categoryRepoAdapter)
	profileService := profile.NewService(profileRepoAdapter)

	return &Service{
		AuthUC:     authService,
		BalanceUC:  balanceService,
		BudgetUC:   budgetService,
		OpUC:       opService,
		CategoryUC: categoryService,
		ProfileUC:  profileService,
	}
}
