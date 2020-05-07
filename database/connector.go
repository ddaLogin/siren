package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// Конфиг подключения к БД
type Config struct {
	User string
	Pass string
	Name string
}

// Конструктор коннектора к базе данных
func NewConnector(config Config) *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@/%v", config.User, config.Pass, config.Name))
	if err != nil {
		log.Println("Не удалось подключиться к базе данных", err)
		return nil
	}

	return db
}
