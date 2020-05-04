package worker

import (
	"database/sql"
	"fmt"
	"github.com/ddalogin/siren/database"
	"log"
)

const INTERVAL_5M = 1
const INTERVAL_15M = 2
const INTERVAL_30M = 3
const INTERVAL_60M = 4

// Общая модель задачи
type Task struct {
	Id         int64  // Идентификатор
	Title      string // Заголовок задачи
	IsEnable   bool   // Включена ли задача
	ObjectType int    // Тип задачи
	ObjectId   int64  // Объект с параметрами задачи по типу
	Interval   int    // Интервал циклического выполнения
	Time       string // Точное время для выполнения
}

// Интерфейс задачи
type TaskInterface interface {
	Do() TaskResult // Выполнение задачу
}

// Выполнение задачу
func (task Task) Do() TaskResult {
	fmt.Println("Выполнение задачи: " + task.Title)

	graylogTask := GetTaskGraylogById(task.ObjectId)
	result := graylogTask.Do()
	result.Task = task

	fmt.Println("Выполнение задачи: " + task.Title + " завершилось")

	return result
}

// Получить все задачи по времени
func GetAllTasksByTime(time string) (tasks []Task) {
	db := database.Db()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM tasks WHERE time = ? AND is_enable = true", time)
	if err != nil {
		log.Println("Не удалось найти задачи по времени. ", err)
		return
	}
	defer rows.Close()

	tasks, err = scanArray(rows)
	if err != nil {
		log.Println("При поиске по времени, не удалось собрать модель задачи. ", err)
	}

	return
}

// Получить все задачи по интервалу
func GetAllTasksByInterval(minutes int) (tasks []Task) {

	mInterval := 0

	switch minutes {
	case 0:
		mInterval = INTERVAL_60M
		break
	case 30:
		mInterval = INTERVAL_30M
		break
	case 45:
		fallthrough
	case 15:
		mInterval = INTERVAL_15M
		break
	case 55:
		fallthrough
	case 50:
		fallthrough
	case 40:
		fallthrough
	case 35:
		fallthrough
	case 25:
		fallthrough
	case 20:
		fallthrough
	case 10:
		fallthrough
	case 5:
		mInterval = INTERVAL_5M
		break
	}

	if mInterval == 0 {
		return
	}

	db := database.Db()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM tasks WHERE `interval` > 0 AND `interval` <= ? AND is_enable = true", mInterval)
	if err != nil {
		log.Println("Не удалось найти задачи по интервалу. ", err)
		return
	}
	defer rows.Close()

	tasks, err = scanArray(rows)
	if err != nil {
		log.Println("При поиске по интервалу, не удалось собрать модель задачи. ", err)
	}

	return
}

// Получить все задачи
func GetAllTasks() (tasks []Task) {
	db := database.Db()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM tasks")
	if err != nil {
		log.Println("Не удалось найти все задачи. ", err)
		return
	}
	defer rows.Close()

	tasks, err = scanArray(rows)
	if err != nil {
		log.Println("При поиске всех задач, не удалось собрать модель задачи. ", err)
	}

	return
}

// Получить 1 задачу
func GetTaskById(id int) (task Task) {
	db := database.Db()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM tasks WHERE id = ?", id)
	if err != nil {
		log.Println("Не удалось найти задачу. ", err)
		return
	}
	defer rows.Close()

	tasks, err := scanArray(rows)
	if err != nil {
		log.Println("При поиске задачи, не удалось собрать модель. ", err)
	}

	if len(tasks) >= 0 {
		task = tasks[0]
	}

	return
}

// Сохранить Task
func (task *Task) Save() bool {
	db := database.Db()
	defer db.Close()

	if task.Id == 0 {
		result, err := db.Exec("INSERT INTO tasks (title, is_enable, object_type, object_id, `interval`, time) VALUES (?, ?, ?, ?, ?, ?)",
			task.Title, task.IsEnable, task.ObjectType, task.ObjectId, task.Interval, task.Time)
		if err != nil {
			log.Println("Can't insert task. ", err, task)
			return false
		}

		task.Id, _ = result.LastInsertId()

		return true
	} else {
		_, err := db.Exec("UPDATE tasks SET title = ?, is_enable = ?, object_type = ?, object_id = ?, `interval` = ?, time = ? WHERE id = ?",
			task.Title, task.IsEnable, task.ObjectType, task.ObjectId, task.Interval, task.Time, task.Id)
		if err != nil {
			log.Println("Can't update task. ", err, task)
			return false
		}

		return true
	}

	return false
}

// Парсит маасив задач в модели
func scanArray(rows *sql.Rows) (tasks []Task, err error) {
	for rows.Next() {
		task := Task{}
		err = rows.Scan(
			&task.Id,
			&task.Title,
			&task.IsEnable,
			&task.ObjectType,
			&task.ObjectId,
			&task.Interval,
			&task.Time,
		)
		if err != nil {
			return
		}

		tasks = append(tasks, task)
	}

	return
}
