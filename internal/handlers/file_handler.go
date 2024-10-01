package handlers

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"statements/internal/config"
	"statements/internal/python"
	"statements/internal/transactions"
	"statements/internal/utils"
	"sync"

	"github.com/gin-gonic/gin"
)

// HandleFileUploadGin обрабатывает загрузку файлов через Gin
func HandleFileUploadGin(c *gin.Context, cfg *config.Config) {
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка при обработке формы")
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.String(http.StatusBadRequest, "Не выбрано ни одного файла")
		return
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	successfulFiles := 0
	errorsOccurred := false

	// Буферизированный канал для асинхронной обработки каждого файла
	resultChan := make(chan error, len(files))

	for _, fileHeader := range files {
		wg.Add(1)

		go func(fileHeader *multipart.FileHeader) {
			defer wg.Done()

			// Открытие файла и его сохранение
			filePath, err := utils.SaveFile(fileHeader, cfg.FileUpload.UploadDir)
			if err != nil {
				log.Printf("Ошибка сохранения файла: %v", err)
				resultChan <- fmt.Errorf("Ошибка сохранения файла: %v", err)
				return
			}

			// Асинхронный вызов Python-скрипта для получения данных
			resultChanPython := python.GetDataFromPythonAsync([]string{filePath}, cfg)

			// Ожидание результата
			result := <-resultChanPython
			if result.Error != nil {
				log.Printf("Ошибка выполнения Python скрипта: %v", result.Error)
				resultChan <- fmt.Errorf("Ошибка обработки файла: %v", result.Error)
				return
			}

			// Логируем извлеченные данные для отладки
			for accountNumber, transactionsList := range result.Result.AccountTransactions {
				log.Printf("Номер счета: %s", accountNumber)
				log.Printf("Транзакции до очистки: %v", transactionsList)

				// Очищаем транзакции через функцию из пакета transactions
				cleanedTransactions := transactions.CleanTransactionList(transactionsList, result.Result.StatementType, accountNumber)

				// Логируем очищенные транзакции
				log.Printf("Очищенные транзакции для счета %s: %v", accountNumber, cleanedTransactions)

				// Проверка на наличие очищенных транзакций
				if len(cleanedTransactions) > 0 {
					// Сохраняем очищенные транзакции в базу данных через функцию из пакета transactions
					transactions.SaveTransactionsToDB(result.Result.StatementType, map[string][]map[string]interface{}{
						accountNumber: cleanedTransactions,
					})
				} else {
					log.Printf("Нет транзакций для сохранения в базу данных для счета %s", accountNumber)
				}
			}

			mu.Lock()
			successfulFiles++
			mu.Unlock()

			resultChan <- nil // Успешная обработка
		}(fileHeader)
	}

	// Ожидаем завершения всех горутин
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Чтение результатов обработки
	for err := range resultChan {
		if err != nil {
			errorsOccurred = true
			log.Printf("Ошибка: %v", err)
		}
	}

	// Возвращаем результат пользователю
	if errorsOccurred {
		c.String(http.StatusInternalServerError, "Произошли ошибки при обработке файлов")
	} else {
		c.String(http.StatusOK, fmt.Sprintf("Файлы успешно загружены и обработаны: %d", successfulFiles))
	}
}
