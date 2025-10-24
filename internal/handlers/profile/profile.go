package profile

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/profile"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
	httputils "github.com/go-park-mail-ru/2025_2_VKarmane/pkg/http"
)

type Handler struct {
	profileUC profile.ProfileUseCase
}

func NewHandler(profileUC profile.ProfileUseCase) *Handler {
	return &Handler{
		profileUC: profileUC,
	}
}

func (h *Handler) getUserID(r *http.Request) (int, bool) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	return userID, ok
}

// GetProfile godoc
// @Summary Получение профиля пользователя
// @Description Возвращает информацию о профиле текущего пользователя
// @Tags profile
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.ProfileResponse "Профиль пользователя"
// @Failure 401 {object} models.ErrorResponse "Требуется аутентификация (UNAUTHORIZED, TOKEN_MISSING, TOKEN_INVALID, TOKEN_EXPIRED)"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера (INTERNAL_ERROR, DATABASE_ERROR)"
// @Router /profile [get]
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r)
	if !ok {
		httputils.UnauthorizedError(w, r, "User not authenticated", models.ErrCodeUnauthorized)
		return
	}

	profile, err := h.profileUC.GetProfile(r.Context(), userID)
	if err != nil {
		httputils.InternalError(w, r, "Failed to get profile")
		return
	}

	httputils.Success(w, r, profile)
}

// UpdateProfile godoc
// @Summary Обновление профиля пользователя
// @Description Обновляет информацию о профиле текущего пользователя
// @Tags profile
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body models.UpdateProfileRequest true "Данные для обновления профиля"
// @Success 200 {object} models.ProfileResponse "Обновленный профиль пользователя"
// @Failure 400 {object} models.ErrorResponse "Некорректные данные (INVALID_REQUEST, MISSING_FIELDS)"
// @Failure 401 {object} models.ErrorResponse "Требуется аутентификация (UNAUTHORIZED, TOKEN_MISSING, TOKEN_INVALID, TOKEN_EXPIRED)"
// @Failure 409 {object} models.ErrorResponse "Email уже используется (RESOURCE_EXISTS)"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера (INTERNAL_ERROR, DATABASE_ERROR)"
// @Router /profile/edit [put]
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r)
	if !ok {
		httputils.UnauthorizedError(w, r, "User not authenticated", models.ErrCodeUnauthorized)
		return
	}

	var req models.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.Error(w, r, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		httputils.ValidationErrors(w, r, err)
		return
	}

	profile, err := h.profileUC.UpdateProfile(r.Context(), req, userID)
	if err != nil {
		httputils.InternalError(w, r, "Failed to update profile")
		return
	}

	httputils.Success(w, r, profile)
}
