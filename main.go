package main

import (
	"github.com/BurntSushi/toml"
	"github.com/ddalogin/siren/app/domain/service"
	"github.com/ddalogin/siren/container"
	"github.com/ddalogin/siren/database"
	"github.com/ddalogin/siren/http"
	"github.com/ddalogin/siren/http/views"
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

	cnt := container.NewContainer(config.Graylog, config.Notify, connector)

	wrk := worker.NewWorker(cnt.TaskService(), cnt.TaskRepository(), cnt.NotifyService())
	go wrk.Run()

	views.TelegramBot = config.Notify.Telegram.Bot
	views.NotifyStart = config.Notify.Start
	views.NotifyEnd = config.Notify.End

	server := http.NewServer(config.Http, cnt)
	server.Run()
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
