BEGIN;

-- Создание уникального индекса для ИНН и КПП
CREATE UNIQUE INDEX IF NOT EXISTS idx_inn_kpp ON public.counterparties (inn, kpp);

COMMIT;