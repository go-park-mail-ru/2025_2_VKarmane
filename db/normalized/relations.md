## Таблица USER

Таблица `USER` содержит информацию о пользователях сервиса.

<p> Функциональные зависимости: </p>

- `{id} -> {name, surname, email, logo, login, password, created_at, updated_at}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты id, name, surname, email, logo, login, password, created_at, updated_at являются атомарными.
- 2 НФ: Атрибуты id, name, surname, email, logo, login, password, created_at, updated_at полностью функционально зависят от первичного ключа id.
- 3 НФ: Атрибуты id, name, surname, email, logo, login, password, created_at, updated_at не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    USER {
        int id PK
        string name 
        string surname 
        string email 
        string logo
        string login 
        string password 
        timestamptz
        created_at 
        timestamptz updated_at 
    }
```
## Таблица ACCOUNT

Таблица `ACCOUNT` содержит информацию о счете пользователя.

<p> Функциональные зависимости: </p>

- `{id} -> {balance, type, created_at, updated_at}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты id, balance, type, created_at, updated_at являются атомарными.
- 2 НФ: Атрибуты id, balance, type, created_at, updated_at полностью функционально зависят от первичного ключа id.
- 3 НФ: Атрибуты id, balance, type, created_at, updated_at не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
     ACCOUNT {
        int id PK
        decimal balance
        string type
        timestamptz created_at
        timestamptz updated_at
    }
```
## Таблица SHARINGS

Таблица `SHARINGS` является связующей между пользователем и счетом, позволяя таким образои нескольким пользователям иметь один счет.

<p> Функциональные зависимости: </p>

- `{id} -> {user_id, account_id, created_at, updated_at}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты id, user_id, account_id, created_at, updated_at являются атомарными.
- 2 НФ: Атрибуты id, user_id, account_id, created_at, updated_at полностью функционально зависят от первичного ключа id.
- 3 НФ: Атрибуты id, user_id, account_id, created_at, updated_at не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    SHARINGS{ 
        int id PK
        int user_id FK
        int account_id FK
        timestamptz created_at
        timestamptz updated_at
        }
```
## Таблица CURRENCY

Таблица `CURRENCY` содержит информацию о валюте.

<p> Функциональные зависимости: </p>

- `{id} -> {code, name, logo, created_at}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты id, code, name, logo, created_at являются атомарными.
- 2 НФ: Атрибуты id, code, name, logo, created_at полностью функционально зависят от первичного ключа id.
- 3 НФ: Атрибуты id, code, name, logo, created_at не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram
    CURRENCY {
        int id PK
        string code
        string name
        string logo
        timestamptz created_at
    }
```
## Таблица BUDGET

Таблица `BUDGET` содержит информацию о бюджетах пользователя.

<p> Функциональные зависимости: </p>

- `{id} -> {user_id, amount, currency_id, type, is_failed, created_at, updated_at, closed_at, period_start, period_end}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты id, user_id, amount, currency_id, type, is_failed, created_at, updated_at, closed_at, period_start, period_end являются атомарными.
- 2 НФ: Атрибуты id, user_id, amount, currency_id, type, is_failed, created_at, updated_at, closed_at, period_start, period_end полностью функционально зависят от первичного ключа id.
- 3 НФ: Атрибуты id, user_id, amount, currency_id, type, is_failed, created_at, updated_at, closed_at, period_start, period_end не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram

BUDGET {
    int id PK
    int user_id FK
    decimal amount
    int currency_id FK
    string type
    bool is_failed
    timestamptz created_at
    timestamptz updated_at
    timestampptz closed_at
    date period_start
    date period_end
}
```
## Таблица CATEGORY

Таблица `CATEGORY` содержит информацию о категориях транзакций.

<p> Функциональные зависимости: </p>

- `{id} -> {user_id, name, logo, created_at, updated_at}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты id, user_id, name, logo, created_at, updated_at являются атомарными.
- 2 НФ: Атрибуты id, user_id, name, logo, created_at, updated_at полностью функционально зависят от первичного ключа id.
- 3 НФ: Атрибуты id, user_id, name, logo, created_at, updated_at не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram

    CATEGORY {
        int id PK
        int user_id FK
        string name
        string logo
        timestamptz created_at
        timestamptz updated_at
    }
```
## Таблица OPERATION

Таблица `OPERATION` содержит информацию о транзакциях пользователей.

<p> Функциональные зависимости: </p>

- `{id} -> {account_id, category_id, type, name, sum, created_at, updated_at}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты id, account_id, category_id, type, name, sum, created_at, updated_at являются атомарными.
- 2 НФ: Атрибуты id, account_id, category_id, type, name, sum, created_at, updated_at полностью функционально зависят от первичного ключа id.
- 3 НФ: Атрибуты id, account_id, category_id, type, name, sum, created_at, updated_at не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram

  OPERATION{
        int id PK
        int account_id FK
        int category_id FK
        string type
        string name
        decimal sum
        timestampptz created_at
    }
```
## Таблица TRANSFER

Таблица `TRANSFER` содержит информацию о о переводах между пользовательскими счетами.

<p> Функциональные зависимости: </p>

- `{id} -> {from_account_id}`

<p> Нормальные формы: <p>

- 1 НФ: Атрибуты id, from_account_id являются атомарными.
- 2 НФ: Атрибуты id, from_account_id полностью функционально зависят от первичного ключа id.
- 3 НФ: Атрибуты id, from_account_id не зависят от других атрибутов.
- НФБК: 3 НФ + в таблице отсутствуют составные ключи.

```mermaid
erDiagram

    TRANSFER {
        int id PK
        int from_account_id FK
    }
```