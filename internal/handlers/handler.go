package handlers

import (
	"github.com/gorilla/mux"

	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/proto"
	bdgpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/handlers/auth"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/account/handlers"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/handlers"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/category/handlers"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/operations/handlers"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/handlers/profile"
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

func NewHandler(uc *usecase.UseCase, logger logger.Logger, authClient authpb.AuthServiceClient, budgetClient bdgpb.BudgetServiceClient) *Handler {
	realClock := clock.RealClock{}
	return &Handler{
		balanceHandler:  balance.NewHandler(uc.BalanceUC, realClock),
		budgetHandler:   budget.NewHandler(realClock, budgetClient),
		authHandler:     auth.NewHandler(realClock, logger, authClient),
		opHandler:       operation.NewHandler(uc.OpUC, uc.ImageUC, realClock),
		categoryHandler: category.NewHandler(uc.CategoryUC, uc.ImageUC),
		profileHandler:  profile.NewHandler(uc.ImageUC, authClient),
		logger:          logger,
		registrator:     NewRegistrator(uc, logger),
	}
}

func (h *Handler) Register(publicRouter *mux.Router, protectedRouter *mux.Router, authCleint authpb.AuthServiceClient, budgetClient bdgpb.BudgetServiceClient) {
	h.registrator.RegisterAll(publicRouter, protectedRouter, h.registrator.uc, h.logger, authCleint, budgetClient)
}
