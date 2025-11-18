package handlers

import (
	"github.com/gorilla/mux"

	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/proto"
	bdgpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/proto"
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

func (r *Registrator) RegisterAll(publicRouter *mux.Router, protectedRouter *mux.Router, uc *usecase.UseCase, log logger.Logger, authClient authpb.AuthServiceClient, budgetClient bdgpb.BudgetServiceClient) {
	auth.Register(publicRouter, protectedRouter, log, authClient)
	balance.Register(protectedRouter, uc.BalanceUC)
	budget.Register(protectedRouter, budgetClient)
	operation.Register(protectedRouter, uc.OpUC, uc.ImageUC)
	category.Register(protectedRouter, uc.CategoryUC, uc.ImageUC)
	profile.Register(protectedRouter, uc.ImageUC, authClient)
	image.Register(protectedRouter, uc.ImageUC)
}
