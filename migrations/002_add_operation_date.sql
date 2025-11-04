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
