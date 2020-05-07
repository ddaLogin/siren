package model

import (
	"database/sql"
	"log"
)

// Задача для поиска ошибок в грейлоге
type TaskGraylog struct {
	id            int    // Идентификатор
	pattern       string // Текст для поиска
	aggregateTime string // Период для аггрегации сообщений
	min           int    // Минимально допустимое кол-во
	max           int    // Максимально допустимое кол-во
}

// Создает модель задачи для грейлога по строке из базы
func ScanTaskGraylog(row *sql.Row) (task TaskGraylog) {
	err := row.Scan(
		&task.id,
		&task.pattern,
		&task.aggregateTime,
		&task.min,
		&task.max,
	)
	if err != nil {
		log.Println("Не удалось собрать модель задачи для грейлога", row)
	}

	return
}

func (t *TaskGraylog) Id() int {
	return t.id
}

func (t *TaskGraylog) SetId(id int) {
	t.id = id
}

func (t *TaskGraylog) Pattern() string {
	return t.pattern
}

func (t *TaskGraylog) SetPattern(pattern string) {
	t.pattern = pattern
}

func (t *TaskGraylog) AggregateTime() string {
	return t.aggregateTime
}

func (t *TaskGraylog) SetAggregateTime(aggregateTime string) {
	t.aggregateTime = aggregateTime
}

func (t *TaskGraylog) Min() int {
	return t.min
}

func (t *TaskGraylog) SetMin(min int) {
	t.min = min
}

func (t *TaskGraylog) Max() int {
	return t.max
}

func (t *TaskGraylog) SetMax(max int) {
	t.max = max
}
