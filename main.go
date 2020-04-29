package main

import (
	"github.com/BurntSushi/toml"
	"github.com/ddalogin/siren/alert"
	"github.com/ddalogin/siren/database"
	"github.com/ddalogin/siren/http"
	"github.com/ddalogin/siren/worker"
	"github.com/ddalogin/siren/worker/graylog"
	"io"
	"log"
	"os"
)

type Config struct {
	Http     http.Config
	Worker   worker.Config
	Db       database.Config
	Telegram alert.Config
	Eshost   string
}

func main() {
	initLogger()

	config := loadConfig()
	database.InitDatabase(config.Db)
	alert.InitTelegram(config.Telegram)
	graylog.InitTaskGraylog(config.Eshost)
	go worker.StartWorker(config.Worker)
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
