package interfaces

import (
	"github.com/ddalogin/siren/app/domain/model"
)

// Интерфейс результатов задачи
type TaskResult interface {
	IsNeedNotify() bool                         // Требуется ли уведомление
	BuildTelegramNotify() *model.TelegramNotify // Собрать сообщение для уведомления в телеграме
}
