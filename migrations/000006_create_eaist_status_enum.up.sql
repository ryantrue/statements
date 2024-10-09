-- Создание ENUM типа для статуса ЕАИСТ, если он ещё не существует
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'eaist_status_enum') THEN
CREATE TYPE public.eaist_status_enum AS ENUM ('Активный', 'Завершен', 'Аннулирован');
END IF;
END $$;

-- Добавление нового ENUM поля для статуса ЕАИСТ в таблицу contracts
ALTER TABLE public.contracts
ALTER COLUMN eaist_status TYPE public.eaist_status_enum USING eaist_status::public.eaist_status_enum;

-- Добавление описания для ENUM типа
COMMENT ON TYPE public.eaist_status_enum IS 'Перечисление статусов в системе ЕАИСТ';