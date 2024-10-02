package handlers

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleHomePageGin обрабатывает запросы к главной странице через Gin
func HandleHomePageGin(c *gin.Context) {
	// Измените путь к шаблону на абсолютный путь внутри контейнера
	tmpl, err := template.ParseFiles("/app/templates/index.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка загрузки шаблона")
		return
	}
	if err := tmpl.Execute(c.Writer, nil); err != nil {
		c.String(http.StatusInternalServerError, "Ошибка рендеринга шаблона")
	}
}
