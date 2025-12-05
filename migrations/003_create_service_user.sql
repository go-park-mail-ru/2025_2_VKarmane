-- ========================================================
-- Миграция: Создание сервисной учетной записи для приложения
-- ========================================================
-- 
-- Этот скрипт создает отдельного пользователя для работы приложения с БД.
--
-- ========================================================

-- Создание сервисного пользователя (если еще не существует)

DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = 'vkarmane_app') THEN
        CREATE USER vkarmane_app WITH PASSWORD 'vkarmane_app_password';
    END IF;
END
$$;

-- Предоставление прав на подключение к базе данных
GRANT CONNECT ON DATABASE vkarmane TO vkarmane_app;

-- Предоставление прав на использование схемы public
GRANT USAGE ON SCHEMA public TO vkarmane_app;

-- ========================================================
-- Права на таблицы
-- ========================================================

-- Таблица user: SELECT, INSERT, UPDATE (приложение не удаляет пользователей)
GRANT SELECT, INSERT, UPDATE ON TABLE public."user" TO vkarmane_app;
GRANT USAGE, SELECT ON SEQUENCE user__id_seq TO vkarmane_app;

-- Таблица account: SELECT, INSERT, UPDATE, DELETE
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE public.account TO vkarmane_app;
GRANT USAGE, SELECT ON SEQUENCE account__id_seq TO vkarmane_app;

-- Таблица sharings: SELECT, INSERT, UPDATE, DELETE
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE public.sharings TO vkarmane_app;
GRANT USAGE, SELECT ON SEQUENCE sharings__id_seq TO vkarmane_app;

-- Таблица category: SELECT, INSERT, UPDATE, DELETE
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE public.category TO vkarmane_app;
GRANT USAGE, SELECT ON SEQUENCE category__id_seq TO vkarmane_app;

-- Таблица operation: SELECT, INSERT, UPDATE (удаление через UPDATE статуса)
GRANT SELECT, INSERT, UPDATE ON TABLE public.operation TO vkarmane_app;
GRANT USAGE, SELECT ON SEQUENCE operation__id_seq TO vkarmane_app;

-- Таблица budget: SELECT, INSERT, UPDATE (удаление через UPDATE closed_at)
GRANT SELECT, INSERT, UPDATE ON TABLE public.budget TO vkarmane_app;
GRANT USAGE, SELECT ON SEQUENCE budget__id_seq TO vkarmane_app;

-- Таблица currency: только SELECT (справочная таблица, только чтение)
GRANT SELECT ON TABLE public.currency TO vkarmane_app;
GRANT USAGE, SELECT ON SEQUENCE currency__id_seq TO vkarmane_app;

-- ========================================================
-- Права на функции
-- ========================================================

-- Функция moddatetime() используется триггерами для обновления updated_at
GRANT EXECUTE ON FUNCTION public.moddatetime() TO vkarmane_app;

-- ========================================================
-- Права на типы (ENUM)
-- ========================================================

-- Типы OPERATION_STATUS и OPERATION_TYPE используются в таблице operation
GRANT USAGE ON TYPE public.OPERATION_STATUS TO vkarmane_app;
GRANT USAGE ON TYPE public.OPERATION_TYPE TO vkarmane_app;

-- ========================================================
-- Комментарии для документации
-- ========================================================

COMMENT ON ROLE vkarmane_app IS 'Сервисная учетная запись для приложения VKarmane. Имеет минимальные права, необходимые для работы приложения.';

