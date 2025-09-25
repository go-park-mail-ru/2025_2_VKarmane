package service

import "github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"

type Service struct {
	store *repository.Store
}

func NewService(store *repository.Store) *Service {
	return &Service{store: store}
}
