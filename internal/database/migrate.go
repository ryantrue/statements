package database

import (
	"log"
	"statements/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations выполняет миграции базы данных
func RunMigrations(cfg *config.Config) {
	// Подключаем драйвер postgres для миграций
	driver, err := postgres.WithInstance(DB, &postgres.Config{})
	if err != nil {
		log.Fatalf("Ошибка создания драйвера миграции: %v", err)
	}

	// Используем путь к миграциям из конфигурации
	migrationDir := cfg.Database.MigrationsDir

	// Создаем экземпляр мигратора
	m, err := migrate.NewWithDatabaseInstance(migrationDir, "postgres", driver)
	if err != nil {
		log.Fatalf("Ошибка создания мигратора: %v", err)
	}

	// Выполняем миграции вверх
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Ошибка выполнения миграций: %v", err)
	}

	log.Println("Миграции успешно выполнены")
}
