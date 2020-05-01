package controllers

import (
	"github.com/ddalogin/siren/http/views"
	"github.com/ddalogin/siren/worker"
	"github.com/ddalogin/siren/worker/graylog"
	"net/http"
	"strconv"
)

// Главная страница
type IndexPage struct {
	Tasks []worker.Task
}

// Страница формы
type FormPage struct {
	Task        worker.Task
	TaskGraylog graylog.TaskGraylog
	Message     string
}

// Главная страница
func IndexAction(w http.ResponseWriter, req *http.Request) {
	views.Render(w, "http/views/index.html", IndexPage{
		Tasks: worker.GetAllTasks(),
	})
}

// Страница формы
func FormAction(w http.ResponseWriter, req *http.Request) {
	message := ""
	taskId := req.URL.Query().Get("id")
	task := worker.Task{}
	graylogTask := graylog.TaskGraylog{}

	if taskId != "" && taskId != "0" {
		intTaskId, _ := strconv.Atoi(taskId)
		task = worker.GetTaskById(intTaskId)
		graylogTask = graylog.GetTaskGraylogById(task.ObjectId)
	}

	if http.MethodPost == req.Method {
		title := req.FormValue("title")
		isEnable := req.FormValue("enable")
		interval, _ := strconv.Atoi(req.FormValue("interval"))
		time := req.FormValue("time")

		pattern := req.FormValue("pattern")
		aggTime := req.FormValue("agg_time")
		minCount, _ := strconv.Atoi(req.FormValue("min_count"))
		maxCount, _ := strconv.Atoi(req.FormValue("max_count"))

		task.Title = title
		task.IsEnable = "true" == isEnable
		task.Interval = interval
		task.Time = time

		graylogTask.Pattern = pattern
		graylogTask.AggTime = aggTime
		graylogTask.MinCount = &minCount
		graylogTask.MaxCount = &maxCount

		if graylogTask.Save() {
			task.ObjectType = 1
			task.ObjectId = graylogTask.Id

			if task.Save() {
				http.Redirect(w, req, "/", http.StatusSeeOther)
			} else {
				message = "Не удалось сохранить задачу"
			}
		} else {
			message = "Не удалось сохранить задачу"
		}
	}

	views.Render(w, "http/views/form.html", FormPage{
		Task:        task,
		TaskGraylog: graylogTask,
		Message:     message,
	})
}
