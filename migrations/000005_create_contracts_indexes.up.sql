BEGIN;

-- Индекс для ускорения поиска по номеру контракта
CREATE INDEX IF NOT EXISTS idx_contract_number ON public.contracts (contract_number);

-- Индекс для ускорения поиска по дате заключения контракта
CREATE INDEX IF NOT EXISTS idx_contract_date ON public.contracts (contract_date);

COMMIT;