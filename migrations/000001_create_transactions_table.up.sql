-- Создание таблицы transactions, если она не существует
CREATE TABLE IF NOT EXISTS transactions (
                                            id SERIAL PRIMARY KEY,                                      -- Первичный ключ
                                            account_number VARCHAR(20) NOT NULL,                        -- Номер счёта, обязателен
    bank VARCHAR(10) NOT NULL,                                  -- Название банка, обязательно
    date DATE NOT NULL CHECK (date <= CURRENT_DATE),            -- Дата транзакции, не может быть в будущем
    debit_account VARCHAR(20),                                  -- Счет дебета, может быть NULL
    inn VARCHAR(12) CHECK (inn IS NULL OR char_length(inn) IN (10, 12)), -- ИНН, должен быть 10 или 12 символов или NULL
    name TEXT,                                                  -- Имя контрагента, может быть NULL
    credit_account VARCHAR(20),                                 -- Счет кредита, может быть NULL
    inn_c VARCHAR(12) CHECK (inn_c IS NULL OR char_length(inn_c) IN (10, 12)), -- ИНН контрагента, может быть 10 или 12 символов или NULL
    name_c TEXT NOT NULL,                                       -- Имя контрагента кредита, обязательно
    debit NUMERIC(15, 2) DEFAULT 0.00 CHECK (debit >= 0),       -- Дебет, не может быть отрицательным
    credit NUMERIC(15, 2) DEFAULT 0.00 CHECK (credit >= 0),     -- Кредит, не может быть отрицательным
    document_number VARCHAR(20) NOT NULL,                       -- Номер документа, обязателен
    payment_description TEXT,                                   -- Описание платежа (опционально)
-- Предотвращение дублирования транзакций
    UNIQUE (account_number, date, document_number, debit_account, credit_account)
    );

-- Создание индекса для ускорения поиска по номеру счета и дате
CREATE INDEX IF NOT EXISTS idx_account_date ON transactions (account_number, date);

-- Создание индекса для ускорения поиска по дебету и кредиту
CREATE INDEX IF NOT EXISTS idx_debit_credit ON transactions (debit_account, credit_account);

COMMIT;