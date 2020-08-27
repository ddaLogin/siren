package repository

import (
	"github.com/ddalogin/siren/app/domain/model"
	"github.com/ddalogin/siren/database"
	"log"
)

var resultsGraylogRepository *ResultsGraylogRepository

// Репозиторий для результатов задач грейлога
type ResultsGraylogRepository struct{}

// Фабричный метод для репозитория результатов задач грейлога
func GetResultsGraylogRepository() *ResultsGraylogRepository {
	if resultsGraylogRepository == nil {
		resultsGraylogRepository = &ResultsGraylogRepository{}
	}

	return resultsGraylogRepository
}

// Получить результат по ID
func (r *ResultsGraylogRepository) GetById(id int) *model.ResultGraylog {
	db := database.Db()
	defer db.Close()

	row := db.QueryRow("SELECT * FROM results_graylog WHERE id = ?", id)
	if row == nil {
		log.Println("Не удалось найти результат задачи для грейлога по ID", id)
		return nil
	}

	result := model.ScanResultGraylog(row)

	return &result
}

// Получить результаты по ID грейлог задачи
func (r *ResultsGraylogRepository) GetByTaskGraylogId(id int, limit int) []*model.ResultGraylog {
	db := database.Db()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM results_graylog WHERE task_graylog_id = ? ORDER BY created_at DESC LIMIT ?", id, limit)
	if err != nil {
		log.Println("Не удалось найти результаты задачи для грейлога по ID", id)
		return nil
	}
	defer rows.Close()

	results := model.ScanResultsGraylog(rows)

	return results
}

// Получить результаты после определенной даты
func (r *ResultsGraylogRepository) GetReportsForDate(id int, date string) []*model.ReportGraylog {
	db := database.Db()
	defer db.Close()

	rows, err := db.Query("SELECT message, status, count(*) FROM results_graylog WHERE task_graylog_id = ? AND created_at > ? GROUP BY status, message", id, date)
	if err != nil {
		log.Println("Не удалось найти результаты позднее даты", date)
		return nil
	}
	defer rows.Close()

	results := model.ScanReportsGraylog(rows)

	return results
}

// Сохранить результат
func (r *ResultsGraylogRepository) Save(resultTask *model.ResultGraylog) bool {
	db := database.Db()
	defer db.Close()

	if resultTask.Id() == 0 {
		result, err := db.Exec(
			"INSERT INTO results_graylog (task_graylog_id, title, status, message, text, count, graylog_link) VALUE (?, ?, ?, ?, ?, ?, ?)",
			resultTask.TaskGraylogId(), resultTask.Title(), resultTask.Status(), resultTask.Message(), resultTask.Text(), resultTask.Count(), resultTask.GraylogLink(),
		)
		if err != nil {
			log.Println("Не удалось сохранить результат задачи для грейлога.", err, resultTask)
			return false
		}

		id, err := result.LastInsertId()
		if err != nil {
			log.Println("Не удалось получить ID результата задачи для грейлога.", err, resultTask)
			return false
		}

		resultTask.SetId(int(id))

		r.DeleteOldestByTaskId(resultTask.Id())
	} else {
		_, err := db.Exec(
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

// Удалить результаты задачи по Id
func (r *ResultsGraylogRepository) DeleteByTaskId(id int) bool {
	db := database.Db()
	defer db.Close()

	result, err := db.Exec("DELETE FROM results_graylog WHERE task_graylog_id = ?", id)
	if err != nil {
		log.Println("Не удалось удалить результаты задачи по её ID")
		return false
	}

	count, _ := result.RowsAffected()

	return count > 0
}

// Удалить слишком старые результаты задачи
func (r *ResultsGraylogRepository) DeleteOldestByTaskId(id int) bool {
	db := database.Db()
	defer db.Close()

	result, err := db.Exec("DELETE FROM results_graylog WHERE task_graylog_id = ? AND created_at < UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 14 DAY))", id)
	if err != nil {
		log.Println("Не удалось удалить старые результаты задачи по её ID")
		return false
	}

	count, _ := result.RowsAffected()

	return count > 0
}
