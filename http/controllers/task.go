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

//// Страница результата выполнения задачи
//type ResultPage struct {
//	TaskResult *model.ResultGraylog
//}
//
//// Страница всех результатов выполнения задачи
//type ResultListPage struct {
//	Task        *model.Task
//	TaskResults []*model.ResultGraylog
//}

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
func RunAction(w http.ResponseWriter, req *http.Request) {
	//taskId := req.URL.Query().Get("id")
	//
	//if taskId == "" || taskId == "0" {
	//	http.NotFound(w, req)
	//	return
	//}
	//
	//intTaskId, _ := strconv.Atoi(taskId)
	//task := worker.GetTaskById(intTaskId)
	//
	//if (worker.Task{}) == task {
	//	http.NotFound(w, req)
	//	return
	//}
	//
	//result := task.Do()
	//
	//http.Redirect(w, req, "/task/result?id="+strconv.Itoa(int(result.Id)), http.StatusSeeOther)
}

// Страница результата выполнения задачи
func ResultAction(w http.ResponseWriter, req *http.Request) {
	//resultId := req.URL.Query().Get("id")
	//
	//if resultId == "" || resultId == "0" {
	//	http.NotFound(w, req)
	//	return
	//}
	//
	//intResultId, _ := strconv.Atoi(resultId)
	//result := worker.GetResultById(intResultId)
	//
	//if (worker.TaskResult{}) == result {
	//	http.NotFound(w, req)
	//	return
	//}
	//
	//views.Render(w, "http/views/task/result.html", ResultPage{
	//	TaskResult: result,
	//})
}

// Страница всех результатов выполнения задачи
func ResultListAction(w http.ResponseWriter, req *http.Request) {
	//taskId := req.URL.Query().Get("id")
	//
	//if taskId == "" || taskId == "0" {
	//	http.NotFound(w, req)
	//	return
	//}
	//
	//intTaskId, _ := strconv.Atoi(taskId)
	//task := worker.GetTaskById(intTaskId)
	//
	//if (worker.Task{}) == task {
	//	http.NotFound(w, req)
	//	return
	//}
	//
	//results := worker.GetResultByTaskId(intTaskId, 100)
	//
	//views.Render(w, "http/views/task/result_list.html", ResultListPage{
	//	Task:        task,
	//	TaskResults: results,
	//})
}
