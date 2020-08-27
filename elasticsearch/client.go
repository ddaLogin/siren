package elasticsearch

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Формат ответа эластика
type ElasticResponse struct {
	Hits struct {
		Count int `json:"total"`
	} `json:"hits"`
}

// Клиент для ElasticSearch
type Client struct {
	host string
}

// Конструктор ElasticSearch клиента
func NewClient(host string) *Client {
	return &Client{host: host}
}

// Поиск
func (c *Client) Search(pattern string, aggTime string) (response ElasticResponse, err error) {
	pattern = strings.ReplaceAll(pattern, `"`, `\"`)

	var jsonStr = []byte(`{
		"query": {
			"bool": {
				"must": {
					"query_string": {
						"query": "` + pattern + `"
					}
				},
				"filter": {
					"bool": {
						"must": {
							"range": {
								"timestamp": {
									"gte": "now-` + aggTime + `",
									"lte": "now"
								}
							}
						}
					}
				}
			}
		}
	}`)
	req, err := http.NewRequest("POST", c.host+"_search", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Println("Не удалось собрать запрос в ElasticSearch.", err, req)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Не удалось выполнить запрос ElasticSearch.", err, req)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Не прочитать ответ от ElasticSearch.", err, string(body))
		return
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Не удалось распарсить ответ от ElasticSearch.", err, string(body))
		return
	}

	return
}
