package worker

import (
	"github.com/ddalogin/siren/alert"
	"time"
)

type Config struct {
}

// Запустить воркера
func StartWorker(config Config) {
	ticker := time.Tick(time.Minute)
	for now := range ticker {
		runTasksByTime(now)
		runTasksByInterval(now)
	}
}

// Запуск комманд по времени
func runTasksByTime(now time.Time) {
	tasks := GetAllTasksByTime(now.Format("15:04"))

	for _, task := range tasks {
		go runTask(task)
	}
}

// Зпуск команд по интервалу
func runTasksByInterval(now time.Time) {
	tasks := GetAllTasksByInterval(now.Minute())

	for _, task := range tasks {
		go runTask(task)
	}
}

// Выполнение задачи
func runTask(task Task) {
	result := task.Do()
	parseResult(result)
}

// Чтение результатов выполнения задачи
func parseResult(result TaskResult) {
	if result.Status != STATUS_SILENT {
		alert.SendMessage(result.GetMessage())
	}
}
