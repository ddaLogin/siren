package controllers

import (
	"github.com/ddalogin/siren/app/domain/model"
	"github.com/ddalogin/siren/container"
	"github.com/ddalogin/siren/http/views"
	"net/http"
	"strconv"
)

type TaskController struct {
	Container *container.Container
}

// Страница формы
type FormPage struct {
	Task        *model.Task
	TaskGraylog *model.TaskGraylog
	Message     string
}

// Страница формы
func (c *TaskController) FormAction(w http.ResponseWriter, req *http.Request) {
	message := ""
	taskId := req.URL.Query().Get("id")
	task := &model.Task{}
	taskGraylog := &model.TaskGraylog{}

	if taskId != "" {
		intTaskId, _ := strconv.Atoi(taskId)
		task = c.Container.TaskRepository().GetById(intTaskId)

		if task == nil || *task == (model.Task{}) {
			http.NotFound(w, req)
			return
		}

		taskGraylog = c.Container.TaskGraylogRepository().GetById(task.ObjectId())

		if taskGraylog == nil || *taskGraylog == (model.TaskGraylog{}) {
			http.NotFound(w, req)
			return
		}
	}

	if http.MethodPost == req.Method {
		task.SetTitle(req.FormValue("title"))
		task.SetIsEnabled("true" == req.FormValue("enabled"))

		interval, _ := strconv.Atoi(req.FormValue("interval"))
		task.SetInterval(interval)

		usernames := req.FormValue("usernames")

		if usernames == "" {
			task.SetUsernames(nil)
		} else {
			task.SetUsernames(&usernames)
		}

		taskGraylog.SetPattern(req.FormValue("pattern"))
		taskGraylog.SetAggregateTime(req.FormValue("aggregate_time"))
		min, _ := strconv.Atoi(req.FormValue("min"))
		taskGraylog.SetMin(min)
		max, _ := strconv.Atoi(req.FormValue("max"))
		taskGraylog.SetMax(max)

		if c.Container.TaskGraylogRepository().Save(taskGraylog) {
			task.SetObjectType(model.TYPE_GRAYLOG)
			task.SetObjectId(taskGraylog.Id())

			if c.Container.TaskRepository().Save(task) {
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
		TaskGraylog: taskGraylog,
		Message:     message,
	})
}

// Метод удаления задачи
func (c *TaskController) DeleteAction(w http.ResponseWriter, req *http.Request) {
	taskId := req.URL.Query().Get("id")

	if taskId == "" || taskId == "0" || http.MethodPost != req.Method {
		http.NotFound(w, req)
		return
	}

	intTaskId, _ := strconv.Atoi(taskId)
	task := c.Container.TaskRepository().GetById(intTaskId)

	if task == nil || *task == (model.Task{}) {
		http.NotFound(w, req)
		return
	}

	c.Container.TaskService().DeleteTask(task)

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

// Страница выполнения задачи
func (c *TaskController) RunAction(w http.ResponseWriter, req *http.Request) {
	taskId := req.URL.Query().Get("id")

	if taskId == "" || taskId == "0" {
		http.NotFound(w, req)
		return
	}

	intTaskId, _ := strconv.Atoi(taskId)
	task := c.Container.TaskRepository().GetById(intTaskId)

	if (&model.Task{}) == task {
		http.NotFound(w, req)
		return
	}

	if task.ObjectType() == model.TYPE_GRAYLOG {
		result := c.Container.GraylogService().RunTaskGraylog(task)

		if result == (model.ResultGraylog{}) {
			http.NotFound(w, req)
			return
		}

		http.Redirect(w, req, "/task/graylog/result?id="+strconv.Itoa(result.Id()), http.StatusSeeOther)
		return
	}

	http.NotFound(w, req)
	return
}

// Страница всех результатов выполнения задачи
func (c *TaskController) ResultListAction(w http.ResponseWriter, req *http.Request) {
	taskId := req.URL.Query().Get("id")

	if taskId == "" || taskId == "0" {
		http.NotFound(w, req)
		return
	}

	intTaskId, _ := strconv.Atoi(taskId)
	task := c.Container.TaskRepository().GetById(intTaskId)

	if (&model.Task{}) == task {
		http.NotFound(w, req)
		return
	}

	if task.ObjectType() == model.TYPE_GRAYLOG {
		http.Redirect(w, req, "/task/graylog/result/list?id="+strconv.Itoa(task.ObjectId()), http.StatusSeeOther)
		return
	}

	http.NotFound(w, req)
	return
}
