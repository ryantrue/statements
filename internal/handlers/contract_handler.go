package handlers

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

// renderTemplate рендерит HTML-шаблон и отправляет его пользователю
func renderTemplate(c *gin.Context, templatePath string) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка загрузки шаблона")
		return
	}
	if err := tmpl.Execute(c.Writer, nil); err != nil {
		c.String(http.StatusInternalServerError, "Ошибка рендеринга шаблона")
	}
}

// HandleHomePageGin обрабатывает запросы к главной странице через Gin
func HandleHomePageGin(c *gin.Context) {
	renderTemplate(c, "/app/assets/index.html")
}

// HandleAddContractPage рендерит страницу для добавления контрактов
func HandleAddContractPage(c *gin.Context) {
	renderTemplate(c, "/app/assets/add_contract.html")
}
