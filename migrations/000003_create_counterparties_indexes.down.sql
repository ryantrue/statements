BEGIN;

-- Удаление уникального индекса на ИНН и КПП
DROP INDEX IF EXISTS idx_inn_kpp;

COMMIT;
