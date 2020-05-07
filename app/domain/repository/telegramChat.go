package repository

import (
	"database/sql"
	"fmt"
	"github.com/ddalogin/siren/app/domain/model"
	"log"
)

var telegramChatRepository *TelegramChatRepository

// Репозиторий для телеграм чатов
type TelegramChatRepository struct {
	db *sql.DB
}

// Фабричный метод для репозитория результатов задач грейлога
func GetTelegramChatRepository(db *sql.DB) *TelegramChatRepository {
	if telegramChatRepository == nil {
		telegramChatRepository = &TelegramChatRepository{
			db: db,
		}
	}

	return telegramChatRepository
}

// Получить чат по юзернейму
func (r *TelegramChatRepository) GetByUserName(username string) *model.TelegramChat {
	row := r.db.QueryRow("SELECT * FROM telegram_chats WHERE username = ?", username)
	if row == nil {
		log.Println("Не удалось найти чат по юзернейму", username)
		return nil
	}

	chat := model.ScanTelegramChat(row)

	return &chat
}

// Сохранить телеграм чат
func (r *TelegramChatRepository) Save(chat *model.TelegramChat) bool {
	if chat.Id() == 0 {
		result, err := r.db.Exec(
			"INSERT INTO telegram_chats (username, chat_id) VALUE (?, ?)",
			chat.Username(), chat.ChatId(),
		)
		if err != nil {
			fmt.Printf(err.Error())
			log.Println("Не удалось сохранить телеграм чат.", err, chat)
			return false
		}

		id, err := result.LastInsertId()
		if err != nil {
			log.Println("Не удалось получить ID телеграм чата.", err, chat)
			return false
		}

		chat.SetId(int(id))
	} else {
		_, err := r.db.Exec(
			"UPDATE telegram_chats SET username = ?, chat_id = ?, created_at = ? WHERE id = ?",
			chat.Username(), chat.ChatId(), chat.CreatedAt(), chat.Id(),
		)
		if err != nil {
			log.Println("Не удалось обновить телеграм чат.", err, chat)
			return false
		}
	}

	return true
}
