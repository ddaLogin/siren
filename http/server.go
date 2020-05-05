package http

import (
	"github.com/ddalogin/siren/http/controllers"
	"net/http"
)

type Config struct {
	Port string
}

// Запуск http сервера для веб морды
func StartServer(config Config) {
	initRoutes()

	fs := http.FileServer(http.Dir("./http/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(`:`+config.Port, nil)
}

func initRoutes() {
	http.HandleFunc("/", controllers.IndexAction)
	http.HandleFunc("/task", controllers.FormAction)
	http.HandleFunc("/task/run", controllers.RunAction)
	http.HandleFunc("/task/delete", controllers.DeleteAction)
	http.HandleFunc("/task/result", controllers.ResultAction)
	http.HandleFunc("/task/result/list", controllers.ResultListAction)
}
