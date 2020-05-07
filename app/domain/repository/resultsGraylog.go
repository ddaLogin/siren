package repository

import (
	"database/sql"
	"fmt"
	"github.com/ddalogin/siren/app/domain/model"
	"log"
)

var resultsGraylogRepository *ResultsGraylogRepository

// Репозиторий для результатов задач грейлога
type ResultsGraylogRepository struct {
	db *sql.DB
}

// Фабричный метод для репозитория результатов задач грейлога
func GetResultsGraylogRepository(db *sql.DB) *ResultsGraylogRepository {
	if resultsGraylogRepository == nil {
		resultsGraylogRepository = &ResultsGraylogRepository{
			db: db,
		}
	}

	return resultsGraylogRepository
}

// Сохранить результат
func (r *ResultsGraylogRepository) Save(resultTask *model.ResultGraylog) bool {
	if resultTask.Id() == 0 {
		result, err := r.db.Exec(
			"INSERT INTO results_graylog (task_graylog_id, title, status, message, text, count, graylog_link) VALUE (?, ?, ?, ?, ?, ?, ?)",
			resultTask.TaskGraylogId(), resultTask.Title(), resultTask.Status(), resultTask.Message(), resultTask.Text(), resultTask.Count(), resultTask.GraylogLink(),
		)
		if err != nil {
			fmt.Printf(err.Error())
			log.Println("Не удалось сохранить результат задачи для грейлога.", err, resultTask)
			return false
		}

		id, err := result.LastInsertId()
		if err != nil {
			log.Println("Не удалось получить ID результата задачи для грейлога.", err, resultTask)
			return false
		}

		resultTask.SetId(int(id))
	} else {
		_, err := r.db.Exec(
			"UPDATE results_graylog SET task_graylog_id = ?, title = ?, status = ?, message = ?, text = ?, count = ?, graylog_link = ? WHERE id = ?",
			resultTask.TaskGraylogId(), resultTask.Title(), resultTask.Status(), resultTask.Message(), resultTask.Text(), resultTask.Count(), resultTask.GraylogLink(), resultTask.Id(),
		)
		if err != nil {
			log.Println("Не удалось обновить результат задачи для грейлога.", err, resultTask)
			return false
		}
	}

	return true
}
