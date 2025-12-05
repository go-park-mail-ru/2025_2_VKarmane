# Скрипты управления базой данных

## Скрипт создания сервисного пользователя

### `create_service_user.sql`

```bash
# Подключитесь к БД под администратором
psql -U postgres -d vkarmane -f scripts/db/create_service_user.sql

# Или если используете пользователя vkarmane
psql -U vkarmane -d vkarmane -f scripts/db/create_service_user.sql
```

### Права сервисного пользователя

Пользователь `vkarmane_app` имеет следующие права:

- **user**: SELECT, INSERT, UPDATE
- **account**: SELECT, INSERT, UPDATE, DELETE
- **sharings**: SELECT, INSERT, UPDATE, DELETE
- **category**: SELECT, INSERT, UPDATE, DELETE
- **operation**: SELECT, INSERT, UPDATE
- **budget**: SELECT, INSERT, UPDATE
- **currency**: SELECT (только чтение, справочная таблица)
- **Функции**: EXECUTE на `moddatetime()` (для триггеров)
- **Типы**: USAGE на `OPERATION_STATUS` и `OPERATION_TYPE`
- **Последовательности**: USAGE, SELECT на все последовательности таблиц
