BEGIN;

-- Удаление индексов для таблицы contracts
DROP INDEX IF EXISTS idx_contract_number;
DROP INDEX IF EXISTS idx_contract_date;

COMMIT;
