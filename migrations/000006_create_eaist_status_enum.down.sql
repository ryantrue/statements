BEGIN;

-- Приведение столбца eaist_status в contracts обратно к типу VARCHAR
ALTER TABLE public.contracts
ALTER COLUMN eaist_status TYPE VARCHAR(50);

-- Удаление типа ENUM для статуса ЕАИСТ
DROP TYPE IF EXISTS public.eaist_status_enum;

COMMIT;