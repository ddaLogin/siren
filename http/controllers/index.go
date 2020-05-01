package controllers

import (
	"github.com/ddalogin/siren/http/views"
	"github.com/ddalogin/siren/worker"
	"net/http"
)

// Главная страница
type IndexPage struct {
	Tasks []worker.Task
}

// Главная страница
func IndexAction(w http.ResponseWriter, req *http.Request) {
	views.Render(w, "http/views/index.html", IndexPage{
		Tasks: worker.GetAllTasks(),
	})
}
