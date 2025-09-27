package handlers

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers/auth"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers/balance"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers/budget"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase"
)

type Handler struct {
	balanceHandler *balance.Handler
	budgetHandler  *budget.Handler
	authHandler    *auth.Handler
	logger         logger.Logger
}

func NewHandler(uc *usecase.UseCase, logger logger.Logger) *Handler {
	return &Handler{
		balanceHandler: balance.NewHandler(uc.BalanceUC),
		budgetHandler:  budget.NewHandler(uc.BudgetUC),
		authHandler:    auth.NewHandler(uc.AuthUC, logger),
		logger:         logger,
	}
}

func (h *Handler) GetListBalance(w http.ResponseWriter, r *http.Request) {
	h.balanceHandler.GetListBalance(w, r)
}

func (h *Handler) GetBalanceByAccountID(w http.ResponseWriter, r *http.Request) {
	h.balanceHandler.GetBalanceByAccountID(w, r)
}

func (h *Handler) GetListBudgets(w http.ResponseWriter, r *http.Request) {
	h.budgetHandler.GetListBudgets(w, r)
}

func (h *Handler) GetBudgetByID(w http.ResponseWriter, r *http.Request) {
	h.budgetHandler.GetBudgetByID(w, r)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	h.authHandler.Register(w, r)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	h.authHandler.Login(w, r)
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	h.authHandler.GetProfile(w, r)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	h.authHandler.Logout(w, r)
}
