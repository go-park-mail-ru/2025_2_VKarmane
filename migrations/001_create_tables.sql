-- Типы
CREATE TYPE OPERATION_STATUS AS ENUM ('finished', 'reverted');
CREATE TYPE OPERATION_TYPE AS ENUM ('expense', 'income');

-- ========================================================
-- Таблица USER
-- ========================================================
CREATE TABLE IF NOT EXISTS "user" (
    _id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_name TEXT NOT NULL CHECK (LENGTH(user_name) BETWEEN 3 AND 40),
    surname TEXT NOT NULL CHECK (LENGTH(surname) <= 40),
    email TEXT NOT NULL UNIQUE CHECK (LENGTH(email) <= 254 AND email ~ '^[^@]+@[^@]+$'),
    logo_hashed_id TEXT NOT NULL DEFAULT '6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b',
    user_login TEXT NOT NULL UNIQUE CHECK (LENGTH(user_login) BETWEEN 5 AND 40),
    user_hashed_password TEXT NOT NULL CHECK (LENGTH(user_hashed_password) > 0),
    user_description TEXT CHECK (LENGTH(user_description) <= 250),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- ========================================================
-- Таблица CURRENCY
-- ========================================================
CREATE TABLE IF NOT EXISTS currency (
    _id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    currency_name TEXT NOT NULL UNIQUE,
    logo_hashed_id TEXT NOT NULL DEFAULT 'e3b0c44298fc1c149afbf4c8996fb924',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- ========================================================
-- Таблица ACCOUNT
-- ========================================================
CREATE TABLE IF NOT EXISTS account (
    _id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    balance DECIMAL(10,2) NOT NULL DEFAULT 0.00 CHECK (balance >= 0),
    account_type TEXT NOT NULL DEFAULT 'default' CHECK (LENGTH(account_type) <= 30),
    currency_id INT REFERENCES currency(_id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- ========================================================
-- Таблица SHARINGS (USER_ACCOUNT)
-- ========================================================
CREATE TABLE IF NOT EXISTS sharings (
    _id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    account_id INT NOT NULL REFERENCES account(_id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES "user"(_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (account_id, user_id)
);

-- ========================================================
-- Таблица CATEGORY
-- ========================================================
CREATE TABLE IF NOT EXISTS category (
    _id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INT NOT NULL REFERENCES "user"(_id) ON DELETE CASCADE,
    category_name TEXT NOT NULL CHECK (LENGTH(category_name) <= 30),
    logo_hashed_id TEXT NOT NULL DEFAULT 'c1dfd96eea8cc2b62785275bca38ac261256e278',
    category_description TEXT CHECK (LENGTH(category_description) <= 60),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, category_name)
);

-- ========================================================
-- Таблица OPERATION
-- ========================================================
CREATE TABLE IF NOT EXISTS operation (
    _id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    account_from_id INT REFERENCES account(_id) ON DELETE SET NULL,
    account_to_id INT REFERENCES account(_id) ON DELETE SET NULL,
    category_id INT REFERENCES category(_id) ON DELETE SET NULL,
    currency_id INT REFERENCES currency(_id) ON DELETE SET NULL,
    operation_status OPERATION_STATUS DEFAULT 'finished',
    operation_type OPERATION_TYPE NOT NULL,
    operation_name TEXT NOT NULL CHECK (LENGTH(operation_name) <= 50),
    operation_description TEXT CHECK (LENGTH(operation_description) <= 100),
    receipt_url TEXT CHECK (LENGTH(receipt_url) <= 200),
    sum DECIMAL(10,2) NOT NULL CHECK (sum > 0),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- ========================================================
-- Таблица BUDGET
-- ========================================================
CREATE TABLE IF NOT EXISTS budget (
    _id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INT NOT NULL REFERENCES "user"(_id) ON DELETE CASCADE,
    category_id INT NOT NULL REFERENCES category(_id) ON DELETE CASCADE,
    currency_id INT NOT NULL REFERENCES currency(_id) ON DELETE CASCADE,
    amount DECIMAL(10,2) NOT NULL CHECK (amount >= 0),
    budget_description TEXT CHECK (LENGTH(budget_description) <= 80),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    closed_at TIMESTAMPTZ,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    UNIQUE (user_id, category_id, currency_id, period_start, period_end)
);

-- ========================================================
-- Таблицы чатов
-- ========================================================
CREATE TABLE IF NOT EXISTS chat (
    _id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    chat_name TEXT NOT NULL CHECK (LENGTH(chat_name) <= 30),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_in_chat (
    _id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INT NOT NULL REFERENCES "user"(_id) ON DELETE CASCADE,
    chat_id INT NOT NULL REFERENCES chat(_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, chat_id)
);

CREATE TABLE IF NOT EXISTS message (
    _id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INT NOT NULL REFERENCES "user"(_id) ON DELETE CASCADE,
    chat_id INT NOT NULL REFERENCES chat(_id) ON DELETE CASCADE,
    message_text TEXT NOT NULL CHECK (LENGTH(message_text) <= 500),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- ========================================================
-- Таблица RECEIVER
-- ========================================================
CREATE TABLE IF NOT EXISTS receiver (
    _id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INT NOT NULL REFERENCES "user"(_id) ON DELETE CASCADE,
    receiver_name TEXT NOT NULL CHECK (LENGTH(receiver_name) <= 60),
    logo_hashed_id TEXT NOT NULL DEFAULT 'b6d81b360a5672d80c27430f39153e2c',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, receiver_name)
);

-- ========================================================
-- Функция обновления updated_at
-- ========================================================
CREATE OR REPLACE FUNCTION public.moddatetime()
RETURNS trigger AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ========================================================
-- Триггеры обновления
-- ========================================================
CREATE TRIGGER modify_user_updated_at
    BEFORE UPDATE ON "user"
    FOR EACH ROW
    EXECUTE PROCEDURE public.moddatetime();

CREATE TRIGGER modify_account_updated_at
    BEFORE UPDATE ON account
    FOR EACH ROW
    EXECUTE PROCEDURE public.moddatetime();

CREATE TRIGGER modify_sharings_updated_at
    BEFORE UPDATE ON sharings
    FOR EACH ROW
    EXECUTE PROCEDURE public.moddatetime();

CREATE TRIGGER modify_category_updated_at
    BEFORE UPDATE ON category
    FOR EACH ROW
    EXECUTE PROCEDURE public.moddatetime();

CREATE TRIGGER modify_budget_updated_at
    BEFORE UPDATE ON budget
    FOR EACH ROW
    EXECUTE PROCEDURE public.moddatetime();

CREATE TRIGGER modify_chat_updated_at
    BEFORE UPDATE ON chat
    FOR EACH ROW
    EXECUTE PROCEDURE public.moddatetime();

CREATE TRIGGER modify_user_in_chat_updated_at
    BEFORE UPDATE ON user_in_chat
    FOR EACH ROW
    EXECUTE PROCEDURE public.moddatetime();

CREATE TRIGGER modify_message_updated_at
    BEFORE UPDATE ON message
    FOR EACH ROW
    EXECUTE PROCEDURE public.moddatetime();

CREATE TRIGGER modify_receiver_updated_at
    BEFORE UPDATE ON receiver
    FOR EACH ROW
    EXECUTE PROCEDURE public.moddatetime();
