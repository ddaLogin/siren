package views

import (
	"html/template"
	"net/http"
)

// Базовый темплейт
type baseTemplate struct {
	Content     interface{}
	NotifyStart string
	NotifyEnd   string
}

// Выполнение указанного шаблона
func Render(w http.ResponseWriter, templateFile string, data interface{}) {
	//start, end := worker.GetNotifyTime()

	base := baseTemplate{
		Content:     data,
		NotifyStart: "start",
		NotifyEnd:   "end",
	}

	view, _ := template.New("").ParseFiles(templateFile, "http/views/base.html")

	view.ExecuteTemplate(w, "base", base)
}
