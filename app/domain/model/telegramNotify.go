package model

// Уведомления для телеграма
type TelegramNotify struct {
	message   string
	userNames []string
}

// Конструктор
func NewTelegramNotify(message string, userNames []string) *TelegramNotify {
	return &TelegramNotify{message: message, userNames: userNames}
}

func (t *TelegramNotify) Message() string {
	return t.message
}

func (t *TelegramNotify) SetMessage(message string) {
	t.message = message
}

func (t *TelegramNotify) UserNames() []string {
	return t.userNames
}

func (t *TelegramNotify) SetUserNames(userNames []string) {
	t.userNames = userNames
}
