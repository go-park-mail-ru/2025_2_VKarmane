package usecase

import (
	image "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/image/usecase"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service"
)

// Repository алиас для service.Repository

type UseCase struct {
	service *service.Service
	ImageUC *image.UseCase
}

func NewUseCase(service *service.Service, jwtSecret string) *UseCase {
	imageUC := image.NewUseCase(service.ImageUC)

	return &UseCase{
		service: service,
		ImageUC: imageUC,
	}
}
