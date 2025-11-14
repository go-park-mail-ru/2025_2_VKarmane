package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	// "github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/user"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/auth"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
	httputil "github.com/go-park-mail-ru/2025_2_VKarmane/pkg/http"
	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/proto"
)

type Handler struct {
	authUC AuthUseCase
	clock  clock.Clock
	logger logger.Logger
	authClient authpb.AuthServiceClient
}

func NewHandler(authUC AuthUseCase, clck clock.Clock, logger logger.Logger, authCLient authpb.AuthServiceClient) *Handler {
	return &Handler{authUC: authUC, clock: clck, logger: logger, authClient: authCLient}
}

// Register godoc
// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя в системе
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Данные для регистрации"
// @Success 201 {object} models.AuthResponse "Пользователь успешно создан"
// @Failure 400 {object} models.ErrorResponse "Некорректные данные (INVALID_REQUEST, MISSING_FIELDS, INVALID_EMAIL, INVALID_PASSWORD, INVALID_LOGIN, WEAK_PASSWORD)"
// @Failure 409 {object} models.ErrorResponse "Конфликт (USER_EXISTS, EMAIL_EXISTS, LOGIN_EXISTS)"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера (INTERNAL_ERROR, DATABASE_ERROR)"
// @Router /auth/register [post]
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

	// response, err := h.authUC.Register(r.Context(), req)
	response, err := h.authClient.Register(r.Context(), &authpb.RegisterRequest{
		Login: req.Login,
		Email: req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, user.ErrEmailExists):
			httputil.ConflictError(w, r, "Пользователь с таким email уже существует", models.ErrCodeEmailExists)
		case errors.Is(err, user.ErrLoginExists):
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

// Login godoc
// @Summary Вход в систему
// @Description Аутентификация пользователя по логину и паролю
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Данные для входа"
// @Success 200 {object} models.AuthResponse "Успешный вход"
// @Failure 400 {object} models.ErrorResponse "Некорректные данные (INVALID_REQUEST, MISSING_FIELDS, INVALID_LOGIN, INVALID_PASSWORD)"
// @Failure 401 {object} models.ErrorResponse "Неверные учетные данные (INVALID_CREDENTIALS, USER_NOT_FOUND, ACCOUNT_LOCKED)"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера (INTERNAL_ERROR, DATABASE_ERROR)"
// @Router /auth/login [post]
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

	// response, err := h.authUC.Login(r.Context(), req)
	response, err := h.authClient.Login(r.Context(), &authpb.LoginRequest{
		Login: req.Login,
		Password: req.Password,
	})
	
	if err != nil {
		log := logger.FromContext(r.Context())
		log.Error("grpc error", err)
		switch {
		case errors.Is(err, user.ErrUserNotFound):
			httputil.UnauthorizedError(w, r, "Пользователь не найден", models.ErrCodeUserNotFound)
		case errors.Is(err, auth.ErrInvalidPassword):
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

// GetCSRFToken godoc
// @Summary Получение CSRF токена
// @Description Возвращает CSRF токен для использования в последующих запросах
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string "CSRF токен"
// @Router /auth/csrf [get]
func (h *Handler) GetCSRFToken(w http.ResponseWriter, r *http.Request) {
	isProduction := os.Getenv("ENV") == "production"

	token, err := h.authUC.GetCSRFToken(r.Context())
	if err != nil {
		httputil.InternalError(w, r, "Failed to get CSRF-Token")
	}

	utils.SetCSRFCookie(w, token, isProduction)
	httputil.Success(w, r, map[string]string{"csrf_token": token})
}

// Logout godoc
// @Summary Выход из системы
// @Description Завершает сессию пользователя
// @Tags auth
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]string "Успешный выход"
// @Failure 401 {object} models.ErrorResponse "Требуется аутентификация (UNAUTHORIZED, TOKEN_MISSING, TOKEN_INVALID, TOKEN_EXPIRED)"
// @Router /auth/logout [post]
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	isProduction := os.Getenv("ENV") == "production"
	utils.ClearAuthCookie(w, isProduction)
	utils.ClearCSRFCookie(w, isProduction)
	httputil.Success(w, r, map[string]string{"message": "Logged out successfully"})
}
