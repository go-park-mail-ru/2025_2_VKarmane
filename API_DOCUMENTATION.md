# API Документация

## Аутентификация

Сервис использует JWT токены для аутентификации. Все защищенные маршруты требуют заголовок `Authorization: Bearer <token>`.

## Маршруты

### Публичные маршруты (без аутентификации)

#### POST /api/v1/auth/register
Регистрация нового пользователя.

**Тело запроса:**
```json
{
  "email": "string (валидный email)",
  "login": "string (3-30 символов)",
  "password": "string (минимум 6 символов)"
}
```

**Ответ:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "ID": 3,
    "Name": "",
    "Surname": "",
    "Email": "newuser@example.com",
    "Login": "newuser",
    "Password": "",
    "CreatedAt": "2025-09-26T18:55:49.319905+03:00",
    "UpdatedAt": "2025-09-26T18:55:49.319905+03:00"
  }
}
```

#### POST /api/v1/auth/login
Вход в систему.

**Тело запроса:**
```json
{
  "login": "string",
  "password": "string"
}
```

**Ответ:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "ID": 1,
    "Name": "Vlad",
    "Surname": "Sigma", 
    "Email": "vlad@example.com",
    "Login": "hello",
    "Password": "",
    "CreatedAt": "2025-09-26T18:46:15.000000+03:00",
    "UpdatedAt": "2025-09-26T18:46:15.000000+03:00"
  }
}
```

### Защищенные маршруты (требуют аутентификации)

#### GET /api/v1/profile
Получение профиля текущего пользователя.

**Заголовки:**
```
Authorization: Bearer <token>
```

**Ответ:**
```json
{
  "ID": 1,
  "Name": "Vlad",
  "Surname": "Sigma",
  "Email": "vlad@example.com", 
  "Login": "hello",
  "Password": "",
  "CreatedAt": "2025-09-26T18:46:15.000000+03:00",
  "UpdatedAt": "2025-09-26T18:46:15.000000+03:00"
}
```

#### GET /api/v1/balance
Получение баланса пользователя.

**Заголовки:**
```
Authorization: Bearer <token>
```

**Ответ:**
```json
{
  "user_id": 1,
  "accounts": [
    {
      "id": 1,
      "balance": 100,
      "type": "card",
      "currency": {
        "id": 1,
        "code": "USD",
        "name": "US Dollar"
      }
    },
    {
      "id": 2,
      "balance": 500,
      "type": "cash",
      "currency": {
        "id": 1,
        "code": "USD",
        "name": "US Dollar"
      }
    }
  ]
}
```

#### GET /api/v1/balance/{id}
Получение баланса конкретного счета.

**Заголовки:**
```
Authorization: Bearer <token>
```

**Параметры:**
- `id` - ID счета

**Ответ:**
```json
{
  "id": 1,
  "balance": 100,
  "type": "card",
  "currency": {
    "id": 1,
    "code": "USD",
    "name": "US Dollar"
  }
}
```

#### GET /api/v1/budgets
Получение списка бюджетов пользователя.

**Заголовки:**
```
Authorization: Bearer <token>
```

**Ответ:**
```json
{
  "user_id": 1,
  "budgets": [
    {
      "id": 1,
      "amount": 100,
      "description": "September food",
      "currency": {
        "id": 1,
        "code": "USD",
        "name": "US Dollar"
      },
      "period_start": "2025-09-01T00:00:00Z",
      "period_end": "2025-09-30T00:00:00Z"
    },
    {
      "id": 2,
      "amount": 500,
      "description": "September relax",
      "currency": {
        "id": 1,
        "code": "USD",
        "name": "US Dollar"
      },
      "period_start": "2025-09-01T00:00:00Z",
      "period_end": "2025-09-30T00:00:00Z"
    }
  ]
}
```

#### GET /api/v1/budget/{id}
Получение конкретного бюджета.

**Заголовки:**
```
Authorization: Bearer <token>
```

**Параметры:**
- `id` - ID бюджета

**Ответ:**
```json
{
  "id": 1,
  "amount": 100,
  "description": "September food",
  "currency": {
    "id": 1,
    "code": "USD",
    "name": "US Dollar"
  },
  "period_start": "2025-09-01T00:00:00Z",
  "period_end": "2025-09-30T00:00:00Z"
}
```

## Коды ошибок

- `400 Bad Request` - Неверные данные запроса
- `401 Unauthorized` - Требуется аутентификация или неверный токен
- `404 Not Found` - Ресурс не найден
- `500 Internal Server Error` - Внутренняя ошибка сервера

## Примеры использования

### Регистрация и вход
```bash
# Регистрация нового пользователя
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "login": "newuser",
    "password": "password123"
  }'

# Вход существующего пользователя (из inmemory данных)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "login": "hello",
    "password": "password123"
  }'
```

### Использование защищенных маршрутов
```bash
# Получение профиля
curl -X GET http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer <your-jwt-token>"

# Получение баланса
curl -X GET http://localhost:8080/api/v1/balance \
  -H "Authorization: Bearer <your-jwt-token>"
```
