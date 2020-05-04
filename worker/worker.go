package worker

import (
	"github.com/ddalogin/siren/alert"
	"time"
)

// Настройки воркера
type Config struct {
	NotifyStart string // Время в формате "08:30", после которого уведомления можно расылать
	NotifyEnd   string // Время в формате "21:00", после которого уведомления нельзя расылать
}

var config Config
var silence = true

// Настройка воркера
func InitWorker(cfg Config) {
	config = cfg
}

// Возвращает время разрешенное для уведомлений
func GetNotifyTime() (start string, end string) {
	start = config.NotifyStart
	end = config.NotifyEnd

	return
}

// Запустить воркера
func StartWorker() {
	ticker := time.Tick(time.Minute)
	for now := range ticker {
		checkSilence(now)
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
	if result.Status != STATUS_SILENT && !silence {
		alert.SendMessage(result.GetMessage())
	}
}

// Проверка тихого режима
func checkSilence(now time.Time) {
	notifyStart, _ := time.Parse("15:04", config.NotifyStart)
	notifyEnd, _ := time.Parse("15:04", config.NotifyEnd)
	current, _ := time.Parse("15:04", now.Format("15:04"))

	silence = !(current.After(notifyStart) && current.Before(notifyEnd))
}
