package repository

import (
	"github.com/ddalogin/siren/app/domain/model"
	"github.com/ddalogin/siren/database"
	"log"
	"time"
)

var tasksRepository *TasksRepository

// Репозиторий для задач
type TasksRepository struct{}

// Фабричный метод для репозитория задач
func GetTasksRepository() *TasksRepository {
	if tasksRepository == nil {
		tasksRepository = &TasksRepository{}
	}

	return tasksRepository
}

// Получить задачу по ID
func (r *TasksRepository) GetById(id int) *model.Task {
	db := database.Db()
	defer db.Close()

	row := db.QueryRow("SELECT * FROM tasks WHERE id = ?", id)
	if row == nil {
		log.Println("Не удалось найти задачу по ID", id)
		return nil
	}

	task := model.ScanTask(row)

	return &task
}

// Получить задачи подходящих по времени к запуску
func (r *TasksRepository) GetForRun(time time.Time) []*model.Task {
	db := database.Db()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM tasks WHERE enabled = 1 AND next_time <= ?", time.Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Println("Не удалось найти задачи для запуска", time)
		return nil
	}
	defer rows.Close()

	return model.ScanTasks(rows)
}

// Получить все задачи
func (r *TasksRepository) GetAll() []*model.Task {
	db := database.Db()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM tasks")
	if err != nil {
		log.Println("Не удалось найти все задачи")
		return nil
	}
	defer rows.Close()

	return model.ScanTasks(rows)
}

// Получить задачи для отчета
func (r *TasksRepository) GetForReport() []*model.Task {
	db := database.Db()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM tasks WHERE enabled = 1 AND usernames IS NULL")
	if err != nil {
		log.Println("Не удалось найти задачи для отчета")
		return nil
	}
	defer rows.Close()

	return model.ScanTasks(rows)
}

// Удалить задачу по Id
func (r *TasksRepository) DeleteById(id int) bool {
	db := database.Db()
	defer db.Close()

	result, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		log.Println("Не удалось удалить задачу по ID")
		return false
	}

	count, _ := result.RowsAffected()

	return count > 0
}

// Сохранить задачу
func (r *TasksRepository) Save(task *model.Task) bool {
	db := database.Db()
	defer db.Close()

	if task.Id() == 0 {
		result, err := db.Exec(
			"INSERT INTO tasks (title, object_type, object_id, `interval`, enabled, usernames) VALUE (?, ?, ?, ?, ?, ?)",
			task.Title(), task.ObjectType(), task.ObjectId(), task.Interval(), task.IsEnabled(), task.Usernames(),
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
		_, err := db.Exec(
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
