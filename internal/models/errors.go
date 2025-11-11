package models

type ErrorCode string

const (
	ErrCodeInvalidRequest  ErrorCode = "INVALID_REQUEST"
	ErrCodeMissingFields   ErrorCode = "MISSING_FIELDS"
	ErrCodeInvalidEmail    ErrorCode = "INVALID_EMAIL"
	ErrCodeInvalidPassword ErrorCode = "INVALID_PASSWORD"
	ErrCodeInvalidLogin    ErrorCode = "INVALID_LOGIN"

	ErrCodeUserExists   ErrorCode = "USER_EXISTS"
	ErrCodeEmailExists  ErrorCode = "EMAIL_EXISTS"
	ErrCodeLoginExists  ErrorCode = "LOGIN_EXISTS"
	ErrCodeWeakPassword ErrorCode = "WEAK_PASSWORD"

	ErrCodeBudgetExists ErrorCode = "BUDGET_EXISTS"

	ErrCodeUserNotFound       ErrorCode = "USER_NOT_FOUND"
	ErrCodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	ErrCodeAccountLocked      ErrorCode = "ACCOUNT_LOCKED"

	ErrCodeTokenExpired ErrorCode = "TOKEN_EXPIRED"
	ErrCodeTokenInvalid ErrorCode = "TOKEN_INVALID"
	ErrCodeTokenMissing ErrorCode = "TOKEN_MISSING"

	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden    ErrorCode = "FORBIDDEN"
	ErrCodeAccessDenied ErrorCode = "ACCESS_DENIED"

	ErrCodeResourceNotFound    ErrorCode = "RESOURCE_NOT_FOUND"
	ErrCodeResourceExists      ErrorCode = "RESOURCE_EXISTS"
	ErrCodeResourceConflict    ErrorCode = "RESOURCE_CONFLICT"
	ErrCodeBudgetNotFound      ErrorCode = "BUDGET_NOT_FOUND"
	ErrCodeAccountNotFound     ErrorCode = "ACCOUNT_NOT_FOUND"
	ErrCodeTransactionNotFound ErrorCode = "OPERATION_NOT_FOUND"

	ErrCodeInvalidAmount   ErrorCode = "INVALID_AMOUNT"
	ErrCodeInvalidCurrency ErrorCode = "INVALID_CURRENCY"
	ErrCodeInvalidDate     ErrorCode = "INVALID_DATE"
	ErrCodeInvalidPeriod   ErrorCode = "INVALID_PERIOD"

	ErrCodeInternalError      ErrorCode = "INTERNAL_ERROR"
	ErrCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
	ErrCodeDatabaseError      ErrorCode = "DATABASE_ERROR"
)

type ErrorResponse struct {
	Error     string    `json:"error"`
	Code      ErrorCode `json:"code"`
	Details   string    `json:"details"`
	Field     string    `json:"field"`
	Timestamp string    `json:"timestamp"`
}

func NewErrorResponse(error, details, field string, code ErrorCode) ErrorResponse {
	return ErrorResponse{
		Error:     error,
		Code:      code,
		Details:   details,
		Field:     field,
		Timestamp: "",
	}
}

func (e ErrorCode) GetErrorMessage() string {
	messages := map[ErrorCode]string{
		ErrCodeInvalidRequest:  "Некорректный запрос",
		ErrCodeMissingFields:   "Отсутствуют обязательные поля",
		ErrCodeInvalidEmail:    "Некорректный email адрес",
		ErrCodeInvalidPassword: "Некорректный пароль",
		ErrCodeInvalidLogin:    "Некорректный логин",

		ErrCodeUserExists:   "Пользователь уже существует",
		ErrCodeEmailExists:  "Email уже используется",
		ErrCodeLoginExists:  "Логин уже занят",
		ErrCodeWeakPassword: "Пароль слишком слабый",

		ErrCodeUserNotFound:       "Пользователь не найден",
		ErrCodeInvalidCredentials: "Неверные учетные данные",
		ErrCodeAccountLocked:      "Аккаунт заблокирован",

		ErrCodeTokenExpired: "Токен истек",
		ErrCodeTokenInvalid: "Недействительный токен",
		ErrCodeTokenMissing: "Токен отсутствует",

		ErrCodeUnauthorized: "Требуется авторизация",
		ErrCodeForbidden:    "Доступ запрещен",
		ErrCodeAccessDenied: "Доступ отклонен",

		ErrCodeResourceNotFound:    "Ресурс не найден",
		ErrCodeResourceExists:      "Ресурс уже существует",
		ErrCodeResourceConflict:    "Конфликт ресурсов",
		ErrCodeBudgetNotFound:      "Бюджет не найден",
		ErrCodeAccountNotFound:     "Счет не найден",
		ErrCodeTransactionNotFound: "Операция не найдена",

		ErrCodeInvalidAmount:   "Некорректная сумма",
		ErrCodeInvalidCurrency: "Некорректная валюта",
		ErrCodeInvalidDate:     "Некорректная дата",
		ErrCodeInvalidPeriod:   "Некорректный период",

		ErrCodeInternalError:      "Внутренняя ошибка сервера",
		ErrCodeServiceUnavailable: "Сервис недоступен",
		ErrCodeDatabaseError:      "Ошибка базы данных",
	}

	if msg, exists := messages[e]; exists {
		return msg
	}
	return "Неизвестная ошибка"
}
