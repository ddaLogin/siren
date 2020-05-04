package alert

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Интерфейс сообщения
type TelegramMessageInterface interface {
	GetMessage() string // Получить сообщение для телеграмма
}

// Настройка бота
type Config struct {
	Token  string
	Host   string
	ChatId string
	Proxy  string
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
	apiUrl := config.Host + "/bot" + config.Token + "/sendMessage"
	query := []byte(`{"chat_id": "` + config.ChatId + `", "text": "` + message + `", "parse_mode": "Markdown"}`)
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(query))
	if err != nil {
		log.Println("Error send telegram message. ", err)
		return
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	proxyUrl, err := url.Parse(config.Proxy)
	client := &http.Client{Timeout: time.Second * 20, Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error response reading from telegram. ", err)
		return
	}
	defer resp.Body.Close()

	response := Response{}

	// Parse response
	response.Status = resp.StatusCode
	response.Response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error body reading from telegram. ", err)
		return
	}

	return
}
