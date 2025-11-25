package handlers

import (
	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/handlers/auth"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/handlers/profile"
	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/proto"
	budget "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/handlers"
	bdgpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/proto"
	balance "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/handlers/account"
	category "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/handlers/category"
	operation "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/handlers/operation"
	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

type Handler struct {
	balanceHandler  *balance.Handler
	budgetHandler   *budget.Handler
	authHandler     *auth.Handler
	opHandler       *operation.Handler
	categoryHandler *category.Handler
	profileHandler  *profile.Handler
	logger          logger.Logger
	registrator     *Registrator
}

func NewHandler(uc *usecase.UseCase, logger logger.Logger, authClient authpb.AuthServiceClient, budgetClient bdgpb.BudgetServiceClient, finClient finpb.FinanceServiceClient) *Handler {
	realClock := clock.RealClock{}
	return &Handler{
		balanceHandler:  balance.NewHandler(finClient, realClock),
		budgetHandler:   budget.NewHandler(realClock, budgetClient),
		authHandler:     auth.NewHandler(realClock, logger, authClient),
		opHandler:       operation.NewHandler(finClient, uc.ImageUC, realClock),
		categoryHandler: category.NewHandler(finClient, uc.ImageUC),
		profileHandler:  profile.NewHandler(uc.ImageUC, authClient),
		logger:          logger,
		registrator:     NewRegistrator(uc, logger),
	}
}

func (h *Handler) Register(publicRouter *mux.Router, protectedRouter *mux.Router, authCleint authpb.AuthServiceClient, budgetClient bdgpb.BudgetServiceClient, finClient finpb.FinanceServiceClient) {
	h.registrator.RegisterAll(publicRouter, protectedRouter, h.registrator.uc, h.logger, authCleint, budgetClient, finClient)
}
