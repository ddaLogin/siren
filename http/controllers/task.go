package controllers

import (
	"net/http"
)

//// Страница формы
//type FormPage struct {
//	Task        worker.Task
//	TaskGraylog worker.TaskGraylog
//	Message     string
//}
//
//// Страница результата выполнения задачи
//type ResultPage struct {
//	TaskResult worker.TaskResult
//}
//
//// Страница всех результатов выполнения задачи
//type ResultListPage struct {
//	Task        worker.Task
//	TaskResults []worker.TaskResult
//}

// Страница формы
func FormAction(w http.ResponseWriter, req *http.Request) {
	//message := ""
	//taskId := req.URL.Query().Get("id")
	//task := worker.Task{}
	//graylogTask := worker.TaskGraylog{}
	//
	//if taskId != "" && taskId != "0" {
	//	intTaskId, _ := strconv.Atoi(taskId)
	//	task = worker.GetTaskById(intTaskId)
	//	graylogTask = worker.GetTaskGraylogById(task.ObjectId)
	//}
	//
	//if http.MethodPost == req.Method {
	//	title := req.FormValue("title")
	//	isEnable := req.FormValue("enable")
	//	interval, _ := strconv.Atoi(req.FormValue("interval"))
	//	time := req.FormValue("time")
	//
	//	pattern := req.FormValue("pattern")
	//	aggTime := req.FormValue("agg_time")
	//	minCount, _ := strconv.Atoi(req.FormValue("min_count"))
	//	maxCount, _ := strconv.Atoi(req.FormValue("max_count"))
	//
	//	task.Title = title
	//	task.IsEnable = "true" == isEnable
	//	task.Interval = interval
	//	task.Time = time
	//
	//	graylogTask.Pattern = pattern
	//	graylogTask.AggTime = aggTime
	//	graylogTask.MinCount = &minCount
	//	graylogTask.MaxCount = &maxCount
	//
	//	if graylogTask.Save() {
	//		task.ObjectType = 1
	//		task.ObjectId = graylogTask.Id
	//
	//		if task.Save() {
	//			http.Redirect(w, req, "/", http.StatusSeeOther)
	//		} else {
	//			message = "Не удалось сохранить задачу"
	//		}
	//	} else {
	//		message = "Не удалось сохранить задачу"
	//	}
	//}
	//
	//views.Render(w, "http/views/task/form.html", FormPage{
	//	Task:        task,
	//	TaskGraylog: graylogTask,
	//	Message:     message,
	//})
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

// Метод удаления задачи
func DeleteAction(w http.ResponseWriter, req *http.Request) {
	//taskId := req.URL.Query().Get("id")
	//
	//if taskId == "" || taskId == "0" || http.MethodPost != req.Method {
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
	//task.Delete()
	//
	//http.Redirect(w, req, "/", http.StatusSeeOther)
}
