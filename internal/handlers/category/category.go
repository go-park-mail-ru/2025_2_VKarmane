package category

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/category"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
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
// @Description Создает новую категорию для пользователя
// @Tags categories
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param category body models.CreateCategoryRequest true "Данные категории"
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

	var req models.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ValidationError(w, r, "Invalid request body", "body")
		return
	}

	validationErrors := utils.ValidateStruct(req)
	if len(validationErrors) > 0 {
		httputils.ValidationErrors(w, r, validationErrors)
		return
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
// @Description Обновляет информацию о категории
// @Tags categories
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID категории"
// @Param category body models.UpdateCategoryRequest true "Данные для обновления"
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

	var req models.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ValidationError(w, r, "Invalid request body", "body")
		return
	}

	validationErrors := utils.ValidateStruct(req)
	if len(validationErrors) > 0 {
		httputils.ValidationErrors(w, r, validationErrors)
		return
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
