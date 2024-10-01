package transactions

import (
	"fmt"
	"regexp"
	"strings"
)

// CleanTransaction очищает транзакцию от лишних символов и форматирует их
func CleanTransaction(transaction map[string]interface{}, bank, accountNumber string) map[string]interface{} {
	cleanedTransaction := make(map[string]interface{})

	// Используем для замены пробелов
	spaceReplacer := regexp.MustCompile(`\s+`)

	for key, value := range transaction {
		strValue := fmt.Sprintf("%v", value)
		strValue = strings.TrimSpace(strValue)
		strValue = spaceReplacer.ReplaceAllString(strValue, " ")

		// Форматируем поля debit и credit
		if key == "debit" || key == "credit" {
			strValue = cleanNumber(strValue) // Очищаем и форматируем как числовое значение

			// Для банка СБЕР, если поле пустое, то устанавливаем его в 0.00
			if bank == "СБЕР" && (strValue == "<nil>" || strValue == "") {
				strValue = "0.00"
			}
		}

		if strValue == "" || strValue == "Кредит" || strValue == "Дебет" {
			cleanedTransaction[key] = nil
		} else {
			cleanedTransaction[key] = strValue
		}
	}

	return cleanedTransaction
}

// CleanTransactionList очищает список транзакций с проверкой наличия индикаторов завершения транзакций для счетов
func CleanTransactionList(transactions []map[string]interface{}, bank, accountNumber string) []map[string]interface{} {
	cleanedTransactions := make([]map[string]interface{}, 0)
	var stopProcessing bool

	for _, transaction := range transactions {
		if IsHeaderRow(transaction) {
			continue
		}

		// Проверяем индикаторы завершения для разных банков
		if bank == "СБЕР" && containsStopPhrase(transaction, "Количество операций") {
			stopProcessing = true
		} else if bank == "ВТБ" && containsStopPhrase(transaction, "ИТОГО за период с") {
			stopProcessing = true
		}

		// Если нашли фразу для завершения обработки, выходим из цикла
		if stopProcessing {
			break
		}

		// Очищаем транзакцию
		cleanedTransaction := CleanTransaction(transaction, bank, accountNumber)

		// Проверяем наличие валидных значений credit или debit для Сбербанка
		if bank == "СБЕР" && !HasValidCreditOrDebit(cleanedTransaction) {
			continue
		}

		// Проверка, есть ли транзакции для текущего счета
		isEmpty := true
		for _, value := range cleanedTransaction {
			if value != nil {
				isEmpty = false
				break
			}
		}

		if !isEmpty {
			// Обрабатываем транзакции, если есть валидные данные
			cleanedTransactions = append(cleanedTransactions, cleanedTransaction)
		}
	}

	return cleanedTransactions
}

// containsStopPhrase проверяет, содержит ли транзакция ключевые фразы для остановки обработки
func containsStopPhrase(transaction map[string]interface{}, phrase string) bool {
	for _, value := range transaction {
		if strings.Contains(fmt.Sprintf("%v", value), phrase) {
			return true
		}
	}
	return false
}

// IsHeaderRow проверяет, является ли строка заголовком
func IsHeaderRow(transaction map[string]interface{}) bool {
	if transaction["account"] == "Счет" || transaction["bik"] == "БИК банка" || transaction["credit"] == "Кредит" || transaction["debit"] == "Дебет" {
		return true
	}
	return false
}

// HasValidCreditOrDebit проверяет, есть ли значения в полях credit или debit
func HasValidCreditOrDebit(transaction map[string]interface{}) bool {
	credit, creditExists := transaction["credit"].(string)
	debit, debitExists := transaction["debit"].(string)

	if creditExists && credit != "0.00" || debitExists && debit != "0.00" {
		return true
	}
	return false
}

// cleanNumber форматирует строку в правильный числовой формат для базы данных
func cleanNumber(number string) string {
	// Удаляем все виды пробелов, включая неразрывные пробелы (U+00A0)
	number = strings.ReplaceAll(number, " ", "")  // Удаляем неразрывные пробелы (U+00A0)
	number = strings.ReplaceAll(number, " ", "")  // Удаляем обычные пробелы
	number = strings.ReplaceAll(number, ",", ".") // Заменяем запятые на точки
	return number
}
