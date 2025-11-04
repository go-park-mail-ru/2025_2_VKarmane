package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewErrorResponse(t *testing.T) {
	tests := []struct {
		name    string
		error   string
		details string
		field   string
		code    ErrorCode
		want    ErrorResponse
	}{
		{
			name:    "valid error response",
			error:   "User not found",
			details: "User with ID 123 does not exist",
			field:   "user_id",
			code:    ErrCodeUserNotFound,
			want: ErrorResponse{
				Error:     "User not found",
				Details:   "User with ID 123 does not exist",
				Field:     "user_id",
				Code:      ErrCodeUserNotFound,
				Timestamp: "",
			},
		},
		{
			name:    "empty error response",
			error:   "",
			details: "",
			field:   "",
			code:    ErrCodeInternalError,
			want: ErrorResponse{
				Error:     "",
				Details:   "",
				Field:     "",
				Code:      ErrCodeInternalError,
				Timestamp: "",
			},
		},
		{
			name:    "invalid request error",
			error:   "Invalid request",
			details: "Missing required field",
			field:   "email",
			code:    ErrCodeInvalidRequest,
			want: ErrorResponse{
				Error:     "Invalid request",
				Details:   "Missing required field",
				Field:     "email",
				Code:      ErrCodeInvalidRequest,
				Timestamp: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewErrorResponse(tt.error, tt.details, tt.field, tt.code)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestErrorCode_GetErrorMessage(t *testing.T) {
	tests := []struct {
		name     string
		code     ErrorCode
		expected string
	}{
		// Request errors
		{"InvalidRequest", ErrCodeInvalidRequest, "Некорректный запрос"},
		{"MissingFields", ErrCodeMissingFields, "Отсутствуют обязательные поля"},
		{"InvalidEmail", ErrCodeInvalidEmail, "Некорректный email адрес"},
		{"InvalidPassword", ErrCodeInvalidPassword, "Некорректный пароль"},
		{"InvalidLogin", ErrCodeInvalidLogin, "Некорректный логин"},

		// User errors
		{"UserExists", ErrCodeUserExists, "Пользователь уже существует"},
		{"EmailExists", ErrCodeEmailExists, "Email уже используется"},
		{"LoginExists", ErrCodeLoginExists, "Логин уже занят"},
		{"WeakPassword", ErrCodeWeakPassword, "Пароль слишком слабый"},

		// Authentication errors
		{"UserNotFound", ErrCodeUserNotFound, "Пользователь не найден"},
		{"InvalidCredentials", ErrCodeInvalidCredentials, "Неверные учетные данные"},
		{"AccountLocked", ErrCodeAccountLocked, "Аккаунт заблокирован"},

		// Token errors
		{"TokenExpired", ErrCodeTokenExpired, "Токен истек"},
		{"TokenInvalid", ErrCodeTokenInvalid, "Недействительный токен"},
		{"TokenMissing", ErrCodeTokenMissing, "Токен отсутствует"},

		// Authorization errors
		{"Unauthorized", ErrCodeUnauthorized, "Требуется авторизация"},
		{"Forbidden", ErrCodeForbidden, "Доступ запрещен"},
		{"AccessDenied", ErrCodeAccessDenied, "Доступ отклонен"},

		// Resource errors
		{"ResourceNotFound", ErrCodeResourceNotFound, "Ресурс не найден"},
		{"ResourceExists", ErrCodeResourceExists, "Ресурс уже существует"},
		{"ResourceConflict", ErrCodeResourceConflict, "Конфликт ресурсов"},
		{"BudgetNotFound", ErrCodeBudgetNotFound, "Бюджет не найден"},
		{"AccountNotFound", ErrCodeAccountNotFound, "Счет не найден"},
		{"TransactionNotFound", ErrCodeTransactionNotFound, "Операция не найдена"},

		// Validation errors
		{"InvalidAmount", ErrCodeInvalidAmount, "Некорректная сумма"},
		{"InvalidCurrency", ErrCodeInvalidCurrency, "Некорректная валюта"},
		{"InvalidDate", ErrCodeInvalidDate, "Некорректная дата"},
		{"InvalidPeriod", ErrCodeInvalidPeriod, "Некорректный период"},

		// Server errors
		{"InternalError", ErrCodeInternalError, "Внутренняя ошибка сервера"},
		{"ServiceUnavailable", ErrCodeServiceUnavailable, "Сервис недоступен"},
		{"DatabaseError", ErrCodeDatabaseError, "Ошибка базы данных"},

		// Unknown error code
		{"UnknownError", ErrorCode("UNKNOWN_ERROR"), "Неизвестная ошибка"},
		{"EmptyError", ErrorCode(""), "Неизвестная ошибка"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.code.GetErrorMessage()
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestErrorCode_GetErrorMessage_AllCodes(t *testing.T) {
	// Проверяем, что все определенные коды ошибок имеют сообщения
	allCodes := []ErrorCode{
		ErrCodeInvalidRequest,
		ErrCodeMissingFields,
		ErrCodeInvalidEmail,
		ErrCodeInvalidPassword,
		ErrCodeInvalidLogin,
		ErrCodeUserExists,
		ErrCodeEmailExists,
		ErrCodeLoginExists,
		ErrCodeWeakPassword,
		ErrCodeUserNotFound,
		ErrCodeInvalidCredentials,
		ErrCodeAccountLocked,
		ErrCodeTokenExpired,
		ErrCodeTokenInvalid,
		ErrCodeTokenMissing,
		ErrCodeUnauthorized,
		ErrCodeForbidden,
		ErrCodeAccessDenied,
		ErrCodeResourceNotFound,
		ErrCodeResourceExists,
		ErrCodeResourceConflict,
		ErrCodeBudgetNotFound,
		ErrCodeAccountNotFound,
		ErrCodeTransactionNotFound,
		ErrCodeInvalidAmount,
		ErrCodeInvalidCurrency,
		ErrCodeInvalidDate,
		ErrCodeInvalidPeriod,
		ErrCodeInternalError,
		ErrCodeServiceUnavailable,
		ErrCodeDatabaseError,
	}

	for _, code := range allCodes {
		t.Run(string(code), func(t *testing.T) {
			msg := code.GetErrorMessage()
			assert.NotEmpty(t, msg, "Error code %s should have a message", code)
			assert.NotEqual(t, "Неизвестная ошибка", msg, "Error code %s should have a specific message, not default", code)
		})
	}
}

