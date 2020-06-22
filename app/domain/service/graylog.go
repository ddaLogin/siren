package service

import (
	"fmt"
	"github.com/ddalogin/siren/app/domain/model"
	"github.com/ddalogin/siren/app/domain/repository"
	"github.com/ddalogin/siren/elasticsearch"
	"log"
	"net/url"
	"strconv"
	"time"
)

// Конфиг для грейлог задач
type GraylogConfig struct {
	Es      string // Хост грейлог эластика
	BaseUrl string // Хост грейлога, для построения ссылки
}

// Сервис задач
type GraylogService struct {
	config                   GraylogConfig
	tasksGraylogRepository   *repository.TasksGraylogRepository
	resultsGraylogRepository *repository.ResultsGraylogRepository
}

// Конструктор грейлог сервиса
func NewGraylogService(config GraylogConfig, tasksGraylogRepository *repository.TasksGraylogRepository, resultsGraylogRepository *repository.ResultsGraylogRepository) *GraylogService {
	return &GraylogService{config: config, tasksGraylogRepository: tasksGraylogRepository, resultsGraylogRepository: resultsGraylogRepository}
}

// Выполнить задачу для грейлога
func (s *GraylogService) RunTaskGraylog(task *model.Task) (result model.ResultGraylog) {
	if task.ObjectType() != model.TYPE_GRAYLOG {
		log.Println("Задача не является задачей для грейлога", task)
		return
	}

	taskGraylog := s.tasksGraylogRepository.GetById(task.ObjectId())
	if taskGraylog == nil {
		return
	}

	result.SetTask(task)
	result.SetTaskGraylogId(taskGraylog.Id())
	result.SetTitle(task.Title())
	result.SetGraylogLink(s.buildGraylogUrl(taskGraylog))

	esClient := elasticsearch.NewClient(s.config.Es)
	response, err := esClient.Search(taskGraylog.Pattern(), taskGraylog.AggregateTime())
	if err != nil {
		result.SetStatus(model.STATUS_ERROR)
		result.SetMessage("Техническая ошибка")
		result.SetText(err.Error())
		result.SetCount(0)
	} else {
		result.SetStatus(model.STATUS_OK)
		result.SetCount(response.Hits.Count)
		result.SetMessage("Кол-во сообщений в норме")

		if result.Count() < taskGraylog.Min() {
			result.SetStatus(model.STATUS_ALERT)
			result.SetMessage("Кол-во сообщений уменьшилось")
		}

		if result.Count() > taskGraylog.Max() {
			result.SetStatus(model.STATUS_ALERT)
			result.SetMessage("Кол-во сообщений увеличилось")
		}

		result.SetText(fmt.Sprintf(
			"Ожидаемое кол-во сообщений от %d до %d \r\n Текущее кол-во сообщений %d",
			taskGraylog.Min(), taskGraylog.Max(), result.Count(),
		))
	}

	s.resultsGraylogRepository.Save(&result)

	return
}

// Удалить задачу
func (s *GraylogService) DeleteTask(task *model.Task) bool {

	if !s.resultsGraylogRepository.DeleteByTaskId(task.Id()) {
		return false
	}

	if !s.tasksGraylogRepository.DeleteById(task.Id()) {
		return false
	}

	return true
}

// Собирает ссылку на грейлог
func (s *GraylogService) buildGraylogUrl(task *model.TaskGraylog) string {
	baseUrl, _ := url.Parse(s.config.BaseUrl + "search")
	timeMarker := task.AggregateTime()[len(task.AggregateTime())-1:]
	aggTime, _ := strconv.Atoi(task.AggregateTime()[:len(task.AggregateTime())-1])

	switch timeMarker {
	case `m`:
		aggTime = aggTime * 60
	case `h`:
		aggTime = aggTime * 60 * 60
	}

	query := baseUrl.Query()
	query.Add("q", task.Pattern())
	query.Add("rangetype", "absolute")
	query.Add("from", time.Now().Add((time.Duration(aggTime)*time.Second)*-1).Format("2006-01-02T15:04:05Z07:00"))
	query.Add("to", time.Now().Format("2006-01-02T15:04:05Z07:00"))
	baseUrl.RawQuery = query.Encode()

	return baseUrl.String()
}
