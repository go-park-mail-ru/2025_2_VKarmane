package budget

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	bdgpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	httputils "github.com/go-park-mail-ru/2025_2_VKarmane/pkg/http"
)

type Handler struct {
	budgetClient bdgpb.BudgetServiceClient
	clock    clock.Clock
}

func NewHandler(clck clock.Clock, budgetClient bdgpb.BudgetServiceClient) *Handler {
	return &Handler{clock: clck, budgetClient: budgetClient}
}

func (h *Handler) getUserID(r *http.Request) (int, bool) {
	return middleware.GetUserIDFromContext(r.Context())
}

func (h *Handler) parseIDFromURL(r *http.Request, paramName string) (int, error) {
	vars := mux.Vars(r)
	idStr := vars[paramName]
	return strconv.Atoi(idStr)
}

// GetListBudgets godoc
// @Summary Получение списка бюджетов пользователя
// @Description Возвращает список всех бюджетов пользователя с расчетом фактических расходов
// @Tags budget
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "Список бюджетов пользователя"
// @Failure 401 {object} models.ErrorResponse "Требуется аутентификация (UNAUTHORIZED, TOKEN_MISSING, TOKEN_INVALID, TOKEN_EXPIRED)"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера (INTERNAL_ERROR, DATABASE_ERROR)"
// @Router /budgets [get]
func (h *Handler) GetListBudgets(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r)
	if !ok {
		httputils.Error(w, r, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// budgets, err := h.budgetUC.GetBudgetsForUser(r.Context(), userID)
	budgets, err := h.budgetClient.GetListBudgets(r.Context(), &bdgpb.UserID{
    	UserID: int32(userID),
	})
	if err != nil {
		log := logger.FromContext(r.Context())

		st, ok := status.FromError(err)
		if log != nil {
			if ok {
				log.Error("grpc GetListBudgets failed",
					"code", st.Code(),
					"message", st.Message(),
				)
			} else {
				log.Error("grpc GetListBudgets unknown error", "error", err)
			}
		}

		httputils.InternalError(w, r, "Failed to get budgets for user")
		return
	}

	// budgetsDTO := BudgetsToAPI(userID, budgets)
	httputils.Success(w, r, BudgetsToAPI(userID, budgets))
}

// GetBudgetByID godoc
// @Summary Получение конкретного бюджета
// @Description Возвращает информацию о конкретном бюджете пользователя
// @Tags budget
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID бюджета"
// @Success 200 {object} map[string]interface{} "Информация о бюджете"
// @Failure 400 {object} models.ErrorResponse "Некорректный ID бюджета (INVALID_REQUEST)"
// @Failure 401 {object} models.ErrorResponse "Требуется аутентификация (UNAUTHORIZED, TOKEN_MISSING, TOKEN_INVALID, TOKEN_EXPIRED)"
// @Failure 404 {object} models.ErrorResponse "Бюджет не найден (BUDGET_NOT_FOUND)"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера (INTERNAL_ERROR, DATABASE_ERROR)"
// @Router /budget/{id} [get]
func (h *Handler) GetBudgetByID(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromURL(r, "id")
	if err != nil {
		httputils.ValidationError(w, r, "Invalid budget ID format", "id")
		return
	}

	userID, ok := h.getUserID(r)
	if !ok {
		httputils.Error(w, r, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// budget, err := h.budgetUC.GetBudgetByID(r.Context(), userID, id)
	budget, err := h.budgetClient.GetBudget(r.Context(), IDsToBudgetRequest(id, userID))
	if err != nil {
		st, ok := status.FromError(err)
		log := logger.FromContext(r.Context())
		if !ok {
			if log != nil {
				log.Error("grpc GetListBudgets unknown error", "error", err)
			}
			httputils.InternalError(w, r, "failed to get budget")
			return
		}
		switch st.Code() {
			case codes.NotFound:
				if log != nil {
					log.Error("grpc GetBudget invalid arg", "error", err)
				}
				httputils.NotFoundError(w, r, "Бюджет не найден")
				return
			default:
				if log != nil {
					log.Error("grpc GetBudget error", "error", err)
				}
				httputils.InternalError(w, r, "failed to get budget")
				return
		}
	}

	budgetDTO := BudgetToAPI(budget)
	httputils.Success(w, r, budgetDTO)
}

func (h *Handler) CreateBudget(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r)
	if !ok {
		httputils.Error(w, r, "User not authenticated", http.StatusUnauthorized)
		return
	}
	var req models.CreateBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ValidationError(w, r, "Некорректный формат данных", "body")
		return
	}

	validationErrors := utils.ValidateStruct(req)
	if len(validationErrors) > 0 {
		httputils.ValidationErrors(w, r, validationErrors)
		return
	}

	// budget, err := h.budgetUC.CreateBudget(r.Context(), req, userID)
	budget, err := h.budgetClient.CreateBudget(r.Context(), ModelCreateReqtoProtoReq(req, userID))
	if err != nil {
		st, ok := status.FromError(err)
		log := logger.FromContext(r.Context())
		if !ok {
			if log != nil {
				log.Error("grpc Createbudget unknown error", "error", err)
			}
			httputils.InternalError(w, r, "failed to create budget")
			return
		}
		
		switch st.Code() {
		case codes.InvalidArgument:
			if log != nil {
				log.Error("grpc Createbudget invalid arg", "error", err)
			}
			httputils.Error(w, r, "failed to create budget", http.StatusBadRequest)
			return
		case codes.AlreadyExists:
			if log != nil {
				log.Error("grpc Createbudget exists error", "error", err)
			}
			httputils.ConflictError(w, r, "Незакрытый бюджет с такой категорией уже существует", models.ErrCodeBudgetExists)
			return
		default:
			if log != nil {
				log.Error("grpc CreateBudget error", "error", err)
			}
			httputils.InternalError(w, r, "failed to create budget")
			return
		}
	}
	budgetDTO := BudgetToAPI(budget)
	httputils.Success(w, r, budgetDTO)
}

func (h *Handler) UpdateBudget(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r)
	if !ok {
		httputils.Error(w, r, "User not authenticated", http.StatusUnauthorized)
		return
	}

	budgetID, err := h.parseIDFromURL(r, "id")
	if err != nil {
		httputils.ValidationError(w, r, "Некорректный ID бюджета", "budgetID")
		return
	}
	var req models.UpdatedBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ValidationError(w, r, "Некорректный формат данных", "body")
		return
	}

	// budget, err := h.budgetUC.UpdateBudget(r.Context(), req, userID, budgetID)
	budget, err := h.budgetClient.UpdateBudget(r.Context(), ModelUpdateReqtoProtoReq(req, budgetID, userID))
	if err != nil {
		st, ok := status.FromError(err)
		log := logger.FromContext(r.Context())
		if !ok {
			if log != nil {
				log.Error("grpc Updatebudget unknown error", "error", err)
			}
			httputils.InternalError(w, r, "failed to update budget")
			return
		}
		switch st.Code() {
		case codes.NotFound:
			if log != nil {
				log.Error("grpc UpdateBudget invalid arg", "error", err)
			}
			httputils.NotFoundError(w, r, "Бюджет не найден")
			return
		case codes.PermissionDenied:
			if log != nil {
				log.Error("grpc UpdateBudget forbidden", "error", err)
			}
			httputils.UnauthorizedError(w, r, "Отказано в доступе", models.ErrCodeForbidden)
			return
		default:
			if log != nil {
				log.Error("grpc UpdateBudget error", "error", err)
			}
			httputils.InternalError(w, r, "failed to update budget")
			return
		}
	}
	budgetDTO := BudgetToAPI(budget)
	httputils.Success(w, r, budgetDTO)
}

func (h *Handler) DeleteBudget(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r)
	if !ok {
		httputils.Error(w, r, "User not authenticated", http.StatusUnauthorized)
		return
	}

	budgetID, err := h.parseIDFromURL(r, "id")
	if err != nil {
		httputils.ValidationError(w, r, "Некорректный ID бюджета", "budgetID")
		return
	}

	// budget, err := h.budgetUC.DeleteBudget(r.Context(), userID, budgetID)
	budget, err := h.budgetClient.DeleteBudget(r.Context(), IDsToBudgetRequest(budgetID, userID))
	if err != nil {
		st, ok := status.FromError(err)
		log := logger.FromContext(r.Context())
		if !ok {
			if log != nil {
				log.Error("grpc Deletebudget unknown error", "error", err)
			}
			httputils.InternalError(w, r, "failed to delete budget")
			return
		}
		switch st.Code() {
		case codes.NotFound:
			if log != nil {
				log.Error("grpc DeleteBudget invalid arg", "error", err)
			}
			httputils.NotFoundError(w, r, "Бюджет не найден")
			return
		case codes.PermissionDenied:
			if log != nil {
				log.Error("grpc DeleteBudget forbidden", "error", err)
			}
			httputils.UnauthorizedError(w, r, "Отказано в доступе", models.ErrCodeForbidden)
			return
		default:
			if log != nil {
				log.Error("grpc DeleteBudget error", "error", err)
			}
			httputils.InternalError(w, r, "failed to update budget")
			return
		}
	}
	budgetDTO := BudgetToAPI(budget)
	httputils.Success(w, r, budgetDTO)
}
