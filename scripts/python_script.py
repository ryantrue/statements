import sys
import pdfplumber
import json
import re
import logging
from abc import ABC, abstractmethod
from typing import List, Dict, Tuple, Optional

# Настройка логирования
def setup_logging():
    logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

class StatementProcessor(ABC):
    """Абстрактный класс для обработки выписок банков."""

    @abstractmethod
    def detect_accounts(self, text: str) -> List[str]:
        """Метод для извлечения всех номеров счетов из текста выписки."""
        pass

    @abstractmethod
    def process_transaction_row(self, row: List[str]) -> Dict[str, str]:
        """Метод для обработки одной строки транзакций."""
        pass

    @staticmethod
    def clean_text(text: Optional[str]) -> str:
        """Очищает текст от лишних символов и пробелов."""
        return text.strip().replace("\n", " ") if text else ""

class SberbankProcessor(StatementProcessor):
    """Класс для обработки выписок Сбербанка."""
    account_re = re.compile(r'ВЫПИСКА ОПЕРАЦИЙ ПО ЛИЦЕВОМУ СЧЕТУ\s(\d{20})')

    def detect_accounts(self, text: str) -> List[str]:
        """Извлекает все номера счетов из текста выписки."""
        return self.account_re.findall(text) if text else []

    def process_transaction_row(self, row: List[str]) -> Dict[str, str]:
        if len(row) < 9:
            return {}
        return {
            'date': self.clean_text(row[0]),
            'debit_account': self.clean_text(row[1]),
            'credit_account': self.clean_text(row[2]),
            'debit': self.clean_text(row[3]),
            'credit': self.clean_text(row[4]),
            'document_number': self.clean_text(row[5]),
            'vo_code': self.clean_text(row[6]),
            'bik': self.clean_text(row[7]),
            'payment_description': self.clean_text(row[8]) if len(row) > 8 else ''
        }

class VTBProcessor(StatementProcessor):
    """Класс для обработки выписок ВТБ."""
    account_re = re.compile(r'Счет\s(\d{20})\s\(Валюта\s\d{3},\sРоссийский\sрубль\)')

    def detect_accounts(self, text: str) -> List[str]:
        """Извлекает все номера счетов из текста выписки."""
        if not text:
            logging.warning("Текст страницы пустой, не удалось извлечь счета.")
            return []
        accounts = self.account_re.findall(text)
        logging.info(f"Найдены счета: {accounts}")
        return accounts

    def process_transaction_row(self, row: List[str]) -> Dict[str, str]:
        """Обрабатывает одну строку транзакций."""
        if len(row) < 9:
            return {}
        return {
            'date': self.clean_text(row[0]),
            'transaction_number': self.clean_text(row[1]),
            'operation_code': self.clean_text(row[2]),
            'inn': self.clean_text(row[3]),
            'bik': self.clean_text(row[4]),
            'account': self.clean_text(row[5]),
            'name': self.clean_text(row[6]),
            'debit': self.clean_text(row[7]),
            'credit': self.clean_text(row[8]),
            'description': self.clean_text(row[9]) if len(row) > 9 else ''
        }

class StatementFactory:
    """Фабрика для создания процессора выписки в зависимости от типа."""
    @staticmethod
    def get_processor(statement_type: str) -> Optional[StatementProcessor]:
        processors = {
            'СБЕР': SberbankProcessor,
            'ВТБ': VTBProcessor
        }
        return processors.get(statement_type)()

def detect_statement_type(text: str) -> str:
    """Определяет тип банковской выписки на основе текста первой страницы."""
    if not text:
        return 'unknown'
    if SberbankProcessor.account_re.search(text):
        return 'СБЕР'
    elif VTBProcessor.account_re.search(text):
        return 'ВТБ'
    return 'unknown'

def process_transaction_table(table: Optional[List[List[str]]], processor: StatementProcessor) -> List[Dict[str, str]]:
    """Обрабатывает таблицу транзакций с помощью переданного процессора."""
    transactions = []
    if not table:
        logging.warning("Таблица транзакций пустая или не найдена.")
        return transactions
    for row in table:
        if len(row) < 9 or row[0] == "Дата" or row[1] == "Счет":
            continue
        transaction = processor.process_transaction_row(row)
        if any(transaction.values()):
            transactions.append(transaction)
    return transactions

def extract_data_from_pdf(pdf: pdfplumber.PDF, processor: StatementProcessor) -> Dict[str, List[Dict[str, str]]]:
    """Извлекает данные транзакций для всех счетов."""
    account_transactions = {}
    current_account = None

    for page in pdf.pages:
        text = page.extract_text()
        if not text:
            logging.warning("Не удалось извлечь текст с одной из страниц.")
            continue

        # Извлекаем номера счетов
        new_accounts = processor.detect_accounts(text)
        if new_accounts:
            current_account = new_accounts[0]  # Обновляем текущий счет
            if current_account not in account_transactions:
                account_transactions[current_account] = []

        tables = page.extract_tables()
        if not tables:
            logging.warning("Таблицы на странице не найдены.")
            continue

        for table in tables:
            transactions = process_transaction_table(table, processor)
            if current_account:
                account_transactions[current_account].extend(transactions)

            # Пропускаем строки с ИТОГО или Количество операций
            if any(keyword in row[0] for row in table if row[0] for keyword in ['ИТОГО', 'Количество операций']):
                logging.info(f"Завершение транзакций для счета {current_account}")
                current_account = None  # Завершаем текущий счет

    logging.info(f"Извлеченные транзакции: {account_transactions}")
    return account_transactions

def extract_transaction_data(pdf_path: str) -> Tuple[Dict[str, List[Dict[str, str]]], str, str]:
    """Извлекает номера счетов и транзакции из PDF файла."""
    try:
        with pdfplumber.open(pdf_path) as pdf:
            first_page_text = pdf.pages[0].extract_text()
            statement_type = detect_statement_type(first_page_text)
            processor = StatementFactory.get_processor(statement_type)

            if processor:
                account_transactions = extract_data_from_pdf(pdf, processor)
                return account_transactions, first_page_text, statement_type
            else:
                logging.warning("Не удалось определить тип выписки")
                return {}, first_page_text, 'unknown'
    except Exception as e:
        logging.error(f"Неожиданная ошибка: {e}")
        raise

def main(pdf_path: str):
    """Основная функция программы. Извлекает и выводит данные из PDF файла."""
    try:
        account_transactions, first_page_text, statement_type = extract_transaction_data(pdf_path)
        if account_transactions:
            result = {
                'account_transactions': account_transactions,
                'first_page_text': first_page_text,
                'statement_type': statement_type
            }
            sys.stdout.buffer.write(json.dumps(result, ensure_ascii=False, indent=4).encode('utf-8'))
        else:
            logging.error("Не удалось определить номера счетов или транзакции.")
            print("Ошибка: не удалось извлечь данные.", file=sys.stderr)
    except Exception as e:
        logging.error(f"Ошибка в главной функции: {e}")
        print(f"Ошибка при обработке PDF: {e}", file=sys.stderr)

if __name__ == "__main__":
    setup_logging()
    if len(sys.argv) < 2:
        logging.error("Пожалуйста, укажите путь к файлу PDF.")
        print("Пожалуйста, укажите путь к файлу PDF.", file=sys.stderr)
    else:
        pdf_path = sys.argv[1]
        main(pdf_path)