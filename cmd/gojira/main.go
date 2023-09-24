package main

import (
	"database/sql"
	"fmt"
	"net/http"
    "os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
    "github.com/joho/godotenv"

    "github.com/lawrencedrums/gojira/api/v1"
)

var db *sql.DB
var err error

const (
    dbHost = "127.0.0.1"
    dbPort = "3306"
)

func main() {
    err = godotenv.Load(".env")
    if err != nil {
        panic(err.Error())
    }
    dbUser := os.Getenv("DBUser")
    dbPass := os.Getenv("DBPass")

    dbSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/", dbUser, dbPass, dbHost, dbPort)
    db, err = sql.Open("mysql", dbSource)
    if err != nil {
        panic(err.Error())
    }
    defer db.Close()

    if err = db.Ping(); err != nil {
        panic(err.Error())
    }

    _, err = db.Exec("CREATE DATABASE IF NOT EXISTS gojira")
    if err != nil {
        panic(err.Error())
    }

    _, err = db.Exec("USE gojira")
    if err != nil {
        panic(err.Error())
    }

    fmt.Println("Connection to DB established")

    var stmt *sql.Stmt
    stmt, err = db.Prepare(`
        CREATE TABLE IF NOT EXISTS issues(
            id INT NOT NULL AUTO_INCREMENT,
            title VARCHAR(255),
            body VARCHAR(1020),
            is_archived BOOLEAN,
            PRIMARY KEY (id)
        );
    `)
    if err != nil {
        panic(err.Error())
    }

    _, err = stmt.Exec()
    if err != nil {
        panic(err.Error())
    }

    router := mux.NewRouter()

    router.HandleFunc("/issues", handlers.GetIssues).Methods("GET")
    router.HandleFunc("/issues", handlers.CreateIssue).Methods("POST")
    router.HandleFunc("/issues/{id}", handlers.GetIssue).Methods("GET")
    router.HandleFunc("/issues/{id}", handlers.UpdateIssue).Methods("PUT")

    http.ListenAndServe(":8000", router)
}

