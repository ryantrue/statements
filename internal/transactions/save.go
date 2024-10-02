package transactions

import (
	"context"
	"fmt"
	"log"
	"statements/internal/database"
	"strings"
	"time"
)

// isValidAccount проверяет корректность номера счета (ожидаемая длина — 20 символов)
func isValidAccount(account string) bool {
	return len(account) == 20
}

// isValidInn проверяет корректность ИНН (ожидаемая длина — 10 или 12 символов)
func isValidInn(inn string) bool {
	return len(inn) == 10 || len(inn) == 12
}

// splitAccountInfo разбивает информацию о счете на компоненты (счет, ИНН, имя)
func splitAccountInfo(accountInfo string) (account, inn, name string) {
	parts := strings.Fields(accountInfo)

	if len(parts) == 0 {
		return "", "", ""
	}

	// Определяем, является ли первая часть счетом или ИНН
	if isValidAccount(parts[0]) {
		account = parts[0]
		if len(parts) > 1 {
			inn = parts[1]
		}
		if len(parts) > 2 {
			name = strings.Join(parts[2:], " ")
		}
	} else if isValidInn(parts[0]) {
		inn = parts[0]
		if len(parts) > 1 {
			name = strings.Join(parts[1:], " ")
		}
	}

	return
}

// SaveTransactionsToDB сохраняет очищенные транзакции для всех счетов в базе данных PostgreSQL
func SaveTransactionsToDB(bank string, accountTransactions map[string][]map[string]interface{}) {
	for accountNumber, transactions := range accountTransactions {
		if len(transactions) == 0 {
			log.Printf("Нет транзакций для сохранения в базу данных для счета %s", accountNumber)
			continue
		}

		log.Printf("Начало записи транзакций для счета %s и банка %s", accountNumber, bank)
		for _, transaction := range transactions {
			switch bank {
			case "СБЕР":
				saveTransaction(accountNumber, bank, transaction, processSberTransaction)
			case "ВТБ":
				saveTransaction(accountNumber, bank, transaction, processVTBTransaction)
			default:
				log.Printf("Неизвестный банк: %s", bank)
			}
		}
	}
}

// saveTransaction — шаблонная функция для сохранения транзакций
func saveTransaction(accountNumber, bank string, transaction map[string]interface{}, process func(string, map[string]interface{}) (string, string, string, string, string, string)) {
	debitAccount, inn, name, creditAccount, innC, nameC := process(accountNumber, transaction)

	documentNumber := extractDocumentNumber(transaction)
	paymentDescription := extractPaymentDescription(transaction)

	exists, err := existsTransaction(
		accountNumber,
		getStringValue(transaction, "date"),
		getStringValue(transaction, "debit"),
		getStringValue(transaction, "credit"),
		documentNumber,
		paymentDescription,
		debitAccount,
		creditAccount,
		inn,
		name,
		innC,
		nameC,
	)
	if err != nil {
		log.Printf("Ошибка проверки дубликата транзакции для счета %s: %v", accountNumber, err)
		return
	}

	if exists {
		log.Printf("Транзакция для счета %s уже существует, пропускаем.", accountNumber)
		return
	}

	err = insertTransaction(accountNumber, bank, transaction, debitAccount, creditAccount, inn, name, innC, nameC, documentNumber, paymentDescription)
	if err != nil {
		log.Printf("Ошибка вставки транзакции для счета %s: %v", accountNumber, err)
	}
}

// processSberTransaction обрабатывает транзакции для Сбербанка
func processSberTransaction(accountNumber string, transaction map[string]interface{}) (debitAccount, inn, name, creditAccount, innC, nameC string) {
	debitAccount, inn, name = splitAccountInfo(getStringValue(transaction, "debit_account"))
	creditAccount, innC, nameC = splitAccountInfo(getStringValue(transaction, "credit_account"))
	return
}

// processVTBTransaction обрабатывает транзакции для ВТБ
func processVTBTransaction(accountNumber string, transaction map[string]interface{}) (debitAccount, inn, name, creditAccount, innC, nameC string) {
	// Если сумма дебета равна 0, значит это приход на счет
	if getStringValue(transaction, "debit") == "0.00" {
		creditAccount = accountNumber
		innC = getStringValue(transaction, "inn")
		nameC = getStringValue(transaction, "name")

		debitAccount = getStringValue(transaction, "account")
		inn = "7719034354"
		name = `КАЗЕННОЕ ПРЕДПРИЯТИЕ "МОСКОВСКАЯ ЭНЕРГЕТИЧЕСКАЯ ДИРЕКЦИЯ"`
	} else {
		debitAccount = accountNumber
		inn = getStringValue(transaction, "inn")
		name = getStringValue(transaction, "name")

		creditAccount = getStringValue(transaction, "account")
		innC = "7719034354"
		nameC = `КАЗЕННОЕ ПРЕДПРИЯТИЕ "МОСКОВСКАЯ ЭНЕРГЕТИЧЕСКАЯ ДИРЕКЦИЯ"`
	}
	return
}

// convertDateToISO преобразует дату из формата DD.MM.YYYY в формат YYYY-MM-DD
func convertDateToISO(date string) (string, error) {
	// Если дата уже в формате YYYY-MM-DD, просто возвращаем её
	if len(date) == 10 && date[4] == '-' && date[7] == '-' {
		return date, nil
	}

	// Преобразуем из формата DD.MM.YYYY в формат YYYY-MM-DD
	parsedDate, err := time.Parse("02.01.2006", date) // Формат для DD.MM.YYYY
	if err != nil {
		return "", fmt.Errorf("ошибка преобразования даты: %w", err)
	}
	return parsedDate.Format("2006-01-02"), nil // Возвращаем дату в формате YYYY-MM-DD
}

// insertTransaction вставляет транзакцию в базу данных
func insertTransaction(accountNumber, bank string, transaction map[string]interface{}, debitAccount, creditAccount, inn, name, innC, nameC, documentNumber, paymentDescription string) error {
	log.Printf("Вставляем транзакцию для счета %s, банк %s, дата %s", accountNumber, bank, getStringValue(transaction, "date"))

	// Преобразование даты в формат YYYY-MM-DD
	isoDate, err := convertDateToISO(getStringValue(transaction, "date"))
	if err != nil {
		log.Printf("Ошибка преобразования даты для счета %s: %v", accountNumber, err)
		return err
	}

	log.Printf("Дата после преобразования: %s", isoDate)

	_, err = database.DB.ExecContext(context.Background(),
		`INSERT INTO transactions (account_number, bank, date, debit_account, credit_account, debit, credit, inn, name, inn_c, name_c, document_number, payment_description)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		accountNumber,
		bank,
		isoDate, // Используем преобразованную дату
		debitAccount,
		creditAccount,
		getStringValue(transaction, "debit"),
		getStringValue(transaction, "credit"),
		inn,
		name,
		innC,
		nameC,
		documentNumber,
		paymentDescription)
	if err != nil {
		log.Printf("Ошибка вставки транзакции для счета %s: %v", accountNumber, err)
		return err
	}

	log.Printf("Транзакция для счета %s успешно вставлена в базу данных.", accountNumber)
	return nil
}

// extractDocumentNumber извлекает номер документа
func extractDocumentNumber(transaction map[string]interface{}) string {
	if docNum, ok := transaction["document_number"].(string); ok {
		return docNum
	}
	return getStringValue(transaction, "transaction_number")
}

// extractPaymentDescription извлекает описание платежа
func extractPaymentDescription(transaction map[string]interface{}) string {
	if desc, ok := transaction["payment_description"].(string); ok {
		return desc
	}
	return getStringValue(transaction, "description")
}

// getStringValue безопасно извлекает строковое значение из карты
func getStringValue(transaction map[string]interface{}, key string) string {
	if val, ok := transaction[key]; ok && val != nil {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return ""
}
