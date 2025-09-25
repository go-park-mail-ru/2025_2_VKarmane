package handlers

import "github.com/go-park-mail-ru/2025_2_VKarmane/internal/service"

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}
