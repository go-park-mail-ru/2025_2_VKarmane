package usecase

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/balance"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/budget"
)

type UseCase struct {
	service   *service.Service
	BalanceUC *balance.UseCase
	BudgetUC  *budget.UseCase
}

func NewUseCase(service *service.Service, store *repository.Store) *UseCase {
	return &UseCase{
		service:   service,
		BalanceUC: balance.NewUseCase(store),
		BudgetUC:  budget.NewUseCase(store),
	}
}
