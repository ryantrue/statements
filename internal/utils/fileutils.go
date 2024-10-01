package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
)

// SaveFile сохраняет файл на диск в указанную директорию
func SaveFile(fileHeader *multipart.FileHeader, uploadDir string) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("ошибка открытия файла: %v", err)
	}
	defer file.Close()

	filePath := filepath.Join(uploadDir, fileHeader.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("ошибка создания файла: %v", err)
	}
	defer out.Close()

	if _, err := out.ReadFrom(file); err != nil {
		return "", fmt.Errorf("ошибка сохранения файла: %v", err)
	}

	return filePath, nil
}
