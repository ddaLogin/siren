package service

import (
	"github.com/ddalogin/siren/app/domain/model"
	"github.com/ddalogin/siren/app/domain/repository"
	"github.com/ddalogin/siren/telegram"
)

// Настройка уведомлений
type NotifyConfig struct {
	Token  string
	Host   string
	ChatId string
	Proxy  string
}

// Сервис уведомления
type NotifyService struct {
	config                 NotifyConfig
	telegramChatRepository *repository.TelegramChatRepository
}

// Конструктор сервиса уведомлений
func NewNotifyService(config NotifyConfig, telegramChatRepository *repository.TelegramChatRepository) *NotifyService {
	return &NotifyService{config: config, telegramChatRepository: telegramChatRepository}
}

// Отправить уведомление в телеграм
func (s *NotifyService) NotifyTelegram(n *model.TelegramNotify) {
	var chatIds []string
	client := telegram.NewClient(s.config.Host, s.config.Token, s.config.Proxy)

	if len(n.UserNames()) == 0 {
		chatIds = append(chatIds, s.config.ChatId)
	} else {
		for _, username := range n.UserNames() {
			chat := s.telegramChatRepository.GetByUserName(username)

			if chat == nil || (model.TelegramChat{}) == *chat {
				chat = s.UpdateTelegramChats(username)
			}

			if chat != nil && (model.TelegramChat{}) != *chat {
				chatIds = append(chatIds, chat.ChatId())
			}
		}
	}

	for _, chatId := range chatIds {
		client.SendMessage(n.Message(), chatId)
	}
}

// Забирает поледние обновления от бота, и регистрирует чаты в базе данных
func (s *NotifyService) UpdateTelegramChats(searchedUsername string) (searchedChat *model.TelegramChat) {
	client := telegram.NewClient(s.config.Host, s.config.Token, s.config.Proxy)
	updates := client.UpdateChats()

	for _, update := range updates.Result {
		if update.Message.Chat.Type == "private" && update.Message.Chat.Username != "" {
			chat := model.TelegramChat{}
			chat.SetUsername(update.Message.Chat.Username)
			chat.SetChatId(update.Message.Chat.Id.String())

			s.telegramChatRepository.Save(&chat)

			if chat.Username() == searchedUsername {
				searchedChat = &chat
			}
		}
	}

	return
}
