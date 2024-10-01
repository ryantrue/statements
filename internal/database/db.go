package database

import (
	"database/sql"
	"fmt"
	"log"
	"statements/internal/config"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var DB *sql.DB

// CreateDatabaseIfNotExists проверяет, существует ли база данных, и создает её, если не существует.
func CreateDatabaseIfNotExists(cfg *config.Config) error {
	// Подключаемся к базе данных, указанной в конфигурации
	contractkeepConnStr := cfg.Database.URL
	contractkeepDB, err := sql.Open("pgx", contractkeepConnStr)
	if err != nil {
		return fmt.Errorf("ошибка подключения к базе данных contractkeep: %w", err)
	}
	defer func() {
		if err := contractkeepDB.Close(); err != nil {
			log.Printf("Ошибка закрытия соединения с базой данных contractkeep: %v", err)
		}
	}()

	// Проверяем, существует ли база данных
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)"
	err = contractkeepDB.QueryRow(query, cfg.Database.DatabaseName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ошибка проверки существования базы данных: %w", err)
	}

	if exists {
		log.Printf("База данных '%s' уже существует.\n", cfg.Database.DatabaseName)
		return nil
	}

	// Создаем базу данных, если она не существует
	createDBQuery := fmt.Sprintf("CREATE DATABASE %s", cfg.Database.DatabaseName)
	_, err = contractkeepDB.Exec(createDBQuery)
	if err != nil {
		return fmt.Errorf("ошибка создания базы данных '%s': %w", cfg.Database.DatabaseName, err)
	}

	log.Printf("База данных '%s' успешно создана.\n", cfg.Database.DatabaseName)
	return nil
}

// ConnectDB устанавливает соединение с базой данных PostgreSQL через database/sql
func ConnectDB(cfg *config.Config) error {
	var err error
	DB, err = sql.Open("pgx", cfg.Database.URL)
	if err != nil {
		return fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	// Настраиваем пул соединений
	DB.SetMaxOpenConns(cfg.Database.MaxConnections)     // Устанавливаем максимальное количество открытых соединений
	DB.SetMaxIdleConns(cfg.Database.MaxIdleConnections) // Устанавливаем максимальное количество неактивных соединений

	// Проверяем соединение
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("ошибка проверки соединения с базой данных: %w", err)
	}

	log.Println("Успешное подключение к базе данных PostgreSQL!")
	return nil
}

// CloseDB закрывает соединение с базой данных
func CloseDB() error {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			return fmt.Errorf("ошибка закрытия соединения с базой данных: %w", err)
		}
		log.Println("Соединение с базой данных успешно закрыто.")
	}
	return nil
}
