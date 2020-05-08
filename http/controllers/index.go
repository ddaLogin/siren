package controllers

import (
	"github.com/ddalogin/siren/app/domain/model"
	"github.com/ddalogin/siren/container"
	"github.com/ddalogin/siren/http/views"
	"net/http"
)

type IndexController struct {
	Container *container.Container
}

// Главная страница
type IndexPage struct {
	Tasks []*model.Task
}

// Главная страница
func (c *IndexController) IndexAction(w http.ResponseWriter, req *http.Request) {
	views.Render(w, "http/views/index.html", IndexPage{
		Tasks: c.Container.TaskRepository().GetAll(),
	})
}
