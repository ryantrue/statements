package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"net/http"
	"statements/internal/database"
)

// DataExporter интерфейс, описывающий структуру данных для выгрузки
type DataExporter interface {
	GetHeaders() []string
	GetRows() ([]map[string]interface{}, error)
}

// ExcelFileExporter создает Excel-файл из данных DataExporter
func ExcelFileExporter(c *gin.Context, exporter DataExporter, filename string) {
	f, err := createExcelFile(exporter)
	if err != nil {
		logrus.Errorf("Error creating Excel file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации файла Excel"})
		return
	}

	// Устанавливаем заголовки для скачивания файла
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xlsx", filename))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	// Отправляем файл пользователю
	if err := f.Write(c.Writer); err != nil {
		logrus.Errorf("Error writing Excel file to response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка отправки файла Excel"})
	}
}

// createExcelFile создает Excel-файл из данных
func createExcelFile(exporter DataExporter) (*excelize.File, error) {
	f := excelize.NewFile()
	headers := exporter.GetHeaders()

	// Добавляем заголовки
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue("Sheet1", cell, header)
	}

	// Получаем строки данных
	rows, err := exporter.GetRows()
	if err != nil {
		return nil, err
	}

	// Заполняем строки данными
	for i, row := range rows {
		rowIndex := i + 2
		for j, header := range headers {
			cell := fmt.Sprintf("%s%d", string(rune('A'+j)), rowIndex)
			f.SetCellValue("Sheet1", cell, row[header])
		}
	}

	return f, nil
}

// TransactionsExporter экспорт данных для таблицы transactions
type TransactionsExporter struct{}

// GetHeaders возвращает заголовки для таблицы transactions
func (e *TransactionsExporter) GetHeaders() []string {
	return []string{
		"Account Number", "Bank", "Date", "Debit Account", "Credit Account",
		"Debit", "Credit", "INN", "Name", "INN C", "Name C", "Document Number", "Payment Description",
	}
}

// GetRows возвращает строки данных из таблицы transactions
func (e *TransactionsExporter) GetRows() ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(`SELECT account_number, bank, date, debit_account, credit_account, debit, credit, inn, name, inn_c, name_c, document_number, payment_description FROM transactions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var accountNumber, bank, date, debitAccount, creditAccount, debit, credit, inn, name, innC, nameC, documentNumber, paymentDescription string
		err = rows.Scan(&accountNumber, &bank, &date, &debitAccount, &creditAccount, &debit, &credit, &inn, &name, &innC, &nameC, &documentNumber, &paymentDescription)
		if err != nil {
			return nil, err
		}
		results = append(results, map[string]interface{}{
			"Account Number":      accountNumber,
			"Bank":                bank,
			"Date":                date,
			"Debit Account":       debitAccount,
			"Credit Account":      creditAccount,
			"Debit":               debit,
			"Credit":              credit,
			"INN":                 inn,
			"Name":                name,
			"INN C":               innC,
			"Name C":              nameC,
			"Document Number":     documentNumber,
			"Payment Description": paymentDescription,
		})
	}

	return results, nil
}

// HandleDownloadTransactionsExcel обработчик для скачивания Excel с транзакциями
func HandleDownloadTransactionsExcel(c *gin.Context) {
	exporter := &TransactionsExporter{}
	ExcelFileExporter(c, exporter, "transactions")
}

// Здесь можно добавить новые экспортеры для других таблиц
// Например, для выгрузки другой таблицы можно реализовать аналогичный экспорт
