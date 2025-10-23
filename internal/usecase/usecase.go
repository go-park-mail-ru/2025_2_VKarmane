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
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

type UseCase struct {
	service   *service.Service
	BalanceUC *balance.UseCase
	OpUC      *operation.UseCase
	BudgetUC  *budget.UseCase
	AuthUC    *authUC.UseCase
}

func NewUseCase(service *service.Service, store *repository.Store, jwtSecret string) *UseCase {
	realClock := clock.RealClock{}
	userRepo := user.NewRepository(store.Users, realClock)
	authService := authService.NewService(userRepo, jwtSecret, realClock)
	authUC := authUC.NewUseCase(authService, realClock)

	return &UseCase{
		service:   service,
		BalanceUC: balance.NewUseCase(store, realClock),
		BudgetUC:  budget.NewUseCase(store, realClock),
		OpUC:      operation.NewUseCase(store, realClock),
		AuthUC:    authUC,
	}
}
