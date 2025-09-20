-- 001_create_tables.sql

--Таблица USER

CREATE TABLE IF NOT EXISTS user (
    _id INT GENERATED ALWAYS BY IDENTITY PRIMARY KEY
    user_name TEXT NOT NULL CHECK (LENGTH(user_name) <= 40)
    surname TEXT NOT NULL CHECK (LENGTH(surname) <= 40)
    email TEXT NOT NULL UNIQUE CHECK (LENGTH(email) <= 25)
    logo text NOT NULL DEFAULT '.../default_user.png'
    user_login TEXT NOT NULL UNIQUE CHECK (LENGTH(user_login) <= 40 AND LENGTH(user_login) >= 5)
    user_password TEXT NOT NULL CHECK (LENGTH(user_password) <= 40 AND LENGTH(user_password) >= 10)
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
)


--Таблица ACCOUNT

CREATE TABLE IF NOT EXISTS account (
    _id INT GENERATED ALWAYS BY IDENTITY PRIMARY KEY
    balance DECIMAL(10, 2) NOT NULL DEFAULT 0.00 check (balance >= 0)
    account_type TEXT NOT NULL
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
)

--Таблица SHARINGS

CREATE TABLE IF NOT EXISTS sharings (
    _id INT GENERATED ALWAYS BY IDENTITY PRIMARY KEY
    account_id INT NOT NULL REFERENCES account(_id) ON DELETE CASCADE
    user_id INT NOT NULL REFERENCES user(_id) ON DELETE CASCADE
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    UNIQUE (account_id, user_id)
)

--Таблица CURRENCY

CREATE TABLE IF NOT EXISTS currency (
    _id INT GENERATED ALWAYS BY IDENTITY PRIMARY KEY
    code TEXT NOT NULL UNIQUE
    currency_name TEXT NOT NULL UNIQUE
    logo TEXT NOT NULL
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
)


--Таблица CATEGORY

CREATE TABLE IF NOT EXISTS category (
    _id INT GENERATED ALWAYS BY IDENTITY PRIMARY KEY
    user_id INT REFERENCES USER(_id) ON DELETE CASCADE
    category_name TEXT NOT NULL CHECK (LENGTH(user_name) <= 15)
    logo TEXT DEFAULT '.../default_category.png'
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    UNIQUE (user_id, category_name)
)

--Таблица OPERATION

CREATE TABLE IF NOT EXISTS operation (
    _id INT GENERATED ALWAYS BY IDENTITY PRIMARY KEY
    account_id INT NOT NULL REFERENCES account(_id) ON DELETE CASCADE
    category_id INT NOT NULL REFERENCES category(_id) ON DELETE CASCADE
    operation_type TEXT NOT NULL
    operation_name TEXT NOT NULL
    sum DECIMAL(10, 2) NOT NULL CHECK (sum > 0)
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
)

--Таблица TRANSFER

CREATE TABLE IF NOT EXISTS transfer_op (
    _id PRIMARY KEY REFERENCES operation(_id)
    from_account_id INT NOT NULL REFERENCES account(_id) ON DELETE CASCADE
)

--Таблица BUDGET

CREATE TABLE IF NOT EXISTS budget (
    _id INT GENERATED ALWAYS BY IDENTITY PRIMARY KEY
    user_id INT NOT NULL REFERENCES user(_id) ON DELETE CASCADE
    category_id INT NOT NULL REFERENCES category(_id) ON DELETE CASCADE
    currency_id INT NOT NULL REFERENCES currency(_id) ON DELETE CASCADE
    amount DECIMAL(10, 2) NOT NULL CHECK (amount >= 0)
    budget_type TEXT NOT NULL
    is_failed BOOLEAN NOT NULL DEFAULT FALSE
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    closed_at TIMESTAMPTZ
    period_start DATE NOT NULL
    period_end DATE NOT NULL
    UNIQUE (user_id, category_id, currency_id, period_start, period_end)
)
