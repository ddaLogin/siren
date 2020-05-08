package model

import (
	"fmt"
	"strings"
)

const STATUS_OK = 0      // Все в порядке
const STATUS_ALERT = 1   // Что то не так, требуется уведомление
const STATUS_ERROR = 100 // Технический сбой во время задачи

// Результат грейлог задачи
type ResultGraylog struct {
	id            int    // Идентификатор результата
	taskGraylogId int    // Идентификатор грейлог задачи
	status        int    // Статус результат
	title         string // Заголовок задачи
	message       string // Сообщение
	text          string // Полное описание
	count         int    // Кол-во сообщений
	graylogLink   string // Ссылка на грейлог
	createdAt     string // Дата результатов
	task          *Task  // Задача
}

func (r *ResultGraylog) Id() int {
	return r.id
}

func (r *ResultGraylog) SetId(id int) {
	r.id = id
}

func (r *ResultGraylog) TaskGraylogId() int {
	return r.taskGraylogId
}

func (r *ResultGraylog) SetTaskGraylogId(taskGraylogId int) {
	r.taskGraylogId = taskGraylogId
}

func (r *ResultGraylog) Status() int {
	return r.status
}

func (r *ResultGraylog) SetStatus(status int) {
	r.status = status
}

func (r *ResultGraylog) Title() string {
	return r.title
}

func (r *ResultGraylog) SetTitle(title string) {
	r.title = title
}

func (r *ResultGraylog) Message() string {
	return r.message
}

func (r *ResultGraylog) SetMessage(message string) {
	r.message = message
}

func (r *ResultGraylog) Text() string {
	return r.text
}

func (r *ResultGraylog) SetText(text string) {
	r.text = text
}

func (r *ResultGraylog) Count() int {
	return r.count
}

func (r *ResultGraylog) SetCount(count int) {
	r.count = count
}

func (r *ResultGraylog) GraylogLink() string {
	return r.graylogLink
}

func (r *ResultGraylog) SetGraylogLink(graylogLink string) {
	r.graylogLink = graylogLink
}

func (r *ResultGraylog) CreatedAt() string {
	return r.createdAt
}

func (r *ResultGraylog) SetCreatedAt(createdAt string) {
	r.createdAt = createdAt
}

func (r *ResultGraylog) Task() *Task {
	return r.task
}

func (r *ResultGraylog) SetTask(task *Task) {
	r.task = task
}

// Требуется ли уведомление
func (r *ResultGraylog) IsNeedNotify() bool {
	return r.Status() == STATUS_ALERT || r.Status() == STATUS_ERROR
}

// Собрать сообщение для уведомления в телеграме
func (r *ResultGraylog) BuildTelegramNotify() *TelegramNotify {
	var userNames []string
	text := r.Text()
	userNamesString := r.Task().Usernames()

	if userNamesString != nil {
		userNames = strings.Split(*userNamesString, ",")
	}

	if r.Status() == STATUS_ERROR {
		text = "```" + text + "```"
	}

	message := fmt.Sprintf("*%s* \r\n %s \r\n %s \r\n [Смотреть сообщения](%s)",
		r.Title(), r.Message(), text, r.GraylogLink(),
	)

	return NewTelegramNotify(message, userNames)
}
