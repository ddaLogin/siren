package alert

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Настройка бота
type Config struct {
	Token  string
	Host   string
	ChatId string
}

// Ответ от телеграма
type Response struct {
	Response []byte
	Status   int
}

var config Config

// Установить конфиг
func InitTelegram(cfg Config) {
	config = cfg
}

// Отправить сообщение в чат
func SendMessage(message string) {
	url := config.Host + "/bot" + config.Token + "/sendMessage"
	query := []byte(`{"chat_id": "` + config.ChatId + `", "text": "` + message + `", "parse_mode": "Markdown"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
	if err != nil {
		log.Fatal("Error send telegram message. ", err)
		return
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 10}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error response reading from telegram. ", err)
		return
	}
	defer resp.Body.Close()

	response := Response{}

	// Parse response
	response.Status = resp.StatusCode
	response.Response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error body reading from telegram. ", err)
		return
	}

	return
}
