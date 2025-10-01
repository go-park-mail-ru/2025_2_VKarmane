# API Документация

## Аутентификация

Сервис использует JWT токены для аутентификации через HTTP-only secure cookies. Все защищенные маршруты автоматически получают токен из cookie `auth_token`.

**Настройки cookie:**
- HttpOnly: `true` - защита от XSS атак
- Secure: `true` в продакшене - только HTTPS
- SameSite: `Strict` - защита от CSRF атак
- MaxAge: 86400 секунд (24 часа)

## Структура данных

### Счета (Accounts)
- `account_id` - Уникальный идентификатор счета
- `balance` - Текущий баланс счета
- `type` - Тип счета (`card` - карта, `cash` - наличные)
- `currency_id` - Идентификатор валюты

### Бюджеты (Budgets)
- `budget_id` - Уникальный идентификатор бюджета
- `user_id` - Идентификатор пользователя-владельца
- `amount` - Запланированная сумма бюджета
- `actual` - Фактически потраченная сумма
- `currency_id` - Идентификатор валюты
- `description` - Описание бюджета
- `period_start` - Начало периода (ISO 8601)
- `period_end` - Конец периода (ISO 8601)

## Маршруты

### Публичные маршруты (без аутентификации)

#### POST /api/v1/auth/register
Регистрация нового пользователя.

**Тело запроса:**
```json
{
  "first_name": "string (обязательно, 2-50 символов)",
  "last_name": "string (обязательно, 2-50 символов)", 
  "email": "string (обязательно, валидный email)",
  "login": "string (обязательно, 3-30 символов, только буквы и цифры)",
  "password": "string (обязательно, 6-100 символов)"
}
```

**Ответ:**
```json
{
  "user": {
    "ID": 3,
    "FirstName": "John",
    "LastName": "Doe",
    "Email": "newuser@example.com",
    "Login": "newuser",
    "Password": "",
    "CreatedAt": "2025-09-26T18:55:49.319905+03:00",
    "UpdatedAt": "2025-09-26T18:55:49.319905+03:00"
  }
}
```

**Cookie:** `auth_token` устанавливается автоматически

#### POST /api/v1/auth/login
Вход в систему.

**Тело запроса:**
```json
{
  "login": "string (обязательно, 3-30 символов)",
  "password": "string (обязательно, минимум 6 символов)"
}
```

**Ответ:**
```json
{
  "user": {
    "ID": 1,
    "FirstName": "Vlad",
    "LastName": "Sigma", 
    "Email": "vlad@example.com",
    "Login": "hello",
    "Password": "",
    "CreatedAt": "2025-09-26T18:46:15.000000+03:00",
    "UpdatedAt": "2025-09-26T18:46:15.000000+03:00"
  }
}
```

**Cookie:** `auth_token` устанавливается автоматически

#### POST /api/v1/auth/logout
Выход из системы.

**Ответ:**
```json
{
  "message": "Logged out successfully"
}
```

**Cookie:** `auth_token` удаляется

### Защищенные маршруты (требуют аутентификации)

#### GET /api/v1/profile
Получение профиля текущего пользователя.

**Ответ:**
```json
{
  "ID": 1,
  "FirstName": "Vlad",
  "LastName": "Sigma",
  "Email": "vlad@example.com", 
  "Login": "hello",
  "Password": "",
  "CreatedAt": "2025-09-26T18:46:15.000000+03:00",
  "UpdatedAt": "2025-09-26T18:46:15.000000+03:00"
}
```

#### GET /api/v1/balance
Получение баланса пользователя.

**Поля ответа:**
- `user_id` - ID пользователя
- `accounts` - Массив счетов пользователя
  - `account_id` - ID счета
  - `balance` - Баланс счета
  - `type` - Тип счета (card, cash)
  - `currency_id` - ID валюты

**Ответ:**
```json
{
  "user_id": 1,
  "accounts": [
    {
      "account_id": 1,
      "balance": 100,
      "type": "card",
      "currency_id": 1
    },
    {
      "account_id": 2,
      "balance": 500,
      "type": "cash",
      "currency_id": 1
    }
  ]
}
```

#### GET /api/v1/balance/{id}
Получение баланса конкретного счета.

**Параметры:**
- `id` - ID счета

**Ответ:**
```json
{
  "account_id": 1,
  "balance": 100,
  "type": "card",
  "currency_id": 1
}
```

#### GET /api/v1/budgets
Получение списка бюджетов пользователя.

**Поля ответа:**
- `user_id` - ID пользователя
- `budgets` - Массив бюджетов пользователя
  - `budget_id` - ID бюджета
  - `user_id` - ID пользователя-владельца
  - `amount` - Запланированная сумма бюджета
  - `actual` - Фактически потраченная сумма
  - `currency_id` - ID валюты
  - `description` - Описание бюджета
  - `period_start` - Начало периода (ISO 8601)
  - `period_end` - Конец периода (ISO 8601)

**Ответ:**
```json
{
  "user_id": 1,
  "budgets": [
    {
      "budget_id": 1,
      "user_id": 1,
      "amount": 100,
      "actual": 110,
      "currency_id": 1,
      "description": "September food",
      "period_start": "2025-09-01T00:00:00Z",
      "period_end": "2025-09-30T00:00:00Z"
    },
    {
      "budget_id": 2,
      "user_id": 1,
      "amount": 500,
      "actual": 110,
      "currency_id": 1,
      "description": "September relax",
      "period_start": "2025-09-01T00:00:00Z",
      "period_end": "2025-09-30T00:00:00Z"
    }
  ]
}
```

#### GET /api/v1/budget/{id}
Получение конкретного бюджета.

**Параметры:**
- `id` - ID бюджета

**Поля ответа:**
- `budget_id` - ID бюджета
- `user_id` - ID пользователя-владельца
- `amount` - Запланированная сумма бюджета
- `actual` - Фактически потраченная сумма
- `currency_id` - ID валюты
- `description` - Описание бюджета
- `period_start` - Начало периода (ISO 8601)
- `period_end` - Конец периода (ISO 8601)

**Ответ:**
```json
{
  "budget_id": 1,
  "user_id": 1,
  "amount": 100,
  "actual": 110,
  "currency_id": 1,
  "description": "September food",
  "period_start": "2025-09-01T00:00:00Z",
  "period_end": "2025-09-30T00:00:00Z"
}
```

## Коды ошибок

### HTTP статус коды
- `400 Bad Request` - Неверные данные запроса
- `401 Unauthorized` - Требуется аутентификация
- `404 Not Found` - Ресурс не найден
- `409 Conflict` - Конфликт данных
- `500 Internal Server Error` - Внутренняя ошибка сервера

### Структура ошибки
Все ошибки возвращаются в следующем формате:
```json
{
  "error": "Человекочитаемое сообщение об ошибке",
  "code": "КОД_ОШИБКИ",
  "details": "Дополнительные детали",
  "field": "Поле с ошибкой",
  "timestamp": "2025-01-27T10:30:00Z"
}
```

### Коды ошибок аутентификации
- `INVALID_REQUEST` - Некорректный запрос
- `MISSING_FIELDS` - Отсутствуют обязательные поля
- `INVALID_EMAIL` - Некорректный email
- `INVALID_PASSWORD` - Некорректный пароль
- `INVALID_LOGIN` - Некорректный логин
- `USER_EXISTS` - Пользователь уже существует
- `EMAIL_EXISTS` - Email уже используется
- `LOGIN_EXISTS` - Логин уже занят
- `WEAK_PASSWORD` - Пароль слишком слабый
- `USER_NOT_FOUND` - Пользователь не найден
- `INVALID_CREDENTIALS` - Неверные учетные данные
- `UNAUTHORIZED` - Требуется авторизация
- `TOKEN_EXPIRED` - Токен истек
- `TOKEN_INVALID` - Недействительный токен
- `TOKEN_MISSING` - Токен отсутствует

### Коды ошибок ресурсов
- `RESOURCE_NOT_FOUND` - Ресурс не найден
- `RESOURCE_EXISTS` - Ресурс уже существует
- `RESOURCE_CONFLICT` - Конфликт ресурсов
- `BUDGET_NOT_FOUND` - Бюджет не найден
- `ACCOUNT_NOT_FOUND` - Счет не найден

### Примеры ошибок

#### Ошибка валидации
```json
{
  "error": "Имя обязательно для заполнения",
  "code": "MISSING_FIELDS",
  "details": "",
  "field": "first_name",
  "timestamp": "2025-01-27T10:30:00Z"
}
```

#### Ошибка конфликта
```json
{
  "error": "Пользователь с таким email уже существует",
  "code": "EMAIL_EXISTS",
  "details": "",
  "field": "email",
  "timestamp": "2025-01-27T10:30:00Z"
}
```

#### Ошибка авторизации
```json
{
  "error": "Пользователь не найден",
  "code": "USER_NOT_FOUND",
  "details": "",
  "field": "",
  "timestamp": "2025-01-27T10:30:00Z"
}
```

#### Ошибки валидации
```json
{
  "error": "Ошибки валидации",
  "code": "INVALID_REQUEST",
  "details": "Поле first_name должно содержать минимум 2 символов; Поле email должно содержать корректный email адрес",
  "timestamp": "2025-01-27T10:30:00Z",
  "validation_errors": [
    {
      "field": "first_name",
      "tag": "min",
      "value": "J",
      "message": "Поле first_name должно содержать минимум 2 символов"
    },
    {
      "field": "email",
      "tag": "email",
      "value": "invalid-email",
      "message": "Поле email должно содержать корректный email адрес"
    }
  ]
}
```

## Тестовые пользователи

Для разработки и тестирования доступны следующие пользователи:

| Login | Password | Имя | Email |
|-------|----------|-----|-------|
| `hello` | `Test123` | Vlad Sigma | vlad@example.com |
| `goodbye` | `Test123` | Nikita Go | nikita@example.com |

## Примеры использования

### Регистрация и вход
```bash
# Регистрация нового пользователя
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "newuser@example.com",
    "login": "newuser",
    "password": "password123"
  }' \
  -c cookies.txt

# Вход существующего пользователя (из inmemory данных)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "login": "hello",
    "password": "Test123"
  }' \
  -c cookies.txt
```

### Использование защищенных маршрутов
```bash
# Получение профиля
curl -X GET http://localhost:8080/api/v1/profile \
  -b cookies.txt

# Получение баланса
curl -X GET http://localhost:8080/api/v1/balance \
  -b cookies.txt

# Получение конкретного счета
curl -X GET http://localhost:8080/api/v1/balance/1 \
  -b cookies.txt

# Получение списка бюджетов
curl -X GET http://localhost:8080/api/v1/budgets \
  -b cookies.txt

# Получение конкретного бюджета
curl -X GET http://localhost:8080/api/v1/budget/1 \
  -b cookies.txt

# Выход из системы
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -b cookies.txt
```
