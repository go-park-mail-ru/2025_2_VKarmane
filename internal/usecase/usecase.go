package usecase

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/balance"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/category"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/image"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/operation"
)

// Repository алиас для service.Repository
type Repository = service.Repository

type UseCase struct {
	service    *service.Service
	BalanceUC  *balance.UseCase
	OpUC       *operation.UseCase
	CategoryUC *category.UseCase
	ImageUC    *image.UseCase
}

func NewUseCase(service *service.Service, store Repository, jwtSecret string) *UseCase {
	balanceUC := balance.NewUseCase(service.BalanceUC)
	opUC := operation.NewUseCase(service.OpUC)
	categoryUC := category.NewUseCase(service.CategoryUC)
	imageUC := image.NewUseCase(service.ImageUC)

	return &UseCase{
		service:    service,
		BalanceUC:  balanceUC,
		OpUC:       opUC,
		CategoryUC: categoryUC,
		ImageUC:    imageUC,
	}
}
