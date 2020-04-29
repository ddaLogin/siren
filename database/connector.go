package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// Конфиг подключения к БД
type Config struct {
	User string
	Pass string
	Name string
}

var config Config

// Инициализация базы данных
func InitDatabase(cfg Config) {
	config = cfg
}

// Получить подключение к бд
func Db() (db *sql.DB) {
	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@/%v", config.User, config.Pass, config.Name))

	if err != nil {
		panic(err)
	}

	return
}
