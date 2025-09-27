package balance

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	httputils "github.com/go-park-mail-ru/2025_2_VKarmane/pkg/http"
	"github.com/gorilla/mux"
)

type Handler struct {
	balanceUC BalanceUseCase
}

func NewHandler(balanceUC BalanceUseCase) *Handler {
	return &Handler{balanceUC: balanceUC}
}

func (h *Handler) getUserID(r *http.Request) (int, bool) {
	return middleware.GetUserIDFromContext(r.Context())
}

func (h *Handler) parseIDFromURL(r *http.Request, paramName string) (int, error) {
	vars := mux.Vars(r)
	idStr := vars[paramName]
	return strconv.Atoi(idStr)
}

func (h *Handler) GetListBalance(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r)
	if !ok {
		httputils.Error(w, r, "User not authenticated", http.StatusUnauthorized)
		return
	}

	accounts, err := h.balanceUC.GetBalanceForUser(r.Context(), userID)
	if err != nil {
		httputils.InternalError(w, r, "Failed to get balance for user")
		return
	}

	balanceDTO := AccountsToBalanceAPI(userID, accounts)
	httputils.Success(w, r, balanceDTO)
}

func (h *Handler) GetBalanceByAccountID(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromURL(r, "id")
	if err != nil {
		httputils.ValidationError(w, r, "Invalid account ID format", "id")
		return
	}

	userID, ok := h.getUserID(r)
	if !ok {
		httputils.Error(w, r, "User not authenticated", http.StatusUnauthorized)
		return
	}

	account, err := h.balanceUC.GetAccountByID(r.Context(), userID, id)
	if err != nil {
		httputils.NotFoundError(w, r, "Account not found")
		return
	}

	accountDTO := AccountToAPI(account)
	httputils.Success(w, r, accountDTO)
}
