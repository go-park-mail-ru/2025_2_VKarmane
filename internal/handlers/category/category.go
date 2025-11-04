package category

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/category"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/image"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
	httputils "github.com/go-park-mail-ru/2025_2_VKarmane/pkg/http"
)

type Handler struct {
	categoryUC category.CategoryUseCase
	imageUC    image.ImageUseCase
}

func NewHandler(categoryUC category.CategoryUseCase, imageUC image.ImageUseCase) *Handler {
	return &Handler{
		categoryUC: categoryUC,
		imageUC:    imageUC,
	}
}

func (h *Handler) getUserID(r *http.Request) (int, bool) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	return userID, ok
}

func (h *Handler) parseIDFromURL(r *http.Request, paramName string) (int, error) {
	vars := mux.Vars(r)
	idStr := vars[paramName]
	return strconv.Atoi(idStr)
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

// CreateCategory godoc
// @Summary Создание новой категории
// @Description Создает новую категорию для пользователя. Поддерживает multipart/form-data с опциональным полем image для загрузки картинки категории
// @Tags categories
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param name formData string true "Название категории"
// @Param description formData string false "Описание категории"
// @Param image formData file false "Картинка категории (опционально)"
// @Success 201 {object} models.Category "Созданная категория"
// @Failure 400 {object} models.ErrorResponse "Некорректные данные (VALIDATION_ERROR, INVALID_INPUT)"
// @Failure 401 {object} models.ErrorResponse "Требуется аутентификация (UNAUTHORIZED, TOKEN_MISSING, TOKEN_INVALID, TOKEN_EXPIRED)"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера (INTERNAL_ERROR, DATABASE_ERROR)"
// @Router /categories [post]
func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r)
	if !ok {
		httputils.UnauthorizedError(w, r, "User not authenticated", models.ErrCodeUnauthorized)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		httputils.Error(w, r, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	req := models.CreateCategoryRequest{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
	}

	validationErrors := utils.ValidateStruct(req)
	if len(validationErrors) > 0 {
		httputils.ValidationErrors(w, r, validationErrors)
		return
	}

	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

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
			httputils.InternalError(w, r, "Failed to upload image")
			return
		}

		req.LogoHashedID = imageID
	}

	category, err := h.categoryUC.CreateCategory(r.Context(), req, userID)
	if err != nil {
		httputils.InternalError(w, r, "Failed to create category")
		return
	}

	httputils.Created(w, r, category)
}

// GetCategoryByID godoc
// @Summary Получение категории по ID
// @Description Возвращает информацию о категории по её ID с количеством операций
// @Tags categories
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID категории"
// @Success 200 {object} models.CategoryWithStats "Информация о категории"
// @Failure 401 {object} models.ErrorResponse "Требуется аутентификация (UNAUTHORIZED, TOKEN_MISSING, TOKEN_INVALID, TOKEN_EXPIRED)"
// @Failure 404 {object} models.ErrorResponse "Категория не найдена (RESOURCE_NOT_FOUND)"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера (INTERNAL_ERROR, DATABASE_ERROR)"
// @Router /categories/{id} [get]
func (h *Handler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r)
	if !ok {
		httputils.UnauthorizedError(w, r, "User not authenticated", models.ErrCodeUnauthorized)
		return
	}

	categoryID, err := h.parseIDFromURL(r, "id")
	if err != nil {
		httputils.ValidationError(w, r, "Invalid category ID", "id")
		return
	}

	category, err := h.categoryUC.GetCategoryByID(r.Context(), userID, categoryID)
	if err != nil {
		httputils.NotFoundError(w, r, "Category not found")
		return
	}

	httputils.Success(w, r, category)
}

// UpdateCategory godoc
// @Summary Обновление категории
// @Description Обновляет информацию о категории. Поддерживает multipart/form-data с опциональным полем image для загрузки картинки категории
// @Tags categories
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID категории"
// @Param name formData string false "Название категории"
// @Param description formData string false "Описание категории"
// @Param image formData file false "Картинка категории (опционально)"
// @Success 200 {object} models.Category "Обновленная категория"
// @Failure 400 {object} models.ErrorResponse "Некорректные данные (VALIDATION_ERROR, INVALID_INPUT)"
// @Failure 401 {object} models.ErrorResponse "Требуется аутентификация (UNAUTHORIZED, TOKEN_MISSING, TOKEN_INVALID, TOKEN_EXPIRED)"
// @Failure 404 {object} models.ErrorResponse "Категория не найдена (RESOURCE_NOT_FOUND)"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера (INTERNAL_ERROR, DATABASE_ERROR)"
// @Router /categories/{id} [put]
func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r)
	if !ok {
		httputils.UnauthorizedError(w, r, "User not authenticated", models.ErrCodeUnauthorized)
		return
	}

	categoryID, err := h.parseIDFromURL(r, "id")
	if err != nil {
		httputils.ValidationError(w, r, "Invalid category ID", "id")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		httputils.Error(w, r, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	req := models.UpdateCategoryRequest{}
	if name := r.FormValue("name"); name != "" {
		req.Name = &name
	}
	if description := r.FormValue("description"); description != "" {
		req.Description = &description
	}

	validationErrors := utils.ValidateStruct(req)
	if len(validationErrors) > 0 {
		httputils.ValidationErrors(w, r, validationErrors)
		return
	}

	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

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
			httputils.InternalError(w, r, "Failed to upload image")
			return
		}

		req.LogoHashedID = &imageID
	}

	category, err := h.categoryUC.UpdateCategory(r.Context(), req, userID, categoryID)
	if err != nil {
		httputils.NotFoundError(w, r, "Category not found")
		return
	}

	httputils.Success(w, r, category)
}

// DeleteCategory godoc
// @Summary Удаление категории
// @Description Удаляет категорию пользователя
// @Tags categories
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID категории"
// @Success 204 "Категория успешно удалена"
// @Failure 401 {object} models.ErrorResponse "Требуется аутентификация (UNAUTHORIZED, TOKEN_MISSING, TOKEN_INVALID, TOKEN_EXPIRED)"
// @Failure 404 {object} models.ErrorResponse "Категория не найдена (RESOURCE_NOT_FOUND)"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера (INTERNAL_ERROR, DATABASE_ERROR)"
// @Router /categories/{id} [delete]
func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r)
	if !ok {
		httputils.UnauthorizedError(w, r, "User not authenticated", models.ErrCodeUnauthorized)
		return
	}

	categoryID, err := h.parseIDFromURL(r, "id")
	if err != nil {
		httputils.ValidationError(w, r, "Invalid category ID", "id")
		return
	}

	err = h.categoryUC.DeleteCategory(r.Context(), userID, categoryID)
	if err != nil {
		httputils.NotFoundError(w, r, "Category not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
