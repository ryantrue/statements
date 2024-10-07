BEGIN;

-- Создание ENUM типа для статуса ЕАИСТ
CREATE TYPE public.eaist_status_enum AS ENUM ('Активный', 'Завершен', 'Аннулирован');

-- Добавление нового ENUM поля для статуса ЕАИСТ в таблицу contracts
ALTER TABLE public.contracts
ALTER COLUMN eaist_status TYPE public.eaist_status_enum USING eaist_status::public.eaist_status_enum;

-- Добавление описания для ENUM типа
COMMENT ON TYPE public.eaist_status_enum IS 'Перечисление статусов в системе ЕАИСТ';

COMMIT;