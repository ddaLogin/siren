package views

import (
	"html/template"
	"net/http"
)

// Базовый темплейт
type baseTemplate struct {
	Content     interface{}
	TelegramBot string
	NotifyStart string
	NotifyEnd   string
}

var TelegramBot string
var NotifyStart string
var NotifyEnd string

// Выполнение указанного шаблона
func Render(w http.ResponseWriter, templateFile string, data interface{}) {

	base := baseTemplate{
		Content:     data,
		TelegramBot: TelegramBot,
		NotifyStart: NotifyStart,
		NotifyEnd:   NotifyEnd,
	}

	view, _ := template.New("").ParseFiles(templateFile, "http/views/base.html")

	view.ExecuteTemplate(w, "base", base)
}
