-- Создание базы данных, если она не существует
DO
$do$
BEGIN
   IF NOT EXISTS (
      SELECT 1 FROM pg_database WHERE datname = 'postgres'
   ) THEN
      EXECUTE 'CREATE DATABASE postgres';
END IF;
END
$do$;

-- Подключитесь к новой базе данных вручную через клиент psql или в коде
-- \c postgres

-- Создание таблицы transactions, если она не существует
CREATE TABLE IF NOT EXISTS transactions (
                                            id SERIAL PRIMARY KEY,                        -- Первичный ключ
                                            account_number VARCHAR(20) NOT NULL,          -- Номер счёта, обязателен
    bank VARCHAR(10) NOT NULL,                    -- Название банка, обязательно
    date DATE NOT NULL CHECK (date <= CURRENT_DATE),  -- Дата транзакции, обязательно и не может быть в будущем
    debit_account VARCHAR(20) NOT NULL,           -- Счет дебета, обязателен
    inn VARCHAR(12) CHECK (char_length(inn) = 10 OR char_length(inn) = 12),  -- ИНН, должен быть 10 или 12 символов
    name TEXT NOT NULL,                           -- Имя контрагента, обязательно
    credit_account VARCHAR(20) NOT NULL,          -- Счет кредита, обязателен
    inn_c VARCHAR(12) CHECK (char_length(inn_c) = 10 OR char_length(inn_c) = 12),  -- ИНН контрагента, 10 или 12 символов
    name_c TEXT NOT NULL,                         -- Имя контрагента кредита, обязательно
    debit NUMERIC(15, 2) DEFAULT 0.00 CHECK (debit >= 0),  -- Дебет, по умолчанию 0.00, не может быть отрицательным
    credit NUMERIC(15, 2) DEFAULT 0.00 CHECK (credit >= 0), -- Кредит, по умолчанию 0.00, не может быть отрицательным
    document_number VARCHAR(20) NOT NULL,         -- Номер документа, обязателен
    payment_description TEXT,                     -- Описание платежа (опционально)

-- Ограничения для предотвращения дублирования транзакций
    UNIQUE (account_number, date, document_number, debit_account, credit_account)
    );

-- Создание индекса для ускорения поиска транзакций по номеру счета и дате
CREATE INDEX IF NOT EXISTS idx_account_date ON transactions (account_number, date);

-- Добавление дополнительного индекса для поиска по дебетовому и кредитовому счетам
CREATE INDEX IF NOT EXISTS idx_debit_credit ON transactions (debit_account, credit_account);