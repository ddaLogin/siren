package main

import (
	"github.com/BurntSushi/toml"
	"github.com/ddalogin/siren/alert"
	"github.com/ddalogin/siren/database"
	"github.com/ddalogin/siren/http"
	"github.com/ddalogin/siren/worker"
	"io"
	"log"
	"os"
)

type Config struct {
	Http     http.Config
	Db       database.Config
	Telegram alert.Config
	Graylog  worker.GraylogConfig
}

func main() {
	initLogger()

	config := loadConfig()
	database.InitDatabase(config.Db)
	alert.InitTelegram(config.Telegram)
	worker.InitGraylogConfig(config.Graylog)
	go worker.StartWorker()
	http.StartServer(config.Http)
}

// Инициализация логов
func initLogger() {
	f, err := os.OpenFile("errors.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)
}

// Чтение настроек
func loadConfig() (config Config) {
	if _, err := toml.DecodeFile("config/config.toml", &config); err != nil {
		log.Fatal(err)
	}

	return
}
