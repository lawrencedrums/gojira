package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lawrencedrums/gojira/api/v1/handlers"
)

func NewRouter() *mux.Router {
    router := mux.NewRouter()

    router.HandleFunc("/", handlers.BaseHandler)
    router.HandleFunc("/issues", handlers.GetIssues).Methods("GET")
    router.HandleFunc("/issues", handlers.CreateIssue).Methods("POST")
    router.HandleFunc("/issues/{id}", handlers.GetIssue).Methods("GET")
    router.HandleFunc("/issues/{id}", handlers.UpdateIssue).Methods("PUT")
    router.HandleFunc("/issues/edit/{id}", handlers.EditIssue).Methods("GET")
    router.HandleFunc("/issues/new/", handlers.NewIssueForm).Methods("GET")
    router.HandleFunc("/reset", handlers.Reset).Methods("GET")

    staticDir := "./cmd/gojira/static"
    router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

    return router
}
