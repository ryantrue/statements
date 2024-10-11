package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/jwtauth"
	"go.uber.org/zap"
)

// Инициализация JWT с секретным ключом
var tokenAuth *jwtauth.JWTAuth

func InitJWT(secret string) {
	tokenAuth = jwtauth.New("HS256", []byte(secret), nil)
}

// Инициализация Zap логгера
var logger *zap.Logger

func InitLogger() {
	var err error
	logger, err = zap.NewProduction() // Логгер для продакшна
	if err != nil {
		panic("Не удалось инициализировать логгер: " + err.Error())
	}
	defer logger.Sync() // Запись всех буферов перед завершением
}

// CORSMiddleware добавляет поддержку CORS с возможностью настройки
func CORSMiddleware() gin.HandlerFunc {
	config := cors.Config{
		AllowOrigins:     []string{"https://example.com"}, // Здесь можно указать свои разрешённые домены
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	return cors.New(config)
}

// AuthMiddleware проверяет и валидирует JWT токен
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, claims, err := jwtauth.FromContext(c.Request.Context())
		if err != nil {
			logger.Warn("Невалидный токен", zap.String("path", c.Request.URL.Path))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Невалидный токен"})
			c.Abort()
			return
		}

		logger.Info("Токен валиден", zap.Any("claims", claims))
		c.Set("claims", claims)
		c.Next()
	}
}

// HandleError логирует и возвращает ошибки с использованием Zap
func HandleError(c *gin.Context, err *gin.Error) {
	switch err.Type {
	case gin.ErrorTypePublic:
		logger.Warn("Публичная ошибка", zap.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case gin.ErrorTypeBind:
		logger.Error("Ошибка биндинга данных", zap.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные запроса"})
	default:
		logger.Error("Внутренняя ошибка сервера", zap.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
	}
}

// ErrorHandling middleware для обработки ошибок с улучшенным логированием через Zap
func ErrorHandling() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				HandleError(c, err)
			}
			return
		}

		if c.Writer.Status() == http.StatusNotFound {
			logger.Warn("Страница не найдена", zap.String("path", c.Request.URL.Path))
			c.JSON(http.StatusNotFound, gin.H{"error": "Страница не найдена"})
		}
	}
}
