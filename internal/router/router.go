package router

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"statements/internal/config"
	"statements/internal/database"
	"statements/internal/handlers"
	"statements/internal/middleware"
)

// RegisterRoutes регистрирует маршруты с использованием Gin
func RegisterRoutes(cfg *config.Config) *gin.Engine {
	// Инициализируем роутер
	router := gin.Default()

	// Добавляем middleware для CORS и обработки ошибок
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandling())

	// Добавляем JWT аутентификацию
	router.Use(middleware.AuthMiddleware()) // Защищённые маршруты требуют JWT

	// Регистрация маршрутов
	registerStaticRoutes(router)
	registerAPIRoutes(router, cfg, database.DB) // Используем глобальный объект базы данных
	registerFileUploadRoutes(router, cfg)
	registerDownloadRoutes(router) // Новый маршрут для скачивания Excel

	// Статические файлы
	router.Static("/assets", cfg.FileUpload.StaticDir)

	return router
}

// registerStaticRoutes регистрирует маршруты для статических страниц
func registerStaticRoutes(router *gin.Engine) {
	static := router.Group("/")
	{
		static.GET("/", handlers.HandleHomePageGin)
		static.GET("/add-contract", handlers.HandleAddContractPage)
		static.POST("/submit-contract", func(c *gin.Context) {
			handlers.HandleContractSubmission(c, nil, database.DB)
		})
	}
}

// registerAPIRoutes регистрирует маршруты для API
func registerAPIRoutes(router *gin.Engine, cfg *config.Config, db *sql.DB) {
	api := router.Group("/api/v1")
	{
		// Пример API эндпоинта для получения списка контрагентов
		api.GET("/counterparties", func(c *gin.Context) {
			handlers.HandleCounterpartiesList(c, db)
		})
	}
}

// registerFileUploadRoutes регистрирует маршруты для загрузки файлов
func registerFileUploadRoutes(router *gin.Engine, cfg *config.Config) {
	upload := router.Group("/upload")
	{
		upload.POST("/", func(c *gin.Context) {
			handlers.HandleFileUploadGin(c, cfg)
		})
	}
}

// registerDownloadRoutes регистрирует маршрут для скачивания Excel-файлов
func registerDownloadRoutes(router *gin.Engine) {
	router.GET("/download", handlers.HandleDownloadTransactionsExcel)
}
