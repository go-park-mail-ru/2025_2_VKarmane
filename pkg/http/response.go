package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
)

func JSON(w http.ResponseWriter, r *http.Request, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		Error(w, r, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func Error(w http.ResponseWriter, r *http.Request, message string, statusCode int) {
	if log := logger.FromContext(r.Context()); log != nil {
		log.Error("HTTP Error", "message", message, "status", statusCode, "path", r.URL.Path)
	}

	http.Error(w, message, statusCode)
}

func ErrorWithCode(w http.ResponseWriter, r *http.Request, errorResp models.ErrorResponse, statusCode int) {
	errorResp.Timestamp = time.Now().Format(time.RFC3339)

	if log := logger.FromContext(r.Context()); log != nil {
		log.Error("HTTP Error",
			"message", errorResp.Error,
			"code", errorResp.Code,
			"details", errorResp.Details,
			"field", errorResp.Field,
			"status", statusCode,
			"path", r.URL.Path)
	}

	JSON(w, r, errorResp, statusCode)
}

func ValidationError(w http.ResponseWriter, r *http.Request, message, field string) {
	ErrorWithCode(w, r, models.NewErrorResponse(message, "", field, models.ErrCodeInvalidRequest), http.StatusBadRequest)
}

func UnauthorizedError(w http.ResponseWriter, r *http.Request, message string, code models.ErrorCode) {
	ErrorWithCode(w, r, models.NewErrorResponse(message, "", "", code), http.StatusUnauthorized)
}

func NotFoundError(w http.ResponseWriter, r *http.Request, message string) {
	ErrorWithCode(w, r, models.NewErrorResponse(message, "", "", models.ErrCodeResourceNotFound), http.StatusNotFound)
}

func ConflictError(w http.ResponseWriter, r *http.Request, message string, code models.ErrorCode) {
	ErrorWithCode(w, r, models.NewErrorResponse(message, "", "", code), http.StatusConflict)
}

func InternalError(w http.ResponseWriter, r *http.Request, message string) {
	ErrorWithCode(w, r, models.NewErrorResponse(message, "", "", models.ErrCodeInternalError), http.StatusInternalServerError)
}

func Success(w http.ResponseWriter, r *http.Request, data interface{}) {
	JSON(w, r, data, http.StatusOK)
}

func Created(w http.ResponseWriter, r *http.Request, data interface{}) {
	JSON(w, r, data, http.StatusCreated)
}

func NoContent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func ValidationErrors(w http.ResponseWriter, r *http.Request, validationErrors utils.ValidationErrors) {
	errorResponse := models.ErrorResponse{
		Error:     "Ошибки валидации",
		Code:      models.ErrCodeInvalidRequest,
		Details:   validationErrors.Error(),
		Field:     "",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if log := logger.FromContext(r.Context()); log != nil {
		log.Error("Validation Error",
			"errors", validationErrors,
			"path", r.URL.Path)
	}

	JSON(w, r, map[string]interface{}{
		"error":             errorResponse.Error,
		"code":              errorResponse.Code,
		"details":           errorResponse.Details,
		"timestamp":         errorResponse.Timestamp,
		"validation_errors": validationErrors,
	}, http.StatusBadRequest)
}
