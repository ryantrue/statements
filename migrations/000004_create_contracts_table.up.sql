BEGIN;

-- Создание таблицы contracts с хранением ссылок на файлы
CREATE TABLE public.contracts (
                                  id SERIAL,                                        -- Уникальный идентификатор
                                  counterparty_id INT NOT NULL REFERENCES public.counterparties(id) ON DELETE RESTRICT, -- Привязка к контрагенту
                                  contract_number VARCHAR(50) NOT NULL,             -- Номер контракта, обязателен
                                  contract_date DATE NOT NULL CHECK (contract_date <= CURRENT_DATE), -- Дата заключения контракта
                                  execution_period DATE NOT NULL,                   -- Срок исполнения контракта
                                  amount NUMERIC(15, 2) NOT NULL CHECK (amount >= 0), -- Сумма договора
                                  eaist_registry_number VARCHAR(50),                -- Реестровый номер ЕАИСТ (опционально)
                                  payment_days INT CHECK (payment_days >= 0),       -- Дней на оплату
                                  validity_period DATE,                             -- Срок действия контракта (опционально)
                                  subject TEXT NOT NULL,                            -- Предмет контракта
                                  contract_type VARCHAR(50) NOT NULL,               -- Тип контракта
                                  work_type VARCHAR(100),                           -- Вид работ (опционально)
                                  conclusion_basis TEXT,                            -- Основание заключения
                                  procurement_type VARCHAR(50),                     -- Вид закупки (опционально)
                                  initiator TEXT,                                   -- Инициатор (опционально)
                                  eaist_status VARCHAR(50),                         -- Статус ЕАИСТ (опционально)
                                  eaist_link TEXT,                                  -- Ссылка на ЕАИСТ (опционально),

    -- Колонки для хранения путей к файлам
                                  contract_file_path TEXT,                          -- Путь к файлу контракта (PDF/DOCX)
                                  memo_file_path TEXT,                              -- Путь к служебной записке (PDF/DOCX)
                                  ecp_file_path TEXT,                               -- Путь к файлу ЭЦП (PDF/DOCX)
                                  technical_task_file_path TEXT,                    -- Путь к файлу технического задания (PDF/DOCX)
                                  additional_files_paths TEXT[],                    -- Пути к дополнительным файлам (массив)

                                  PRIMARY KEY (contract_number, contract_date)       -- Первичный ключ
) PARTITION BY RANGE (contract_date);                 -- Партиционирование по дате заключения контракта

-- Создание партиции для контрактов, заключенных в 2023 году
CREATE TABLE public.contracts_2023 PARTITION OF public.contracts
    FOR VALUES FROM ('2023-01-01') TO ('2024-01-01');

-- Создание партиции для контрактов, заключенных в 2024 году
CREATE TABLE public.contracts_2024 PARTITION OF public.contracts
    FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');

COMMIT;