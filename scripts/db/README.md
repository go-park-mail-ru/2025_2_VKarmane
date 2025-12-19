# Скрипты управления базой данных

## Скрипт создания сервисного пользователя

### `create_service_user.sql`

```bash
# Подключитесь к БД под суперпользователем
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


## Скрипт создания учетной записи администратора

```bash
# Подключитесь к БД под суперпользователем (или ролевым администратором)
psql -U postgres -d vkarmane -f scripts/db/create_admin_user.sql
```

### Права администратора `vkarmane_admin`

- **Роль**: LOGIN, CREATEDB, CREATEROLE, INHERIT, без SUPERUSER
- **База**: CONNECT на `vkarmane`
- **Схема `public`**:
  - USAGE, CREATE
  - ALL PRIVILEGES на все существующие таблицы, последовательности и функции
  - ALL PRIVILEGES по умолчанию на будущие таблицы, последовательности и функции

## Настройка `pg_ident.conf`

Файл `pg_ident.conf` используется для маппинга системных пользователей операционной системы в роли PostgreSQL. Это необходимо для **peer authentication** (когда PostgreSQL проверяет системного пользователя ОС).

Чтобы использовать peer authentication для админа, нужно создать системного пользователя в контейнере:

```bash
# Войдите в контейнер PostgreSQL
docker exec -it vkarmane-postgres sh

# Создайте системного пользователя
adduser -D -s /bin/sh vkarmane_admin_os

# Выйдите из контейнера
exit
```

После этого можно подключаться локально:
```bash
docker exec -it --user vkarmane_admin_os vkarmane-postgres psql -U vkarmane_admin -d vkarmane
```
