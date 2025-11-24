package handlers

import (
	"github.com/gorilla/mux"

	balance "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/account/handlers"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/handlers/auth"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/handlers/profile"
	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/proto"
	budget "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/handlers"
	bdgpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/proto"
	category "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/handlers/category"
	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
	image "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/image/handlers"
	operation "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/operations/handlers"
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

func (r *Registrator) RegisterAll(publicRouter *mux.Router, protectedRouter *mux.Router, uc *usecase.UseCase, log logger.Logger, authClient authpb.AuthServiceClient, budgetClient bdgpb.BudgetServiceClient, finClient finpb.FinanceServiceClient) {
	auth.Register(publicRouter, protectedRouter, log, authClient)
	balance.Register(protectedRouter, uc.BalanceUC)
	budget.Register(protectedRouter, budgetClient)
	operation.Register(protectedRouter, uc.OpUC, uc.ImageUC)
	category.Register(protectedRouter, finClient, uc.ImageUC)
	profile.Register(protectedRouter, uc.ImageUC, authClient)
	image.Register(protectedRouter, uc.ImageUC)
}
