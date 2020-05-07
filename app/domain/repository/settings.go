package repository

import (
	"database/sql"
	"github.com/ddalogin/siren/app/domain/model"
	"log"
)

var settingRepository *SettingRepository

// Репозиторий для настроек
type SettingRepository struct {
	db *sql.DB
}

// Фабричный метод для репозитория настроек
func GetSettingRepository(db *sql.DB) *SettingRepository {
	if settingRepository == nil {
		settingRepository = &SettingRepository{
			db: db,
		}
	}

	return settingRepository
}

// Получить настройку по ID
func (r *SettingRepository) GetById(id int) (setting *model.Setting) {
	row := r.db.QueryRow("SELECT * FROM settings WHERE id = ?", id)
	if row == nil {
		log.Println("Не удалось найти настроку по ID", id)
		return
	}

	return model.ScanSetting(row)
}
