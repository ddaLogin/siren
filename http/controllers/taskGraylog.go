package controllers

import (
	"github.com/ddalogin/siren/app/domain/model"
	"github.com/ddalogin/siren/container"
	"github.com/ddalogin/siren/http/views"
	"net/http"
	"strconv"
)

type TaskGraylogController struct {
	Container *container.Container
}

// Страница результата
type ResultPage struct {
	Result *model.ResultGraylog
}

// Страница результатов
type ResultListPage struct {
	Results []*model.ResultGraylog
}

// Страница просмотра одного результата грейлог задачи
func (c *TaskGraylogController) ResultAction(w http.ResponseWriter, req *http.Request) {
	resultId := req.URL.Query().Get("id")

	if resultId == "" || resultId == "0" {
		http.NotFound(w, req)
		return
	}

	intResultId, _ := strconv.Atoi(resultId)
	result := c.Container.ResultsGraylogRepository().GetById(intResultId)

	if result == nil || *result == (model.ResultGraylog{}) {
		http.NotFound(w, req)
		return
	}

	views.Render(w, "http/views/taskGraylog/result.html", ResultPage{
		Result: result,
	})
}

// Страница просмотра последних результатов грейлог задачи
func (c *TaskGraylogController) ResultListAction(w http.ResponseWriter, req *http.Request) {
	graylogTaskId := req.URL.Query().Get("id")

	if graylogTaskId == "" || graylogTaskId == "0" {
		http.NotFound(w, req)
		return
	}

	intGraylogTaskId, _ := strconv.Atoi(graylogTaskId)
	results := c.Container.ResultsGraylogRepository().GetByTaskGraylogId(intGraylogTaskId, 100)

	if results == nil || len(results) == 0 {
		http.NotFound(w, req)
		return
	}

	views.Render(w, "http/views/taskGraylog/result_list.html", ResultListPage{
		Results: results,
	})
}
