package router

import (
    "net/http"

    "github.com/gorilla/mux"

    "github.com/lawrencedrums/gojira/api/v1/handlers"
)

func Router() {
    router := mux.NewRouter()

    router.HandleFunc("/issues", handlers.GetIssues).Methods("GET")
    router.HandleFunc("/issues", handlers.CreateIssue).Methods("POST")
    router.HandleFunc("/issues/{id}", handlers.GetIssue).Methods("GET")
    router.HandleFunc("/issues/{id}", handlers.UpdateIssue).Methods("PUT")

    http.ListenAndServe(":8000", router)
}
