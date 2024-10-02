package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"net/http"
	"statements/internal/database"
)

// HandleDownloadExcel формирует Excel-файл с данными из таблицы transactions и отправляет его пользователю
func HandleDownloadExcel(c *gin.Context) {
	// Создаем новый Excel-файл
	f := excelize.NewFile()

	// Добавляем заголовки в первую строку
	f.SetCellValue("Sheet1", "A1", "Account Number")
	f.SetCellValue("Sheet1", "B1", "Bank")
	f.SetCellValue("Sheet1", "C1", "Date")
	f.SetCellValue("Sheet1", "D1", "Debit Account")
	f.SetCellValue("Sheet1", "E1", "Credit Account")
	f.SetCellValue("Sheet1", "F1", "Debit")
	f.SetCellValue("Sheet1", "G1", "Credit")
	f.SetCellValue("Sheet1", "H1", "INN")
	f.SetCellValue("Sheet1", "I1", "Name")
	f.SetCellValue("Sheet1", "J1", "INN C")
	f.SetCellValue("Sheet1", "K1", "Name C")
	f.SetCellValue("Sheet1", "L1", "Document Number")
	f.SetCellValue("Sheet1", "M1", "Payment Description")

	// Извлекаем данные из базы данных
	rows, err := database.DB.Query(`SELECT account_number, bank, date, debit_account, credit_account, debit, credit, inn, name, inn_c, name_c, document_number, payment_description FROM transactions`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка чтения данных из базы"})
		return
	}
	defer rows.Close()

	// Перебираем строки и записываем данные в файл
	rowIndex := 2
	for rows.Next() {
		var accountNumber, bank, date, debitAccount, creditAccount, debit, credit, inn, name, innC, nameC, documentNumber, paymentDescription string

		err = rows.Scan(&accountNumber, &bank, &date, &debitAccount, &creditAccount, &debit, &credit, &inn, &name, &innC, &nameC, &documentNumber, &paymentDescription)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки данных"})
			return
		}

		// Заполняем ячейки файла
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", rowIndex), accountNumber)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", rowIndex), bank)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", rowIndex), date)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", rowIndex), debitAccount)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", rowIndex), creditAccount)
		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", rowIndex), debit)
		f.SetCellValue("Sheet1", fmt.Sprintf("G%d", rowIndex), credit)
		f.SetCellValue("Sheet1", fmt.Sprintf("H%d", rowIndex), inn)
		f.SetCellValue("Sheet1", fmt.Sprintf("I%d", rowIndex), name)
		f.SetCellValue("Sheet1", fmt.Sprintf("J%d", rowIndex), innC)
		f.SetCellValue("Sheet1", fmt.Sprintf("K%d", rowIndex), nameC)
		f.SetCellValue("Sheet1", fmt.Sprintf("L%d", rowIndex), documentNumber)
		f.SetCellValue("Sheet1", fmt.Sprintf("M%d", rowIndex), paymentDescription)
		rowIndex++
	}

	// Устанавливаем заголовки для скачивания файла
	c.Header("Content-Disposition", "attachment; filename=transactions.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	// Отправляем файл пользователю
	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации файла Excel"})
	}
}
