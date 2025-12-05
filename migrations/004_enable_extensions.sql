-- ========================================================
-- Включение расширений для мониторинга и анализа
-- ========================================================

CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

SELECT * FROM pg_extension WHERE extname = 'pg_stat_statements';
