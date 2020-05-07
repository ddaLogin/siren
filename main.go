package main

import (
	"github.com/BurntSushi/toml"
	"github.com/ddalogin/siren/app/domain/repository"
	"github.com/ddalogin/siren/app/domain/service"
	"github.com/ddalogin/siren/database"
	"github.com/ddalogin/siren/http"
	"github.com/ddalogin/siren/worker"
	"log"
	"os"
)

type Config struct {
	Http    http.Config
	Db      database.Config
	Notify  service.NotifyConfig
	Graylog service.GraylogConfig
}

// Инициализация приложения
func init() {
	initLogger()
}

// Начало работы
func main() {
	config := loadConfig()

	connector := database.NewConnector(config.Db)
	defer connector.Close()

	taskRepository := repository.GetTasksRepository(connector)
	taskGraylogRepository := repository.GetTasksGraylogRepository(connector)
	resultsGraylogRepository := repository.GetResultsGraylogRepository(connector)
	telegramChatRepository := repository.GetTelegramChatRepository(connector)

	graylogService := service.NewGraylogService(config.Graylog, taskGraylogRepository, resultsGraylogRepository)
	taskService := service.NewTaskService(taskRepository, graylogService)
	notifyService := service.NewNotifyService(config.Notify, telegramChatRepository)

	wrk := worker.NewWorker(taskService, taskRepository, notifyService)
	wrk.Run()
}

// Инициализация логов
func initLogger() {
	f, err := os.OpenFile("errors.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Ошибка при открытие файла для логов: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
}

// Чтение настроек
func loadConfig() (config Config) {
	if _, err := toml.DecodeFile("config/config.toml", &config); err != nil {
		log.Fatal(err)
	}

	return
}
