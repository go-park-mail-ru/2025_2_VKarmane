package profile

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/image"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/profile"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
	httputils "github.com/go-park-mail-ru/2025_2_VKarmane/pkg/http"
)

type Handler struct {
	profileUC profile.ProfileUseCase
	imageUC   image.ImageUseCase
}

func NewHandler(profileUC profile.ProfileUseCase, imageUC image.ImageUseCase) *Handler {
	return &Handler{
		profileUC: profileUC,
		imageUC:   imageUC,
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
// @Description Обновляет информацию о профиле текущего пользователя. Поддерживает multipart/form-data с опциональным полем avatar для загрузки аватарки
// @Tags profile
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param first_name formData string true "Имя"
// @Param last_name formData string true "Фамилия"
// @Param email formData string true "Email"
// @Param avatar formData file false "Аватарка (опционально)"
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

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		httputils.Error(w, r, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	req := models.UpdateProfileRequest{
		FirstName: r.FormValue("first_name"),
		LastName:  r.FormValue("last_name"),
		Email:     r.FormValue("email"),
	}

	if err := utils.ValidateStruct(req); err != nil {
		httputils.ValidationErrors(w, r, err)
		return
	}

	file, header, err := r.FormFile("avatar")
	if err == nil {
		defer func() { _ = file.Close() }()

		ext := strings.ToLower(filepath.Ext(header.Filename))
		allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
		allowed := false
		for _, allowedExt := range allowedExts {
			if ext == allowedExt {
				allowed = true
				break
			}
		}

		if !allowed {
			httputils.Error(w, r, "Invalid image format", http.StatusBadRequest)
			return
		}

		contentType := header.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "image/" + ext[1:]
		}

		imageID, err := h.imageUC.UploadImage(r.Context(), file, header.Filename, header.Size, contentType)
		if err != nil {
			httputils.InternalError(w, r, "Failed to upload avatar")
			return
		}

		req.LogoHashedID = imageID
	}

	profile, err := h.profileUC.UpdateProfile(r.Context(), req, userID)
	if err != nil {
		httputils.InternalError(w, r, "Failed to update profile")
		return
	}

	httputils.Success(w, r, profile)
}
