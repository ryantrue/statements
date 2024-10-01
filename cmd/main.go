package main

import (
	"fmt"
	"log"
	"os"
	"statements/internal/config"
	"statements/internal/database"
	"statements/internal/handlers"
	"statements/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.UpdateDependencies()
	startApp()
}

// startApp запускает основную логику приложения
func startApp() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Создание базы данных, если она не существует
	if err := database.CreateDatabaseIfNotExists(cfg); err != nil {
		log.Fatalf("Ошибка создания базы данных: %v", err)
	}

	// Подключение к базе данных
	if err := database.ConnectDB(cfg); err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer func() {
		if err := database.CloseDB(); err != nil {
			log.Printf("Ошибка при закрытии базы данных: %v", err)
		}
	}()

	// Выполнение миграций базы данных
	database.RunMigrations(cfg)

	// Создаем директорию для загрузки файлов, если её нет
	if err := os.MkdirAll(cfg.FileUpload.UploadDir, os.ModePerm); err != nil {
		log.Fatalf("Ошибка создания директории для загрузки файлов: %v", err)
	}

	// Регистрация маршрутов с использованием Gin
	router := registerRoutes(cfg)

	// Запуск сервера
	address := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Сервер запущен на http://%s\n", address)
	if err := router.Run(address); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

// registerRoutes регистрирует маршруты с использованием Gin
func registerRoutes(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// Главная страница
	router.GET("/", func(c *gin.Context) {
		handlers.HandleHomePageGin(c)
	})

	// Загрузка файлов
	router.POST("/upload", func(c *gin.Context) {
		handlers.HandleFileUploadGin(c, cfg)
	})

	// Статические файлы (CSS, JS)
	router.Static("/static", cfg.FileUpload.StaticDir)

	return router
}
