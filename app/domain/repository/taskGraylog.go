package repository

import (
	"database/sql"
	"github.com/ddalogin/siren/app/domain/model"
	"log"
)

var tasksGraylogRepository *TasksGraylogRepository

// Репозиторий для задач грейлога
type TasksGraylogRepository struct {
	db *sql.DB
}

// Фабричный метод для репозитория задач грейлога
func GetTasksGraylogRepository(db *sql.DB) *TasksGraylogRepository {
	if tasksGraylogRepository == nil {
		tasksGraylogRepository = &TasksGraylogRepository{
			db: db,
		}
	}

	return tasksGraylogRepository
}

// Получить задачу по ID
func (r *TasksGraylogRepository) GetById(id int) *model.TaskGraylog {
	row := r.db.QueryRow("SELECT * FROM tasks_graylog WHERE id = ?", id)
	if row == nil {
		log.Println("Не удалось найти задачу для грейлога по ID", id)
		return nil
	}

	task := model.ScanTaskGraylog(row)

	return &task
}

// Сохранить задачу
func (r *TasksGraylogRepository) Save(task *model.TaskGraylog) bool {
	if task.Id() == 0 {
		result, err := r.db.Exec(
			"INSERT INTO tasks_graylog (pattern, aggregate_time, min, max) VALUE (?, ?, ?, ?)",
			task.Pattern(), task.AggregateTime(), task.Min(), task.Max(),
		)
		if err != nil {
			log.Println("Не удалось сохранить задачу для грейлога.", err, task)
			return false
		}

		id, err := result.LastInsertId()
		if err != nil {
			log.Println("Не удалось получить ID новой задачи для грейлога.", err, task)
			return false
		}

		task.SetId(int(id))
	} else {
		_, err := r.db.Exec(
			"UPDATE tasks_graylog SET pattern = ?, aggregate_time = ?, min = ?, max = ? WHERE id = ?",
			task.Pattern(), task.AggregateTime(), task.Min(), task.Max(), task.Id(),
		)
		if err != nil {
			log.Println("Не удалось обновить задачу для грейлога.", err, task)
			return false
		}
	}

	return true
}
