package usecase

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/user"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service"
	authService "github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/auth"
	authUC "github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/auth"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/balance"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/budget"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/operation"
)

type UseCase struct {
	service   *service.Service
	BalanceUC *balance.UseCase
	OpUC      *operation.UseCase
	BudgetUC  *budget.UseCase
	AuthUC    *authUC.UseCase
}

func NewUseCase(service *service.Service, store *repository.Store, jwtSecret string) *UseCase {
	userRepo := user.NewRepository(store.Users)
	authService := authService.NewService(userRepo, jwtSecret)
	authUC := authUC.NewUseCase(authService)

	return &UseCase{
		service:   service,
		BalanceUC: balance.NewUseCase(store),
		BudgetUC:  budget.NewUseCase(store),
		OpUC:      operation.NewUseCase(store),
		AuthUC:    authUC,
	}
}
