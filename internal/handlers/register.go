package handlers

import (
	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers/auth"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers/balance"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers/budget"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers/category"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers/image"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers/operation"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers/profile"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase"
)

type Registrator struct {
	uc  *usecase.UseCase
	log logger.Logger
}

func NewRegistrator(uc *usecase.UseCase, log logger.Logger) *Registrator {
	return &Registrator{
		uc:  uc,
		log: log,
	}

}

func (r *Registrator) RegisterAll(publicRouter *mux.Router, protectedRouter *mux.Router, uc *usecase.UseCase, log logger.Logger) {
	auth.Register(publicRouter, protectedRouter, uc.AuthUC, log)
	balance.Register(protectedRouter, uc.BalanceUC)
	budget.Register(protectedRouter, uc.BudgetUC)
	operation.Register(protectedRouter, uc.OpUC, uc.ImageUC)
	category.Register(protectedRouter, uc.CategoryUC, uc.ImageUC)
	profile.Register(protectedRouter, uc.ProfileUC, uc.ImageUC)
	image.Register(protectedRouter, uc.ImageUC)
}
