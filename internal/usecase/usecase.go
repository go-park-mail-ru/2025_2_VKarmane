package usecase

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service"
	authUC "github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/auth"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/balance"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/budget"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/category"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/image"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/operation"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/profile"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

// Repository алиас для service.Repository
type Repository = service.Repository

type UseCase struct {
	service    *service.Service
	BalanceUC  *balance.UseCase
	OpUC       *operation.UseCase
	BudgetUC   *budget.UseCase
	AuthUC     *authUC.UseCase
	CategoryUC *category.UseCase
	ProfileUC  *profile.UseCase
	ImageUC    *image.UseCase
}

func NewUseCase(service *service.Service, store Repository, jwtSecret string) *UseCase {
	realClock := clock.RealClock{}
	authUC := authUC.NewUseCase(service.AuthUC, realClock)
	balanceUC := balance.NewUseCase(service.BalanceUC)
	budgetUC := budget.NewUseCase(service.BudgetUC)
	opUC := operation.NewUseCase(service.OpUC)
	categoryUC := category.NewUseCase(service.CategoryUC)
	profileUC := profile.NewUseCase(service.ProfileUC)
	imageUC := image.NewUseCase(service.ImageUC)

	return &UseCase{
		service:    service,
		BalanceUC:  balanceUC,
		BudgetUC:   budgetUC,
		OpUC:       opUC,
		AuthUC:     authUC,
		CategoryUC: categoryUC,
		ProfileUC:  profileUC,
		ImageUC:    imageUC,
	}
}
