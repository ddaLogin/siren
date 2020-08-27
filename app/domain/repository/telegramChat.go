package repository

import (
	"github.com/ddalogin/siren/app/domain/model"
	"github.com/ddalogin/siren/database"
	"log"
)

var telegramChatRepository *TelegramChatRepository

// Репозиторий для телеграм чатов
type TelegramChatRepository struct{}

// Фабричный метод для репозитория результатов задач грейлога
func GetTelegramChatRepository() *TelegramChatRepository {
	if telegramChatRepository == nil {
		telegramChatRepository = &TelegramChatRepository{}
	}

	return telegramChatRepository
}

// Получить чат по юзернейму
func (r *TelegramChatRepository) GetByUserName(username string) *model.TelegramChat {
	db := database.Db()
	defer db.Close()

	row := db.QueryRow("SELECT * FROM telegram_chats WHERE username = ?", username)
	if row == nil {
		log.Println("Не удалось найти чат по юзернейму", username)
		return nil
	}

	chat := model.ScanTelegramChat(row)

	return &chat
}

// Сохранить телеграм чат
func (r *TelegramChatRepository) Save(chat *model.TelegramChat) bool {
	db := database.Db()
	defer db.Close()

	if chat.Id() == 0 {
		result, err := db.Exec(
			"INSERT INTO telegram_chats (username, chat_id) VALUE (?, ?)",
			chat.Username(), chat.ChatId(),
		)
		if err != nil {
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
		_, err := db.Exec(
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
