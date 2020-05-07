package model

import (
	"database/sql"
	"log"
)

// Модель настройки
type Setting struct {
	id    int    // Идентификатор настройки
	key   string // Ключ настройки
	value string // Значение настройки
}

// Создает новую модель настройки
func NewSetting(id int, key string, value string) *Setting {
	return &Setting{
		id:    id,
		key:   key,
		value: value,
	}
}

// Создает модель настройки по строке из базы
func ScanSetting(row *sql.Row) (setting *Setting) {
	err := row.Scan(setting)
	if err != nil {
		log.Println("Не удалось собрать модель настроки", row)
	}

	return
}

func (s *Setting) Id() int {
	return s.id
}

func (s *Setting) SetId(id int) {
	s.id = id
}

func (s *Setting) Key() string {
	return s.key
}

func (s *Setting) SetKey(key string) {
	s.key = key
}

func (s *Setting) Value() string {
	return s.value
}

func (s *Setting) SetValue(value string) {
	s.value = value
}
