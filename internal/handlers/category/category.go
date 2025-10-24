package category

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/category"
	httputils "github.com/go-park-mail-ru/2025_2_VKarmane/pkg/http"
)

type Handler struct {
	categoryUC category.CategoryUseCase
}

func NewHandler(categoryUC category.CategoryUseCase) *Handler {
	return &Handler{
		categoryUC: categoryUC,
	}
}

func (h *Handler) getUserID(r *http.Request) (int, bool) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	return userID, ok
}


// GetCategories godoc
// @Summary Получение списка категорий пользователя
// @Description Возвращает список всех категорий пользователя с количеством операций
// @Tags categories
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "Список категорий пользователя"
// @Failure 401 {object} models.ErrorResponse "Требуется аутентификация (UNAUTHORIZED, TOKEN_MISSING, TOKEN_INVALID, TOKEN_EXPIRED)"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера (INTERNAL_ERROR, DATABASE_ERROR)"
// @Router /categories [get]
func (h *Handler) GetCategories(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r)
	if !ok {
		httputils.UnauthorizedError(w, r, "User not authenticated", models.ErrCodeUnauthorized)
		return
	}

	categories, err := h.categoryUC.GetCategoriesByUser(r.Context(), userID)
	if err != nil {
		httputils.InternalError(w, r, "Failed to get categories")
		return
	}

	response := map[string]interface{}{
		"user_id":    userID,
		"categories": categories,
	}

	httputils.Success(w, r, response)
}
