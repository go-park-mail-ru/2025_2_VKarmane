package handlers

import (
	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers/auth"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers/balance"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers/budget"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers/operation"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

type Handler struct {
	balanceHandler *balance.Handler
	budgetHandler  *budget.Handler
	authHandler    *auth.Handler
	opHandler      *operation.Handler
	logger         logger.Logger
	registrator    *Registrator
}

func NewHandler(uc *usecase.UseCase, logger logger.Logger) *Handler {
	realClock := clock.RealClock{}
	return &Handler{
		balanceHandler: balance.NewHandler(uc.BalanceUC, realClock),
		budgetHandler:  budget.NewHandler(uc.BudgetUC, realClock),
		authHandler:    auth.NewHandler(uc.AuthUC, realClock, logger),
		opHandler:      operation.NewHandler(uc.OpUC, realClock),
		logger:         logger,
		registrator:    NewRegistrator(uc, logger),
	}
}

func (h *Handler) Register(publicRouter *mux.Router, protectedRouter *mux.Router) {
	h.registrator.RegisterAll(publicRouter, protectedRouter, h.registrator.uc, h.logger)
}
