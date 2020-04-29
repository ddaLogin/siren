package worker

import (
	"time"
)

type Config struct {
}

// Запустить воркера
func StartWorker(config Config) {
	c := time.Tick(time.Minute)
	for now := range c {
		runTasksByTime(now)
		runTasksByInterval(now)
	}
}

// Запуск комманд по времени
func runTasksByTime(now time.Time) {
	tasks := GetAllTasksByTime(now.Format("15:04"))

	for _, task := range tasks {
		go task.Do()
	}
}

func runTasksByInterval(now time.Time) {
	tasks := GetAllTasksByInterval(now.Minute())

	for _, task := range tasks {
		go task.Do()
	}
}
