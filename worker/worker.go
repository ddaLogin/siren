package worker

import (
	"github.com/ddalogin/siren/app/domain/model"
	"github.com/ddalogin/siren/app/domain/repository"
	"github.com/ddalogin/siren/app/domain/service"
	"time"
)

// Фоновый воркер задач
type Worker struct {
	taskService    *service.TaskService
	taskRepository *repository.TasksRepository
	notifyService  *service.NotifyService
}

// Конструктор воркера
func NewWorker(taskService *service.TaskService, taskRepository *repository.TasksRepository, notifyService *service.NotifyService) *Worker {
	return &Worker{taskService: taskService, taskRepository: taskRepository, notifyService: notifyService}
}

// Запустить воркер
func (w *Worker) Run() {
	ticker := time.Tick(time.Second)

	for now := range ticker {
		tasks := w.taskRepository.GetForRun(now)

		for _, task := range tasks {
			go w.doTask(task)
		}
	}
}

// Выполняет задачу и отправляет уведомление
func (w *Worker) doTask(task *model.Task) {
	result := w.taskService.RunTask(task)

	task.SetNextTime(task.CalculateNextTime())
	w.taskRepository.Save(task)

	if result.IsNeedNotify() {
		w.notifyService.NotifyTelegram(result.BuildTelegramNotify())
	}
}
