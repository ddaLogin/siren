package model

import (
	"database/sql"
	"log"
)

// Модель репорта грейлога
type ReportGraylog struct {
	message string
	status  string
	count   int
}

// Создает массив репортов по задаче
func ScanReportsGraylog(rows *sql.Rows) (results []*ReportGraylog) {
	for rows.Next() {
		result := ReportGraylog{}
		err := rows.Scan(
			&result.message,
			&result.status,
			&result.count,
		)
		if err != nil {
			log.Println("Не удалось собрать модель отчета задачи для грейлога", err)
			return
		}

		results = append(results, &result)
	}

	return
}

func (r *ReportGraylog) Message() string {
	return r.message
}

func (r *ReportGraylog) SetMessage(message string) {
	r.message = message
}

func (r *ReportGraylog) Status() string {
	return r.status
}

func (r *ReportGraylog) SetStatus(status string) {
	r.status = status
}

func (r *ReportGraylog) Count() int {
	return r.count
}

func (r *ReportGraylog) SetCount(count int) {
	r.count = count
}
