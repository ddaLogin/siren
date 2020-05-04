package worker

import "fmt"

const STATUS_SILENT = 0  // Все хорошо, сообщать не о чем
const STATUS_WARNING = 1 // Рано паниковать, но стоит обратить внимание
const STATUS_ALERT = 2   // Пора паниковать
const STATUS_ERROR = 100 // Технический сбой в ходе выполнения задачи

// Модель результатов выполнения задачи
type TaskResult struct {
	Status  int    // Статус выполнения задачи
	Error   error  // Ошибка в случае статуса = 100
	Task    Task   // Задача
	Message string // Подзаголовок
	Body    string // Тело сообщения
	Info    string // Не участвует в уведомлениях, отображается при ручном запуске через веб интерфейс
}

// Получить сообщение для телеграмма
func (t TaskResult) GetMessage() string {
	if t.Status == STATUS_ERROR {
		return fmt.Sprintf("*%s*\r\n%s\r\n%s", t.Task.Title, "Технический сбой во время выполнения", t.Error)
	} else {
		return fmt.Sprintf("*%s*\r\n%s\r\n%s", t.Task.Title, t.Message, t.Body)
	}
}
