package http

import (
	"github.com/ddalogin/siren/container"
	"github.com/ddalogin/siren/http/controllers"
	"log"
	"net/http"
)

// Настройки для сервера
type Config struct {
	Host string
	Port string
}

// Сервер
type Server struct {
	config    Config
	container *container.Container
}

// Конструктор сервера
func NewServer(config Config, container *container.Container) *Server {
	return &Server{config: config, container: container}
}

// Запуск http сервера для веб морды
func (s *Server) Run() {
	s.initRoutes()

	fs := http.FileServer(http.Dir("./http/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	err := http.ListenAndServe(`:`+s.config.Port, nil)

	if err != nil {
		log.Fatal("Не удалось запустить сервер", err)
	}
}

// Регистрация роутов
func (s *Server) initRoutes() {

	indexController := controllers.IndexController{Container: s.container}
	http.HandleFunc("/", indexController.IndexAction)

	taskController := controllers.TaskController{Container: s.container}
	http.HandleFunc("/task", taskController.FormAction)
	http.HandleFunc("/task/run", taskController.RunAction)
	http.HandleFunc("/task/delete", taskController.DeleteAction)
	http.HandleFunc("/task/result/list", taskController.ResultListAction)

	taskGraylogController := controllers.TaskGraylogController{Container: s.container}
	http.HandleFunc("/task/graylog/result", taskGraylogController.ResultAction)
	http.HandleFunc("/task/graylog/result/list", taskGraylogController.ResultListAction)
}
