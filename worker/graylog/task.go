package graylog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ddalogin/siren/alert"
	"github.com/ddalogin/siren/database"
	"io/ioutil"
	"log"
	"net/http"
)

// Задача проверки сообщения в грейлоге
type TaskGraylog struct {
	Id       int64  // Идентификатор
	Pattern  string // Шаблон поиска
	AggTime  string // Время для аггрегации результатов
	MinCount *int   // Минимально допустимое кол-во сообщений
	MaxCount *int   // Максимально допустимое кол-во сообщений
}

type ElasticResponse struct {
	Hits struct {
		Count int `json:"total"`
	} `json:"hits"`
}

var host string

// Установить конфиг
func InitTaskGraylog(h string) {
	host = h
}

// Выполнение задачи
func (t TaskGraylog) Do(title string) {
	fmt.Println("Выполнение задачи: " + title)

	var jsonStr = []byte(`{
		"query": {
			"bool": {
				"must": {
					"query_string": {
						"query": "` + t.Pattern + `"
					}
				},
				"filter": {
					"bool": {
						"must": {
							"range": {
								"timestamp": {
									"gte": "now-` + t.AggTime + `",
									"lte": "now"
								}
							}
						}
					}
				}
			}
		}
	}`)
	req, err := http.NewRequest("POST", host, bytes.NewBuffer(jsonStr))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	result := ElasticResponse{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("Не удалось распарсить ответ от ElasticSearch. ", err, string(body))
		return
	}

	message := "\xF0\x9F\x9A\xA8 *" + title + "*\r\n"

	if t.MinCount != nil && result.Hits.Count <= *t.MinCount {
		message = message + "\xE2\xAC\x87 *Кол-во сообщей уменьшилось*\r\n"
		message = message + t.Pattern
		alert.SendMessage(message)
	}

	if t.MaxCount != nil && result.Hits.Count >= *t.MaxCount {
		message = message + "\xE2\xAC\x86 *Кол-во сообщей увеличилось*\r\n"
		message = message + t.Pattern
		alert.SendMessage(message)
	}
}

// Получить задачу по идентификатору
func GetTaskGraylogById(id int64) (task TaskGraylog) {
	db := database.Db()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM tasks_graylog WHERE id = ?", id)
	if err != nil {
		log.Println("Не удалось найти грейлог задачу по идентификатору. ", err)
		return
	}
	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(
			&task.Id,
			&task.Pattern,
			&task.AggTime,
			&task.MinCount,
			&task.MaxCount,
		)
		if err != nil {
			log.Println("Не удалось распарсить грейлог задачу по идентификатору. ", err)
			continue
		}
	}

	return
}

// Сохранить TaskGraylog
func (t *TaskGraylog) Save() bool {
	db := database.Db()
	defer db.Close()

	if t.Id == 0 {
		result, err := db.Exec("INSERT INTO tasks_graylog (pattern, agg_time, min_count, max_count) VALUES (?, ?, ?, ?)",
			t.Pattern, t.AggTime, t.MinCount, t.MaxCount)
		if err != nil {
			log.Println("Can't insert graylog task. ", err, t)
			return false
		}

		t.Id, _ = result.LastInsertId()

		return true
	} else {
		_, err := db.Exec("UPDATE tasks_graylog SET pattern = ?, agg_time = ?, min_count = ?, max_count = ? WHERE id = ?",
			t.Pattern, t.AggTime, t.MinCount, t.MaxCount, t.Id)
		if err != nil {
			log.Println("Can't update graylog task. ", err, t)
			return false
		}

		return true
	}

	return false
}
