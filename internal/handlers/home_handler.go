package handlers

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleHomePageGin обрабатывает запросы к главной странице через Gin
func HandleHomePageGin(c *gin.Context) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка загрузки шаблона")
		return
	}
	if err := tmpl.Execute(c.Writer, nil); err != nil {
		c.String(http.StatusInternalServerError, "Ошибка рендеринга шаблона")
	}
}
