package worker

import (
	"github.com/ddalogin/siren/app/domain/model"
	"github.com/ddalogin/siren/app/domain/repository"
	"github.com/ddalogin/siren/app/domain/service"
	"time"
)

// Фоновый воркер задач
type Worker struct {
	reportTime     string
	taskService    *service.TaskService
	taskRepository *repository.TasksRepository
	notifyService  *service.NotifyService
	reportService  *service.ReportService
}

// Конструктор воркера
func NewWorker(reportTime string, taskService *service.TaskService, taskRepository *repository.TasksRepository, notifyService *service.NotifyService, reportService *service.ReportService) *Worker {
	return &Worker{reportTime: reportTime, taskService: taskService, taskRepository: taskRepository, notifyService: notifyService, reportService: reportService}
}

// Запустить воркер
func (w *Worker) Run() {
	ticker := time.Tick(time.Minute)

	for now := range ticker {
		if now.Format("15:04") == w.reportTime {
			go w.doReport()
		}

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

// Сформировать отчет и отправить
func (w *Worker) doReport() {
	reportMessage := w.reportService.MakeReport()

	if reportMessage != nil {
		w.notifyService.NotifyTelegram(reportMessage)
	}
}
