-- Добавление поля operation_date в таблицу operation
-- Это поле будет содержать дату операции, которая может отличаться от created_at

ALTER TABLE operation 
ADD COLUMN operation_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP;

-- Обновляем существующие записи, устанавливая operation_date = created_at
UPDATE operation 
SET operation_date = created_at 
WHERE operation_date IS NULL;

-- Делаем поле обязательным
ALTER TABLE operation 
ALTER COLUMN operation_date SET NOT NULL;


insert into currency (code, currency_name, logo_hashed_id, created_at) values ('usd', 'usd', 'ahds', NOW());


ALTER TABLE operation
ADD COLUMN operation_name_tsv tsvector
    GENERATED ALWAYS AS (to_tsvector('russian', operation_name)) STORED;


CREATE INDEX operation_name_tsv_idx
ON operation USING GIN (operation_name_tsv);


CREATE INDEX operation_filters_idx
    ON operation (category_id, operation_type, created_at);
