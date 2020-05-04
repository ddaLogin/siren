package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ddalogin/siren/database"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
func (t TaskGraylog) Do() TaskResult {
	result := TaskResult{}
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

	searchResult := ElasticResponse{}

	err = json.Unmarshal(body, &searchResult)
	if err != nil {
		log.Println("Не удалось распарсить ответ от ElasticSearch. ", err, string(body))
		result.Status = STATUS_ERROR
		result.Error = err
		return result
	}

	if t.MinCount != nil && searchResult.Hits.Count < *t.MinCount {
		result.Status = STATUS_ALERT
		result.Message = "Кол-во сообщей уменьшилось"
		result.Body = fmt.Sprintf("Шаблон поиска: %s\r\nТекущее кол-во сообщений: %s",
			t.Pattern,
			strconv.Itoa(searchResult.Hits.Count),
		)
	}

	if t.MaxCount != nil && searchResult.Hits.Count > *t.MaxCount {
		result.Status = STATUS_ALERT
		result.Message = "Кол-во сообщей увеличилось"
		result.Body = fmt.Sprintf("Шаблон поиска: %s\r\nТекущее кол-во сообщений: %s",
			t.Pattern,
			strconv.Itoa(searchResult.Hits.Count),
		)
	}

	result.Info = fmt.Sprintf("Шаблон поиска: %s\r\nТекущее кол-во сообщений: %s\r\nМаксимальное кол-во сообщений: %s\r\nМинимальное кол-во сообщений: %s",
		t.Pattern,
		strconv.Itoa(searchResult.Hits.Count),
		strconv.Itoa(*t.MaxCount),
		strconv.Itoa(*t.MinCount),
	)

	return result
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
