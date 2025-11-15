package support

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	httputils "github.com/go-park-mail-ru/2025_2_VKarmane/pkg/http"
)

type Handler struct {
	supUC SupportUseCase
}

func NewHandler(supUC SupportUseCase) *Handler {
	return &Handler{supUC: supUC}
}

func (h *Handler) parseUserID(r *http.Request) (int, bool) {
	return middleware.GetUserIDFromContext(r.Context())
}

func (h *Handler) parseIDFromURL(r *http.Request, param string) (int, error) {
	val := middleware.GetParam(r, param)
	return strconv.Atoi(val)
}

func (h *Handler) CreateSupportRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := h.parseUserID(r)
	if !ok {
		httputils.UnauthorizedError(w, r, "Требуется авторизация", models.ErrCodeUnauthorized)
		return
	}

	var input CreateRequestInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httputils.ValidationError(w, r, "Некорректный формат данных", "body")
		return
	}

	ticket, err := h.supUC.CreateSupportRequest(ctx, userID, models.CategoryContacting(input.Category), input.Message)
	if err != nil {
		httputils.InternalError(w, r, err.Error())
		return
	}

	resp := CreateRequestOutput{
		ID:       ticket.ID,
		Status:   string(ticket.StatusRequest),
		Category: string(ticket.CategoryRequest),
		Message:  ticket.Message,
	}

	httputils.Created(w, r, resp)
}

func (h *Handler) GetUserSupportRequests(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := h.parseUserID(r)
	if !ok {
		httputils.UnauthorizedError(w, r, "Требуется авторизация", models.ErrCodeUnauthorized)
		return
	}

	tickets, err := h.supUC.GetUserSupportRequests(ctx, userID)
	if err != nil {
		httputils.InternalError(w, r, err.Error())
		return
	}

	out := UserRequestsOutput{}
	for _, t := range tickets {
		out.Requests = append(out.Requests, SupportItem{
			ID:       t.ID,
			Status:   string(t.StatusRequest),
			Category: string(t.CategoryRequest),
			Message:  t.Message,
		})
	}

	httputils.Success(w, r, out)
}

func (h *Handler) UpdateSupportStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//if isAdmin, _ := middleware.GetIsAdminFromContext(ctx); !isAdmin {
	//	httputils.ForbiddenError(w, r, "Доступ запрещен")
	//	return
	//}

	reqID, err := h.parseIDFromURL(r, "req_id")
	if err != nil {
		httputils.ValidationError(w, r, "Некорректный ID обращения", "req_id")
		return
	}

	var input UpdateStatusInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httputils.ValidationError(w, r, "Некорректный формат данных", "body")
		return
	}

	if err := h.supUC.UpdateSupportStatus(ctx, reqID, models.StatusContacting(input.Status)); err != nil {
		httputils.InternalError(w, r, err.Error())
		return
	}

	httputils.Success(w, r, map[string]string{"status": "ok"})
}

func (h *Handler) GetSupportStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//if isAdmin, _ := middleware.GetIsAdminFromContext(ctx); !isAdmin {
	//	httputils.ForbiddenError(w, r, "Доступ запрещен")
	//	return
	//}

	statsMap, err := h.supUC.GetSupportStats(ctx)
	if err != nil {
		httputils.InternalError(w, r, err.Error())
		return
	}

	out := StatsOutput{}
	for st, cnt := range statsMap {
		out.Items = append(out.Items, StatsItem{
			Status: string(st),
			Count:  cnt,
		})
	}

	httputils.Success(w, r, out)
}
