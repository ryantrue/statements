package utils

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

// UpdateDependencies обновляет зависимости для Go и Python только если есть изменения
func UpdateDependencies() {
	// Шаг 1: Устанавливаем/обновляем Go-зависимости, если есть изменения в go.mod
	fmt.Println("Проверяем необходимость обновления Go-зависимостей...")
	if err := updateGoDependenciesIfNeeded(); err != nil {
		log.Printf("Ошибка обновления Go-зависимостей: %v", err)
	}

	// Шаг 2: Устанавливаем/обновляем Python-зависимости, если есть изменения в requirements.txt
	fmt.Println("Проверяем необходимость обновления Python-зависимостей...")
	if err := updatePythonDependenciesIfNeeded(); err != nil {
		log.Printf("Ошибка обновления Python-зависимостей: %v", err)
	}
}

// updateGoDependenciesIfNeeded обновляет Go-зависимости только если есть изменения в go.mod
func updateGoDependenciesIfNeeded() error {
	if !hasFileChanged("go.mod") {
		fmt.Println("Изменений в Go-зависимостях не найдено, пропускаем обновление.")
		return nil
	}

	fmt.Println("Обновляем Go-зависимости...")
	if err := updateGoDependencies(); err != nil {
		return errors.New("Ошибка обновления Go-зависимостей: " + err.Error())
	}

	fmt.Println("Go-зависимости успешно обновлены!")
	return nil
}

// updateGoDependencies выполняет обновление и установку Go-зависимостей
func updateGoDependencies() error {
	if err := runCommand("go", "mod", "tidy"); err != nil {
		return errors.New("ошибка выполнения go mod tidy: " + err.Error())
	}

	if err := runCommand("go", "get", "-u", "./..."); err != nil {
		return errors.New("ошибка выполнения go get: " + err.Error())
	}

	return nil
}

// updatePythonDependenciesIfNeeded обновляет Python-зависимости, если есть изменения в requirements.txt
func updatePythonDependenciesIfNeeded() error {
	if !isPythonInstalled() {
		fmt.Println("Python не установлен, пропускаем обновление Python-зависимостей.")
		return nil
	}

	if !hasFileChanged("requirements.txt") {
		fmt.Println("Изменений в Python-зависимостях не найдено, пропускаем обновление.")
		return nil
	}

	fmt.Println("Обновляем Python-зависимости...")
	if err := updatePythonDependencies(); err != nil {
		return errors.New("Ошибка обновления Python-зависимостей: " + err.Error())
	}

	fmt.Println("Python-зависимости успешно обновлены!")
	return nil
}

// updatePythonDependencies выполняет установку Python-зависимостей из requirements.txt
func updatePythonDependencies() error {
	if err := runCommand("pip3", "install", "-r", "requirements.txt"); err != nil {
		return errors.New("ошибка выполнения pip install: " + err.Error())
	}

	return nil
}

// isPythonInstalled проверяет, установлен ли Python
func isPythonInstalled() bool {
	_, err := exec.LookPath("python3")
	return err == nil
}

// runCommand выполняет команду в командной строке и обрабатывает ошибки
func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// hasFileChanged проверяет, изменился ли файл с момента последней проверки
// Используем хэширование файла (MD5) для определения изменений
func hasFileChanged(filePath string) bool {
	// Читаем файл
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Ошибка при открытии файла %s: %v", filePath, err)
		return false
	}
	defer file.Close()

	// Вычисляем MD5-хэш файла
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Printf("Ошибка при чтении файла %s: %v", filePath, err)
		return false
	}
	currentHash := fmt.Sprintf("%x", hash.Sum(nil))

	// Сохраняем и проверяем хэш файла для сравнения с предыдущей версией
	cacheFilePath := filePath + ".md5"
	if previousHash, err := os.ReadFile(cacheFilePath); err == nil {
		if string(previousHash) == currentHash {
			// Хэш не изменился, файл не изменялся
			return false
		}
	}

	// Хэш изменился или его нет — записываем новый хэш в файл
	if err := os.WriteFile(cacheFilePath, []byte(currentHash), 0644); err != nil {
		log.Printf("Ошибка при сохранении хэша файла %s: %v", filePath, err)
	}
	return true
}
