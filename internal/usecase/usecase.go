package usecase

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service"
	authUC "github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/auth"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/balance"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/budget"
)

type UseCase struct {
	service   *service.Service
	BalanceUC balance.BalanceService
	BudgetUC  budget.BudgetService
	AuthUC    *authUC.UseCase
}

func NewUseCase(service *service.Service, store repository.Repository, jwtSecret string) *UseCase {
	authUC := authUC.NewUseCase(service.AuthUC)
	balanceUC := balance.NewUseCase(store)
	budgetUC := budget.NewUseCase(store)

	return &UseCase{
		service:   service,
		BalanceUC: balanceUC,
		BudgetUC:  budgetUC,
		AuthUC:    authUC,
	}
}
