package handlers

import (
	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers/auth"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers/balance"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers/budget"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers/category"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers/operation"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers/profile"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/proto"
	bdgpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/proto"
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
