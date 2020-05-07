package telegram

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Ответ от телеграма
type Response struct {
	Response []byte
	Status   int
}

// Формат обновлений телеграм бота
type Updates struct {
	Result []struct {
		Message struct {
			Chat struct {
				Id       json.Number `json:"id"`
				Username string      `json:"username"`
				Type     string      `json:"type"`
			} `json:"chat"`
		} `json:"message"`
	} `json:"result"`
}

// Клиент для телеграма
type Client struct {
	host  string
	token string
	proxy string
}

// Конструктор телеграм клиента
func NewClient(host string, token string, proxy string) *Client {
	return &Client{host: host, token: token, proxy: proxy}
}

// Отправить сообщение в чат
func (c *Client) SendMessage(message string, chatId string) bool {
	apiUrl := c.host + "/bot" + c.token + "/sendMessage"
	query := []byte(`{"chat_id": "` + chatId + `", "text": "` + message + `", "parse_mode": "Markdown"}`)

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(query))
	if err != nil {
		log.Println("Ошибка при формирование запроса на отправку сообщения в телеграм", err)
		return false
	}

	response := c.sendRequest(req)

	return response != nil
}

// Получает обновления телеграм чатов
func (c *Client) UpdateChats() *Updates {
	apiUrl := c.host + "/bot" + c.token + "/getUpdates"

	req, err := http.NewRequest("GET", apiUrl, bytes.NewBuffer([]byte(``)))
	if err != nil {
		log.Println("Ошибка при формирование запроса за обновлениями телеграма", err)
		return nil
	}

	response := c.sendRequest(req)
	updates := Updates{}

	err = json.Unmarshal(response.Response, &updates)
	if err != nil {
		log.Println("Не удалось распарсить обновления телеграм бота", err, string(response.Response))
		return nil
	}

	return &updates
}

// Отправить запрос
func (c *Client) sendRequest(req *http.Request) *Response {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 20}

	if c.proxy != "" {
		proxyUrl, _ := url.Parse(c.proxy)
		client = &http.Client{Timeout: time.Second * 20, Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Ошибка при отправке запроса в телеграм", err)
		return nil
	}
	defer resp.Body.Close()

	response := Response{}

	response.Status = resp.StatusCode
	response.Response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Ошибка при чтение ответа от телеграма", err)
		return nil
	}

	return &response
}
