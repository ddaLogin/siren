package service

import (
	"fmt"
	"github.com/ddalogin/siren/app/domain/model"
	"github.com/ddalogin/siren/app/domain/repository"
	"strconv"
	"time"
)

var Host string

// Сервис для репортов
type ReportService struct {
	taskRepository           *repository.TasksRepository
	taskGraylogRepository    *repository.TasksGraylogRepository
	resultsGraylogRepository *repository.ResultsGraylogRepository
}

// Конструктор сервиса репортов
func NewReportService(taskRepository *repository.TasksRepository, taskGraylogRepository *repository.TasksGraylogRepository, resultsGraylogRepository *repository.ResultsGraylogRepository) *ReportService {
	return &ReportService{taskRepository: taskRepository, taskGraylogRepository: taskGraylogRepository, resultsGraylogRepository: resultsGraylogRepository}
}

// Построенние репорта за последние сутки
func (s *ReportService) MakeReport() *model.TelegramNotify {
	yesterdayReport := time.Now().AddDate(0, 0, -1).Format("2006-01-02 15:04:05")
	text := "*Ежедневный отчет от " + yesterdayReport + "*\r\n"
	tasks := s.taskRepository.GetForReport()

	if len(tasks) == 0 {
		return nil
	}

	for _, task := range tasks {
		message := s.buildReportForTask(task, yesterdayReport)

		if message != "" {
			text = text + "\r\n" + message
		}
	}

	return model.NewTelegramNotify(text, nil)
}

// Собрать репорт по задаче
func (s *ReportService) buildReportForTask(task *model.Task, date string) string {
	text := ""

	if task.ObjectType() == model.TYPE_GRAYLOG {
		graylogReports := s.resultsGraylogRepository.GetReportsForDate(task.ObjectId(), date)

		text = text + fmt.Sprintf("[%s](%s)\r\n", task.Title(), Host+"/task/graylog/result/list?id="+strconv.Itoa(task.ObjectId()))

		for _, graylogReport := range graylogReports {
			text = text + fmt.Sprintf("%s - %d \r\n", graylogReport.Message(), graylogReport.Count())
		}
	}

	return text
}
