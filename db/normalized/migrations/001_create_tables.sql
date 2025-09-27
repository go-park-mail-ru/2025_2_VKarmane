-- 001_create_tables.sql

CREATE TYPE OPERATION_STATUS AS ENUM ('finished', 'reverted')
CREATE TYPE OPERATION_TYPE AS ENUM ('transfer', 'expeense', 'income')

--Таблица USER

CREATE TABLE IF NOT EXISTS user (
    _id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_name TEXT NOT NULL CHECK (LENGTH(user_name) <= 40 AND LENGTH(user_name) >= 3),
    surname TEXT NOT NULL CHECK (LENGTH(surname) <= 40),
    email TEXT NOT NULL UNIQUE CHECK (LENGTH(email) <= 254),
    logo_hasehd_id TEXT NOT NULL DEFAULT 1,
    user_login TEXT NOT NULL UNIQUE CHECK (LENGTH(user_login) <= 40 AND LENGTH(user_login) >= 5),
    user_hashed_password TEXT NOT NULL
    user_description TEXT CHECK (LENGTH(user_description) <= 250)
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
)


--Таблица ACCOUNT

CREATE TABLE IF NOT EXISTS account (
    _id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    balance DECIMAL(10, 2) NOT NULL DEFAULT 0.00 check (balance >= 0),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
)

--Таблица SHARINGS

CREATE TABLE IF NOT EXISTS sharings (
    _id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    account_id INT NOT NULL REFERENCES account(_id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES user(_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (account_id, user_id)
)

--Таблица CURRENCY

CREATE TABLE IF NOT EXISTS currency (
    _id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    currency_name TEXT NOT NULL UNIQUE,
    logo_hasehd_id TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
)


--Таблица CATEGORY

CREATE TABLE IF NOT EXISTS category (
    _id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INT REFERENCES USER(_id) ON DELETE CASCADE,
    category_name TEXT NOT NULL CHECK (LENGTH(category_name) <= 15),
    logo_hasehd_id TEXT NOT NULL DEFAULT 2
    category_description TEXT CHECK (LENGTH(category_description) <= 30)
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, category_name)
)

--Таблица OPERATION

CREATE TABLE IF NOT EXISTS operation (
    _id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    account_id INT NOT NULL REFERENCES account(_id) ON DELETE CASCADE,
    category_id INT NOT NULL REFERENCES category(_id) ON SET DEFAULT 1,
    operation_status OPERATION_STATUS, 
    operation_type OPERATION_TYPE NOT NULL,
    operation_name TEXT NOT NULL CHECK (LENGTH(operation_name) <= 50)
    operation_description TEXT CHECK (LENGTH(operation_description) <= 60)
    receipt_url TEXT CHECK (LENGTH(receipt_url) <= 100),
    sum DECIMAL(10, 2) NOT NULL CHECK (sum > 0),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
)

--Таблица TRANSFER

CREATE TABLE IF NOT EXISTS transfer_op (
    _id PRIMARY KEY REFERENCES operation(_id),
    from_account_id INT NOT NULL REFERENCES account(_id) ON DELETE CASCADE
)

--Таблица BUDGET

CREATE TABLE IF NOT EXISTS budget (
    _id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INT NOT NULL REFERENCES user(_id) ON DELETE CASCADE,
    category_id INT NOT NULL REFERENCES category(_id) ON DELETE CASCADE,
    currency_id INT NOT NULL REFERENCES currency(_id) ON DELETE CASCADE,
    amount DECIMAL(10, 2) NOT NULL CHECK (amount >= 0),
    budget_description TEXT CHECK (LENGTH(budget_description) <= 50)
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    closed_at TIMESTAMPTZ,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    UNIQUE (user_id, category_id, currency_id, period_start, period_end)
)

CREATE TABLE IF NOT EXISTS chat (
    _id GENERATED ALWAYS AS IDENTITY PRIMARY KEY
    chat_name TEXT NOT NULL CHECK (LENGTH(chat_name) <= 30)
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
)

CREATE TABLE IF NOT EXISTS dialogue (
    _id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY
    user_id INT NOT NULL REFERENCES user(_id) ON DELETE CASCADE
    chat_id INT NOT NULL REFERENCES chat(_id) ON DELETE CASCADE
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, chat_id)
)

CREATE TABLE IF NOT EXISTS message  (
    _id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY
    user_id INT NOT NULL REFERENCES user(_id) ON DELETE CASCADE
    chat_id INT NOT NULL REFERENCES chat(_id) ON DELETE CASCADE
    message_text TEXT NOT NULL CHECK (LENGTH(message_text) <= 200)
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
)

CREATE TABLE IF NOT EXISTS receiver (
    _id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY
    user_id INT NOT NULL REFERENCES user(_id) ON DELETE CASCADE
    receiver_name TEXT NOT NULL CHECK (LENGTH(receiver_name) <= 60)
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, receiver_name)
)

CREATE TRIGGER modify_payment_updated_at
    BEFORE UPDATE
    ON user, account, sharings, category, budget, chat, dialogue, message, receiver
    FOR EACH ROW
EXECUTE PROCEDURE public.moddatetime(updated_at);



