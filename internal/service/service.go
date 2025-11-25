package service

import (
	image "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/image/repository"
	imageservice "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/image/service"
)

type Service struct {
	ImageUC imageservice.ImageService
}

func NewService(jwtSecret string, imageStorage image.ImageStorage) *Service {
	imageService := imageservice.NewService(imageStorage)

	return &Service{
		ImageUC: imageService,
	}
}
