package repository

import (
	"database/sql"
	"github.com/ddalogin/siren/app/domain/model"
	"log"
	"time"
)

var tasksRepository *TasksRepository

// Репозиторий для задач
type TasksRepository struct {
	db *sql.DB
}

// Фабричный метод для репозитория задач
func GetTasksRepository(db *sql.DB) *TasksRepository {
	if tasksRepository == nil {
		tasksRepository = &TasksRepository{
			db: db,
		}
	}

	return tasksRepository
}

// Получить задачу по ID
func (r *TasksRepository) GetById(id int) *model.Task {
	row := r.db.QueryRow("SELECT * FROM tasks WHERE id = ?", id)
	if row == nil {
		log.Println("Не удалось найти задачу по ID", id)
		return nil
	}

	return model.ScanTask(row)
}

// Получить задачи подходящих по времени к запуску
func (r *TasksRepository) GetForRun(time time.Time) []*model.Task {
	rows, err := r.db.Query("SELECT * FROM tasks WHERE enabled = 1 AND next_time <= ?", time.Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Println("Не удалось найти задачи для запуска", time)
		return nil
	}
	defer rows.Close()

	return model.ScanTasks(rows)
}

// Сохранить задачу
func (r *TasksRepository) Save(task *model.Task) bool {
	if task.Id() == 0 {
		result, err := r.db.Exec(
			"INSERT INTO tasks (title, object_type, object_id, `interval`, next_time, enabled, usernames) VALUE (?, ?, ?, ?, ?, ?, ?)",
			task.Title(), task.ObjectType(), task.ObjectId(), task.Interval(), task.NextTime(), task.IsEnabled(), task.Usernames(),
		)
		if err != nil {
			log.Println("Не удалось сохранить задачу.", err, task)
			return false
		}

		id, err := result.LastInsertId()
		if err != nil {
			log.Println("Не удалось получить ID новой задачи.", err, task)
			return false
		}

		task.SetId(int(id))
	} else {
		_, err := r.db.Exec(
			"UPDATE tasks SET title = ?, object_type = ?, object_id = ?, `interval` = ?, next_time = ?, enabled = ?, usernames = ? WHERE id = ?",
			task.Title(), task.ObjectType(), task.ObjectId(), task.Interval(), task.NextTime(), task.IsEnabled(), task.Usernames(), task.Id(),
		)
		if err != nil {
			log.Println("Не удалось обновить задачу.", err, task)
			return false
		}
	}

	return true
}
