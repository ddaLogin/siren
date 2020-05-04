package worker

import (
	"database/sql"
	"fmt"
	"github.com/ddalogin/siren/database"
	"log"
)

const STATUS_SILENT = 0  // Все хорошо, сообщать не о чем
const STATUS_WARNING = 1 // Рано паниковать, но стоит обратить внимание
const STATUS_ALERT = 2   // Пора паниковать
const STATUS_ERROR = 100 // Технический сбой в ходе выполнения задачи

// Модель результатов выполнения задачи
type TaskResult struct {
	Id        int64  // Идентификатор
	TaskId    int64  // Идентификатор задачи
	Task      Task   // Задача
	Status    int    // Статус выполнения задачи
	Message   string // Подзаголовок
	Body      string // Тело сообщения
	Info      string // Не участвует в уведомлениях, отображается при ручном запуске через веб интерфейс
	Error     string // Ошибка в случае статуса = 100
	CreatedAt string // Дата выполнения
}

// Получить сообщение для телеграмма
func (res TaskResult) GetMessage() string {
	if res.Status == STATUS_ERROR {
		return fmt.Sprintf("*%s*\r\n%s\r\n%s", res.Task.Title, "Технический сбой во время выполнения", res.Error)
	} else {
		return fmt.Sprintf("*%s*\r\n%s\r\n%s", res.Task.Title, res.Message, res.Body)
	}
}

// Получить все рузльтаты выполнения по ID задачи
func GetResultByTaskId(taskId int, limit int) (results []TaskResult) {
	db := database.Db()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM results WHERE task_id = ? ORDER BY created_at DESC LIMIT ?", taskId, limit)
	if err != nil {
		log.Println("Не удалось найти результаты задачи. ", err)
		return
	}
	defer rows.Close()

	results, err = scanResults(rows)
	if err != nil {
		log.Println("При поиске результатов задачи, не удалось собрать модель. ", err)
	}

	return
}

// Получить 1 результат по ID
func GetResultById(id int) (result TaskResult) {
	db := database.Db()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM results WHERE id = ?", id)
	if err != nil {
		log.Println("Не удалось найти результат. ", err)
		return
	}
	defer rows.Close()

	results, err := scanResults(rows)
	if err != nil {
		log.Println("При поиске результата, не удалось собрать модель. ", err)
	}

	if len(results) >= 0 {
		result = results[0]
	}

	return
}

// Сохранить результат выполнения задачи
func (res *TaskResult) Save() bool {
	db := database.Db()
	defer db.Close()

	if res.Id == 0 {
		result, err := db.Exec("INSERT INTO results (status, task_id, message, body, info, error) VALUES (?, ?, ?, ?, ?, ?)",
			res.Status, res.TaskId, res.Message, res.Body, res.Info, res.Error)
		if err != nil {
			log.Println("Can't insert result. ", err, res)
			return false
		}

		res.Id, _ = result.LastInsertId()

		result, err = db.Exec("DELETE FROM results WHERE task_id = ? AND created_at < UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 14 DAY))", res.TaskId)
		if err != nil {
			log.Println("Не удалось удалить старые записи. ", err, res)
			return false
		}

		return true
	} else {
		_, err := db.Exec("UPDATE results SET status = ?, task_id = ?, message = ?, body = ?, info = ?, error = ? WHERE id = ?",
			res.Status, res.TaskId, res.Message, res.Body, res.Info, res.Error, res.Id)
		if err != nil {
			log.Println("Can't update result. ", err, res)
			return false
		}

		return true
	}

	return false
}

// Парсит маасив результатов в модели
func scanResults(rows *sql.Rows) (results []TaskResult, err error) {
	for rows.Next() {
		result := TaskResult{}
		err = rows.Scan(
			&result.Id,
			&result.Status,
			&result.TaskId,
			&result.Message,
			&result.Body,
			&result.Info,
			&result.Error,
			&result.CreatedAt,
		)
		if err != nil {
			return
		}

		results = append(results, result)
	}

	return
}
