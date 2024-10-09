package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	baseTemplate = "base.html"
	menuTemplate = "menu.html"
	templatesDir = "/app/assets/templates/"
	pagesDir     = "/app/assets/pages/"
)

// renderTemplate рендерит HTML-шаблон с поддержкой базового шаблона и отправляет его пользователю
func renderTemplate(c *gin.Context, pageTemplate string, data gin.H) {
	// Полные пути к файлам шаблонов
	templateFiles := []string{
		filepath.Join(templatesDir, baseTemplate), // базовый шаблон
		filepath.Join(templatesDir, menuTemplate), // меню
		filepath.Join(pagesDir, pageTemplate),     // дочерний шаблон страницы
	}

	// Парсим шаблоны
	tmpl, err := template.ParseFiles(templateFiles...)
	if err != nil {
		logrus.Errorf("Ошибка загрузки шаблонов: %v", err)
		c.String(http.StatusInternalServerError, "Ошибка загрузки шаблона")
		return
	}

	// Рендерим шаблон
	if err := tmpl.Execute(c.Writer, data); err != nil {
		logrus.Errorf("Ошибка рендеринга шаблона: %v", err)
		c.String(http.StatusInternalServerError, "Ошибка рендеринга шаблона")
	}
}

// HandleHomePageGin обрабатывает запрос на главную страницу
func HandleHomePageGin(c *gin.Context) {
	renderTemplate(c, "index.html", gin.H{
		"Title":  "Главная страница",
		"Header": "Загрузка файлов и скачивание данных",
	})
}

// HandleAddContractPage обрабатывает запрос на страницу добавления контракта
func HandleAddContractPage(c *gin.Context) {
	renderTemplate(c, "add_contract.html", gin.H{
		"Title":  "Добавление контракта",
		"Header": "Добавление контракта",
	})
}
