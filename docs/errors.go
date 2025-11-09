package docs

// ErrorCodes содержит все возможные коды ошибок API
// Этот файл используется для генерации документации Swagger

const (
	// Ошибки валидации
	ErrCodeInvalidRequest  = "INVALID_REQUEST"  // Некорректный запрос
	ErrCodeMissingFields   = "MISSING_FIELDS"   // Отсутствуют обязательные поля
	ErrCodeInvalidEmail    = "INVALID_EMAIL"    // Некорректный email
	ErrCodeInvalidPassword = "INVALID_PASSWORD" // Некорректный пароль
	ErrCodeInvalidLogin    = "INVALID_LOGIN"    // Некорректный логин

	// Ошибки существования
	ErrCodeUserExists   = "USER_EXISTS"   // Пользователь уже существует
	ErrCodeEmailExists  = "EMAIL_EXISTS"  // Email уже используется
	ErrCodeLoginExists  = "LOGIN_EXISTS"  // Логин уже занят
	ErrCodeWeakPassword = "WEAK_PASSWORD" // Пароль слишком слабый

	// Ошибки пользователя
	ErrCodeUserNotFound       = "USER_NOT_FOUND"      // Пользователь не найден
	ErrCodeInvalidCredentials = "INVALID_CREDENTIALS" // Неверные учетные данные
	ErrCodeAccountLocked      = "ACCOUNT_LOCKED"      // Аккаунт заблокирован

	// Ошибки токенов
	ErrCodeTokenExpired = "TOKEN_EXPIRED" // Токен истек
	ErrCodeTokenInvalid = "TOKEN_INVALID" // Недействительный токен
	ErrCodeTokenMissing = "TOKEN_MISSING" // Токен отсутствует

	// Ошибки авторизации
	ErrCodeUnauthorized = "UNAUTHORIZED"  // Требуется авторизация
	ErrCodeForbidden    = "FORBIDDEN"     // Доступ запрещен
	ErrCodeAccessDenied = "ACCESS_DENIED" // Доступ отклонен

	// Ошибки ресурсов
	ErrCodeResourceNotFound = "RESOURCE_NOT_FOUND" // Ресурс не найден
	ErrCodeResourceExists   = "RESOURCE_EXISTS"    // Ресурс уже существует
	ErrCodeResourceConflict = "RESOURCE_CONFLICT"  // Конфликт ресурсов
	ErrCodeBudgetNotFound   = "BUDGET_NOT_FOUND"   // Бюджет не найден
	ErrCodeAccountNotFound  = "ACCOUNT_NOT_FOUND"  // Счет не найден

	// Ошибки валидации данных
	ErrCodeInvalidAmount   = "INVALID_AMOUNT"   // Некорректная сумма
	ErrCodeInvalidCurrency = "INVALID_CURRENCY" // Некорректная валюта
	ErrCodeInvalidDate     = "INVALID_DATE"     // Некорректная дата
	ErrCodeInvalidPeriod   = "INVALID_PERIOD"   // Некорректный период

	// Системные ошибки
	ErrCodeInternalError      = "INTERNAL_ERROR"      // Внутренняя ошибка сервера
	ErrCodeServiceUnavailable = "SERVICE_UNAVAILABLE" // Сервис недоступен
	ErrCodeDatabaseError      = "DATABASE_ERROR"      // Ошибка базы данных
)

// ErrorDescriptions содержит описания всех кодов ошибок
var ErrorDescriptions = map[string]string{
	// Ошибки валидации
	ErrCodeInvalidRequest:  "Некорректный запрос",
	ErrCodeMissingFields:   "Отсутствуют обязательные поля",
	ErrCodeInvalidEmail:    "Некорректный email адрес",
	ErrCodeInvalidPassword: "Некорректный пароль",
	ErrCodeInvalidLogin:    "Некорректный логин",

	// Ошибки существования
	ErrCodeUserExists:   "Пользователь уже существует",
	ErrCodeEmailExists:  "Email уже используется",
	ErrCodeLoginExists:  "Логин уже занят",
	ErrCodeWeakPassword: "Пароль слишком слабый",

	// Ошибки пользователя
	ErrCodeUserNotFound:       "Пользователь не найден",
	ErrCodeInvalidCredentials: "Неверные учетные данные",
	ErrCodeAccountLocked:      "Аккаунт заблокирован",

	// Ошибки токенов
	ErrCodeTokenExpired: "Токен истек",
	ErrCodeTokenInvalid: "Недействительный токен",
	ErrCodeTokenMissing: "Токен отсутствует",

	// Ошибки авторизации
	ErrCodeUnauthorized: "Требуется авторизация",
	ErrCodeForbidden:    "Доступ запрещен",
	ErrCodeAccessDenied: "Доступ отклонен",

	// Ошибки ресурсов
	ErrCodeResourceNotFound: "Ресурс не найден",
	ErrCodeResourceExists:   "Ресурс уже существует",
	ErrCodeResourceConflict: "Конфликт ресурсов",
	ErrCodeBudgetNotFound:   "Бюджет не найден",
	ErrCodeAccountNotFound:  "Счет не найден",

	// Ошибки валидации данных
	ErrCodeInvalidAmount:   "Некорректная сумма",
	ErrCodeInvalidCurrency: "Некорректная валюта",
	ErrCodeInvalidDate:     "Некорректная дата",
	ErrCodeInvalidPeriod:   "Некорректный период",

	// Системные ошибки
	ErrCodeInternalError:      "Внутренняя ошибка сервера",
	ErrCodeServiceUnavailable: "Сервис недоступен",
	ErrCodeDatabaseError:      "Ошибка базы данных",
}
