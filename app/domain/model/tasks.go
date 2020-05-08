package model

import (
	"database/sql"
	"log"
	"time"
)

// Тип задач для GrayLog'а
const TYPE_GRAYLOG = "graylog"

// Модель задачи
type Task struct {
	id         int     // Идентификатор задачи
	title      string  // Заголовок задачи
	objectType string  // Тип задачи
	objectId   int     // Идентификатор деталей задачи
	interval   int     // Кол-во минут между запусками
	nextTime   string  // Время следующего запуска по интервалу
	isEnabled  bool    // Активная ли задача
	usernames  *string // Список юзернеймов Telegram'а разделеных запятой, если пустой то отправляем в дефолтный чат
}

// Создает модель задачи по строке из базы
func ScanTask(row *sql.Row) (task Task) {
	err := row.Scan(
		&task.id,
		&task.title,
		&task.objectType,
		&task.objectId,
		&task.interval,
		&task.nextTime,
		&task.isEnabled,
		&task.usernames,
	)
	if err != nil {
		log.Println("Не удалось собрать модель задачи", row)
	}

	return
}

// Создает массив моделей задач по строкам из базы
func ScanTasks(rows *sql.Rows) (tasks []*Task) {
	for rows.Next() {
		task := Task{}
		err := rows.Scan(
			&task.id,
			&task.title,
			&task.objectType,
			&task.objectId,
			&task.interval,
			&task.nextTime,
			&task.isEnabled,
			&task.usernames,
		)
		if err != nil {
			log.Println("Не удалось собрать модель задачи из массив строк", err)
			return
		}

		tasks = append(tasks, &task)
	}

	return
}

func (t *Task) Id() int {
	return t.id
}

func (t *Task) SetId(id int) {
	t.id = id
}

func (t *Task) Title() string {
	return t.title
}

func (t *Task) SetTitle(title string) {
	t.title = title
}

func (t *Task) ObjectType() string {
	return t.objectType
}

func (t *Task) SetObjectType(objectType string) {
	t.objectType = objectType
}

func (t *Task) ObjectId() int {
	return t.objectId
}

func (t *Task) SetObjectId(objectId int) {
	t.objectId = objectId
}

func (t *Task) Interval() int {
	return t.interval
}

func (t *Task) SetInterval(interval int) {
	t.interval = interval
}

func (t *Task) NextTime() string {
	return t.nextTime
}

func (t *Task) SetNextTime(nextTime string) {
	t.nextTime = nextTime
}

func (t *Task) IsEnabled() bool {
	return t.isEnabled
}

func (t *Task) SetIsEnabled(isEnabled bool) {
	t.isEnabled = isEnabled
}

func (t *Task) Usernames() *string {
	return t.usernames
}

func (t *Task) SetUsernames(usernames *string) {
	t.usernames = usernames
}

// Высчитывает время для следующего запуска
func (t *Task) CalculateNextTime() string {
	return time.Now().Add(time.Minute * time.Duration(t.interval)).Format("2006-01-02 15:04:05")
}
