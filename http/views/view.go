package views

import (
	"html/template"
	"net/http"
)

// Базовый темплейт
type baseTemplate struct {
	Content interface{}
}

// Выполнение указанного шаблона
func Render(w http.ResponseWriter, templateFile string, data interface{}) {
	base := baseTemplate{
		Content: data,
	}

	view, _ := template.New("").ParseFiles(templateFile, "http/views/base.html")

	view.ExecuteTemplate(w, "base", base)
}
