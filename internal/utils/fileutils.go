package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
)

// SaveFile сохраняет файл на диск в указанную директорию
func SaveFile(fileHeader *multipart.FileHeader, uploadDir string) (string, error) {
	// Проверка существования и создание директории, если она не существует
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("ошибка создания директории %s: %w", uploadDir, err)
	}

	// Открытие файла из запроса
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("ошибка открытия файла %s: %w", fileHeader.Filename, err)
	}
	defer file.Close()

	// Полный путь до файла
	filePath := filepath.Join(uploadDir, fileHeader.Filename)

	// Создание файла на диске
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("ошибка создания файла %s: %w", filePath, err)
	}
	defer out.Close()

	// Копирование содержимого файла
	if _, err := out.ReadFrom(file); err != nil {
		return "", fmt.Errorf("ошибка записи в файл %s: %w", filePath, err)
	}

	return filePath, nil
}
