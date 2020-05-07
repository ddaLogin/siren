package model

import (
	"database/sql"
	"log"
)

// Телеграм чат
type TelegramChat struct {
	id        int    // Идентификатор
	username  string // Имя пользователя
	chatId    string // Идентификатор чата
	createdAt string // Дата создания
}

// Создает модель телеграм чата по строке из базы
func ScanTelegramChat(row *sql.Row) (chat TelegramChat) {
	err := row.Scan(
		&chat.id,
		&chat.username,
		&chat.chatId,
		&chat.createdAt,
	)
	if err != nil {
		log.Println("Не удалось собрать телеграм чата", row)
	}

	return
}

func (t *TelegramChat) Id() int {
	return t.id
}

func (t *TelegramChat) SetId(id int) {
	t.id = id
}

func (t *TelegramChat) Username() string {
	return t.username
}

func (t *TelegramChat) SetUsername(username string) {
	t.username = username
}

func (t *TelegramChat) ChatId() string {
	return t.chatId
}

func (t *TelegramChat) SetChatId(chatId string) {
	t.chatId = chatId
}

func (t *TelegramChat) CreatedAt() string {
	return t.createdAt
}

func (t *TelegramChat) SetCreatedAt(createdAt string) {
	t.createdAt = createdAt
}
