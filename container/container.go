package container

import (
	"database/sql"
	"github.com/ddalogin/siren/app/domain/repository"
	"github.com/ddalogin/siren/app/domain/service"
)

// Контейнер зависимостей
type Container struct {
	graylogConfig service.GraylogConfig
	notifyConfig  service.NotifyConfig

	connection *sql.DB

	taskRepository           *repository.TasksRepository
	taskGraylogRepository    *repository.TasksGraylogRepository
	resultsGraylogRepository *repository.ResultsGraylogRepository
	telegramChatRepository   *repository.TelegramChatRepository

	graylogService *service.GraylogService
	taskService    *service.TaskService
	notifyService  *service.NotifyService
	reportService  *service.ReportService
}

// Конструктор контейнера
func NewContainer(graylogConfig service.GraylogConfig, notifyConfig service.NotifyConfig, connection *sql.DB) *Container {
	return &Container{graylogConfig: graylogConfig, notifyConfig: notifyConfig, connection: connection}
}

func (c *Container) Connection() *sql.DB {
	return c.connection
}

func (c *Container) SetConnection(connection *sql.DB) {
	c.connection = connection
}

func (c *Container) TaskRepository() *repository.TasksRepository {
	if c.taskRepository == nil {
		c.SetTaskRepository(repository.GetTasksRepository(c.connection))
	}

	return c.taskRepository
}

func (c *Container) SetTaskRepository(taskRepository *repository.TasksRepository) {
	c.taskRepository = taskRepository
}

func (c *Container) TaskGraylogRepository() *repository.TasksGraylogRepository {
	if c.taskGraylogRepository == nil {
		c.SetTaskGraylogRepository(repository.GetTasksGraylogRepository(c.connection))
	}

	return c.taskGraylogRepository
}

func (c *Container) SetTaskGraylogRepository(taskGraylogRepository *repository.TasksGraylogRepository) {
	c.taskGraylogRepository = taskGraylogRepository
}

func (c *Container) ResultsGraylogRepository() *repository.ResultsGraylogRepository {
	if c.resultsGraylogRepository == nil {
		c.SetResultsGraylogRepository(repository.GetResultsGraylogRepository(c.connection))
	}

	return c.resultsGraylogRepository
}

func (c *Container) SetResultsGraylogRepository(resultsGraylogRepository *repository.ResultsGraylogRepository) {
	c.resultsGraylogRepository = resultsGraylogRepository
}

func (c *Container) TelegramChatRepository() *repository.TelegramChatRepository {
	if c.telegramChatRepository == nil {
		c.SetTelegramChatRepository(repository.GetTelegramChatRepository(c.connection))
	}

	return c.telegramChatRepository
}

func (c *Container) SetTelegramChatRepository(telegramChatRepository *repository.TelegramChatRepository) {
	c.telegramChatRepository = telegramChatRepository
}

func (c *Container) GraylogService() *service.GraylogService {
	if c.graylogService == nil {
		c.SetGraylogService(
			service.NewGraylogService(c.graylogConfig, c.TaskGraylogRepository(), c.ResultsGraylogRepository()),
		)
	}

	return c.graylogService
}

func (c *Container) SetGraylogService(graylogService *service.GraylogService) {
	c.graylogService = graylogService
}

func (c *Container) TaskService() *service.TaskService {
	if c.taskService == nil {
		c.SetTaskService(
			service.NewTaskService(c.TaskRepository(), c.GraylogService()),
		)
	}

	return c.taskService
}

func (c *Container) SetTaskService(taskService *service.TaskService) {
	c.taskService = taskService
}

func (c *Container) NotifyService() *service.NotifyService {
	if c.notifyService == nil {
		c.SetNotifyService(
			service.NewNotifyService(c.notifyConfig, c.TelegramChatRepository()),
		)
	}

	return c.notifyService
}

func (c *Container) SetNotifyService(notifyService *service.NotifyService) {
	c.notifyService = notifyService
}

func (c *Container) ReportService() *service.ReportService {
	if c.reportService == nil {
		c.SetReportService(
			service.NewReportService(c.taskRepository, c.taskGraylogRepository, c.resultsGraylogRepository),
		)
	}

	return c.reportService
}

func (c *Container) SetReportService(reportService *service.ReportService) {
	c.reportService = reportService
}
