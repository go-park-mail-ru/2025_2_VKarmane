package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/user"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/auth"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
	httputil "github.com/go-park-mail-ru/2025_2_VKarmane/pkg/http"
)

type Handler struct {
	authUC AuthUseCase
	logger logger.Logger
}

func NewHandler(authUC AuthUseCase, logger logger.Logger) *Handler {
	return &Handler{authUC: authUC, logger: logger}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.ValidationError(w, r, "Некорректный формат данных", "body")
		return
	}

	validationErrors := utils.ValidateStruct(req)
	if len(validationErrors) > 0 {
		httputil.ValidationErrors(w, r, validationErrors)
		return
	}

	response, err := h.authUC.Register(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, user.EmailExistsErr):
			httputil.ConflictError(w, r, "Пользователь с таким email уже существует", models.ErrCodeEmailExists)
		case errors.Is(err, user.LoginExistsErr):
			httputil.ConflictError(w, r, "Пользователь с таким логином уже существует", models.ErrCodeLoginExists)
		default:
			httputil.ConflictError(w, r, "Пользователь уже существует", models.ErrCodeUserExists)
		}
		return
	}

	isProduction := os.Getenv("ENV") == "production"
	utils.SetAuthCookie(w, response.Token, isProduction)

	response.Token = ""
	httputil.Created(w, r, response)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.ValidationError(w, r, "Некорректный формат данных", "body")
		return
	}

	validationErrors := utils.ValidateStruct(req)
	if len(validationErrors) > 0 {
		httputil.ValidationErrors(w, r, validationErrors)
		return
	}

	response, err := h.authUC.Login(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, user.UserNotFound):
			httputil.UnauthorizedError(w, r, "Пользователь не найден", models.ErrCodeUserNotFound)
		case errors.Is(err, auth.InvalidPassword):
			httputil.UnauthorizedError(w, r, "Неверный пароль", models.ErrCodeInvalidCredentials)
		default:
			httputil.UnauthorizedError(w, r, "Неверные учетные данные", models.ErrCodeInvalidCredentials)
		}
		return
	}

	isProduction := os.Getenv("ENV") == "production"
	utils.SetAuthCookie(w, response.Token, isProduction)

	response.Token = ""
	httputil.Success(w, r, response)
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		httputil.UnauthorizedError(w, r, "Требуется авторизация", models.ErrCodeUnauthorized)
		return
	}

	user, err := h.authUC.GetUserByID(r.Context(), userID)
	if err != nil {
		httputil.NotFoundError(w, r, "Пользователь не найден")
		return
	}

	httputil.Success(w, r, user)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	isProduction := os.Getenv("ENV") == "production"
	utils.ClearAuthCookie(w, isProduction)

	httputil.Success(w, r, map[string]string{"message": "Logged out successfully"})
}
