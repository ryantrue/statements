BEGIN;

-- Создание схемы public, если она еще не существует
CREATE SCHEMA IF NOT EXISTS public;

-- Создание таблицы counterparties
CREATE TABLE IF NOT EXISTS public.counterparties (
                                                     id SERIAL PRIMARY KEY,                                     -- Первичный ключ
                                                     name TEXT NOT NULL,                                        -- Наименование контрагента, обязательно
                                                     inn VARCHAR(12) NOT NULL CHECK (char_length(inn) IN (10, 12)), -- ИНН, 10 или 12 символов
    kpp VARCHAR(9) CHECK (char_length(kpp) = 9),               -- КПП, 9 символов (опционально)
    UNIQUE (inn, kpp)                                          -- Уникальность сочетания ИНН и КПП
    );

-- Добавление описаний для таблицы counterparties
COMMENT ON TABLE public.counterparties IS 'Таблица для хранения информации о контрагентах';
COMMENT ON COLUMN public.counterparties.name IS 'Название контрагента';
COMMENT ON COLUMN public.counterparties.inn IS 'ИНН контрагента, 10 или 12 символов';
COMMENT ON COLUMN public.counterparties.kpp IS 'КПП контрагента, 9 символов (опционально)';

COMMIT;