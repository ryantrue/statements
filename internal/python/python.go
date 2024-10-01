package python

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"statements/internal/config"
	"statements/internal/models"
	"strings"
	"sync"
)

// AsyncResult представляет структуру для асинхронного результата и ошибок
type AsyncResult struct {
	FilePath string
	Result   models.Result
	Error    error
}

// GetDataFromPythonAsync вызывает Python-скрипт асинхронно для каждого файла
func GetDataFromPythonAsync(pdfPaths []string, cfg *config.Config) <-chan AsyncResult {
	if len(pdfPaths) == 0 {
		log.Fatal("Список файлов PDF пуст. Необходимо указать хотя бы один файл.")
	}

	resultChan := make(chan AsyncResult, len(pdfPaths)) // Буферизированный канал для асинхронной обработки
	var wg sync.WaitGroup

	for _, pdfPath := range pdfPaths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()

			// Выполняем Python-скрипт и обрабатываем результат
			result, err := runPythonScript(path, cfg)
			if err != nil {
				log.Printf("Ошибка при вызове Python скрипта для файла %s: %v", path, err)
				resultChan <- AsyncResult{FilePath: path, Error: err}
				return
			}

			log.Printf("Успешный вызов Python скрипта для файла %s", path)
			resultChan <- AsyncResult{FilePath: path, Result: result}
		}(pdfPath)
	}

	// Закрываем канал после завершения всех горутин
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return resultChan
}

// runPythonScript запускает Python-скрипт и возвращает результат
func runPythonScript(pdfPath string, cfg *config.Config) (models.Result, error) {
	if pdfPath == "" {
		return models.Result{}, fmt.Errorf("путь к PDF файлу не может быть пустым")
	}

	log.Printf("Запуск Python скрипта для файла: %s", pdfPath)
	cmd := exec.Command(cfg.Python.Interpreter, cfg.Python.ScriptPath, pdfPath)

	// Буферы для стандартного вывода и вывода ошибок
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Запускаем команду и проверяем на ошибки выполнения
	if err := cmd.Run(); err != nil {
		return models.Result{}, fmt.Errorf("ошибка выполнения Python скрипта для файла %s: %w. stderr: %s", pdfPath, err, stderr.String())
	}

	// Логируем предупреждения из stderr
	if stderr.Len() > 0 {
		log.Printf("Предупреждение из Python скрипта для файла %s:\n%s", pdfPath, stderr.String())
	}

	// Получаем и очищаем вывод команды
	output := strings.TrimSpace(stdout.String())
	log.Printf("Вывод Python скрипта для файла %s:\n%s", pdfPath, output)

	// Проверяем, что вывод является корректным JSON
	if !strings.HasPrefix(output, "{") {
		return models.Result{}, fmt.Errorf("вывод Python скрипта для файла %s не является корректным JSON: %s", pdfPath, output)
	}

	// Парсим JSON
	var result models.Result
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return models.Result{}, fmt.Errorf("ошибка парсинга JSON для файла %s: %w. Вывод: %s", pdfPath, err, output)
	}

	log.Printf("Успешный парсинг JSON для файла %s", pdfPath)
	return result, nil
}
