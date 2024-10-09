package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// ErrorHandling middleware для обработки ошибок с улучшенным логированием
func ErrorHandling() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Проверка на наличие ошибок
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				// Логирование ошибок с различными уровнями
				switch err.Type {
				case gin.ErrorTypePublic:
					logrus.Warnf("Публичная ошибка: %s", err.Error())
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				case gin.ErrorTypeBind:
					logrus.Errorf("Ошибка биндинга данных: %s", err.Error())
					c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные запроса"})
				default:
					logrus.Errorf("Внутренняя ошибка сервера: %s", err.Error())
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
				}
			}
			return
		}

		// Обработка 404 ошибки
		if c.Writer.Status() == http.StatusNotFound {
			logrus.Warn("Страница не найдена")
			c.JSON(http.StatusNotFound, gin.H{"error": "Страница не найдена"})
		}
	}
}

// CORSMiddleware добавляет поддержку CORS с возможностью настройки
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*") // Позволяет запросы со всех источников, можно настроить
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		// Обрабатываем preflight-запросы для CORS
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// Можно добавить дополнительные проверки или логику, если требуется
		c.Next()
	}
}

// AuthMiddleware проверяет наличие и валидность токена в заголовке Authorization
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			logrus.Warn("Токен авторизации отсутствует")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Необходим токен авторизации"})
			c.Abort()
			return
		}

		// Здесь можно добавить логику для проверки токена
		// Например, вызов функции для валидации токена

		c.Next()
	}
}
