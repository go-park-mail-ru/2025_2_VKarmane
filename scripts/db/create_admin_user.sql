-- Создание учетной записи администратора для VKarmane

DO
$$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_roles WHERE rolname = 'vkarmane_admin'
    ) THEN
        CREATE ROLE vkarmane_admin
            LOGIN
            PASSWORD 'vkarmane_admin_password'
            CREATEDB
            CREATEROLE
            INHERIT
            NOSUPERUSER
            NOBYPASSRLS
            NOREPLICATION;
    END IF;
END
$$;

-- Доступ к базе vkarmane
GRANT CONNECT ON DATABASE vkarmane TO vkarmane_admin;

-- Доступ к схеме public
GRANT USAGE, CREATE ON SCHEMA public TO vkarmane_admin;

-- Полный доступ к существующим таблицам/последовательностям/функциям схемы public
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO vkarmane_admin;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO vkarmane_admin;
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO vkarmane_admin;

-- Доступ по умолчанию к объектам, созданным в будущем
ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT ALL PRIVILEGES ON TABLES TO vkarmane_admin;
ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT ALL PRIVILEGES ON SEQUENCES TO vkarmane_admin;
ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT ALL PRIVILEGES ON FUNCTIONS TO vkarmane_admin;

COMMENT ON ROLE vkarmane_admin IS 'Административная учетная запись для администрирования БД VKarmane (без SUPERUSER).';


