package handlers

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"statements/internal/config"
)

// HandleContractSubmission обрабатывает форму добавления контракта и загрузку файлов
func HandleContractSubmission(c *gin.Context, cfg *config.Config, db *sql.DB) {
	// Парсинг формы
	counterpartyID := c.PostForm("counterparty_id")
	contractNumber := c.PostForm("contract_number")
	contractDate := c.PostForm("contract_date")
	executionPeriod := c.PostForm("execution_period")
	amount := c.PostForm("amount")
	contractType := c.PostForm("contract_type")
	subject := c.PostForm("subject")

	fmt.Printf("Контракт: %s, Дата: %s, Срок исполнения: %s, Сумма: %s, Тип: %s, Предмет: %s, Контрагент: %s\n",
		contractNumber, contractDate, executionPeriod, amount, contractType, subject, counterpartyID)

	// Функция для сохранения файлов
	saveFile := func(fileHeader *multipart.FileHeader, directory string) (string, error) {
		file, err := fileHeader.Open()
		if err != nil {
			return "", err
		}
		defer file.Close()

		// Путь для сохранения файла
		filePath := filepath.Join(directory, fileHeader.Filename)
		outFile, err := os.Create(filePath)
		if err != nil {
			return "", err
		}
		defer outFile.Close()

		// Запись файла
		_, err = outFile.ReadFrom(file)
		return filePath, err
	}

	// Создание директории для файлов контракта
	baseDir := filepath.Join(cfg.FileUpload.UploadDir, contractNumber)
	err := os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		log.Printf("Ошибка создания директории: %v", err)
		c.String(http.StatusInternalServerError, "Ошибка создания директории для файлов")
		return
	}

	// Сохранение путей к файлам
	var contractFilePath, memoFilePath, ecpFilePath, technicalTaskFilePath string
	var additionalFilePaths []string

	// Загрузка обязательных файлов
	fileHeader, err := c.FormFile("contract_file")
	if err == nil {
		contractFilePath, err = saveFile(fileHeader, baseDir)
		if err != nil {
			log.Printf("Ошибка сохранения файла контракта: %v", err)
			c.String(http.StatusInternalServerError, "Ошибка сохранения файла контракта")
			return
		}
	}

	// Загрузка необязательных файлов
	files := map[string]*string{
		"memo_file":           &memoFilePath,
		"ecp_file":            &ecpFilePath,
		"technical_task_file": &technicalTaskFilePath,
	}
	for key, path := range files {
		fileHeader, err := c.FormFile(key)
		if err == nil {
			*path, err = saveFile(fileHeader, baseDir)
			if err != nil {
				log.Printf("Ошибка сохранения файла %s: %v", key, err)
				c.String(http.StatusInternalServerError, fmt.Sprintf("Ошибка сохранения файла %s", key))
				return
			}
		}
	}

	// Обработка дополнительных файлов
	form, _ := c.MultipartForm()
	additionalFiles := form.File["additional_files[]"]
	for _, fileHeader := range additionalFiles {
		filePath, err := saveFile(fileHeader, baseDir)
		if err != nil {
			log.Printf("Ошибка при сохранении дополнительного файла: %v", err)
			c.String(http.StatusInternalServerError, "Ошибка сохранения дополнительных файлов")
			return
		}
		additionalFilePaths = append(additionalFilePaths, filePath)
	}

	// Сохранение данных в базу данных
	query := `INSERT INTO contracts 
                (counterparty_id, contract_number, contract_date, execution_period, amount, contract_type, subject, 
                contract_file_path, memo_file_path, ecp_file_path, technical_task_file_path, additional_files_paths) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err = db.Exec(query, counterpartyID, contractNumber, contractDate, executionPeriod, amount, contractType, subject,
		contractFilePath, memoFilePath, ecpFilePath, technicalTaskFilePath, additionalFilePaths)
	if err != nil {
		log.Printf("Ошибка сохранения данных контракта в базу: %v", err)
		c.String(http.StatusInternalServerError, "Ошибка сохранения данных контракта")
		return
	}

	// Возвращаем результат пользователю
	c.String(http.StatusOK, "Контракт успешно добавлен!")
}

// HandleCounterpartiesList возвращает список всех контрагентов
func HandleCounterpartiesList(c *gin.Context, db *sql.DB) {
	query := "SELECT id, name FROM counterparties"
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Ошибка получения списка контрагентов: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных контрагентов"})
		return
	}
	defer rows.Close()

	var counterparties []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Printf("Ошибка чтения строки: %v", err)
			continue
		}
		counterparties = append(counterparties, map[string]interface{}{
			"id":   id,
			"name": name,
		})
	}

	c.JSON(http.StatusOK, counterparties)
}
