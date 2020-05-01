package controllers

import (
	"github.com/ddalogin/siren/http/views"
	"github.com/ddalogin/siren/worker"
	"net/http"
	"strconv"
)

// Страница формы
type FormPage struct {
	Task        worker.Task
	TaskGraylog worker.TaskGraylog
	Message     string
}

// Страница запуска задачи
type RunPage struct {
	TaskResult worker.TaskResult
}

// Страница формы
func FormAction(w http.ResponseWriter, req *http.Request) {
	message := ""
	taskId := req.URL.Query().Get("id")
	task := worker.Task{}
	graylogTask := worker.TaskGraylog{}

	if taskId != "" && taskId != "0" {
		intTaskId, _ := strconv.Atoi(taskId)
		task = worker.GetTaskById(intTaskId)
		graylogTask = worker.GetTaskGraylogById(task.ObjectId)
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

	views.Render(w, "http/views/task/form.html", FormPage{
		Task:        task,
		TaskGraylog: graylogTask,
		Message:     message,
	})
}

// Страница выполнения задачи
func RunAction(w http.ResponseWriter, req *http.Request) {
	taskId := req.URL.Query().Get("id")

	if taskId == "" || taskId == "0" {
		http.NotFound(w, req)
		return
	}

	intTaskId, _ := strconv.Atoi(taskId)
	task := worker.GetTaskById(intTaskId)

	if (worker.Task{}) == task {
		http.NotFound(w, req)
		return
	}

	result := task.Do()

	views.Render(w, "http/views/task/result.html", RunPage{
		TaskResult: result,
	})
}
