package service

import (
	"fmt"
	"github.com/ddalogin/siren/app/domain/model"
	"github.com/ddalogin/siren/app/domain/repository"
	"github.com/ddalogin/siren/app/interfaces"
)

// Сервис задач
type TaskService struct {
	taskRepository *repository.TasksRepository
	graylogService *GraylogService
}

// Конструктор сервиса задач
func NewTaskService(taskRepository *repository.TasksRepository, graylogService *GraylogService) *TaskService {
	return &TaskService{taskRepository: taskRepository, graylogService: graylogService}
}

// Выполнить задачу
func (s *TaskService) RunTask(task *model.Task) (result interfaces.TaskResult) {
	fmt.Println("Выполнение задачи " + task.Title())

	if task.ObjectType() == model.TYPE_GRAYLOG {
		resultGraylog := s.graylogService.RunTaskGraylog(task)
		result = &resultGraylog
	}

	fmt.Println("Выполнение задачи " + task.Title() + " завершилось")

	return result
}

// Удалить задачу
func (s *TaskService) DeleteTask(task *model.Task) bool {

	if task.ObjectType() == model.TYPE_GRAYLOG {
		resultGraylog := s.graylogService.RunTaskGraylog(task)
		result = &resultGraylog
	}

	return result
}
