package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
    "os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
    "github.com/joho/godotenv"

    "github.com/lawrencedrums/gojira/internal"
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

    router.HandleFunc("/issues", getIssues).Methods("GET")
    router.HandleFunc("/issues", createIssue).Methods("POST")
    router.HandleFunc("/issues/{id}", getIssue).Methods("GET")
    router.HandleFunc("/issues/{id}", updateIssue).Methods("PUT")

    http.ListenAndServe(":8000", router)
}

func getIssues(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    result, err := db.Query("SELECT id, title, body FROM issues WHERE is_archived=0")
    if err != nil {
        panic(err.Error())
    }
    defer result.Close()

    var issues []models.Issue

    for result.Next() {
        var issue models.Issue
        err := result.Scan(&issue.ID, &issue.Title, &issue.Body)
        if err != nil {
            panic(err.Error())
        }

        issues = append(issues, issue)
    }
    json.NewEncoder(w).Encode(issues)
}

func createIssue(w http.ResponseWriter, r *http.Request) {
    stmt, err := db.Prepare("INSERT INTO issues(title, body, is_archived) VALUES(?,?,?)")
    if err != nil {
        panic(err.Error())
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        panic(err.Error())
    }

    keyVal := make(map[string]string)
    json.Unmarshal(body, &keyVal)
    issueTitle := keyVal["title"]
    issueBody := keyVal["body"]
    issueIsArchived := keyVal["is_archived"]

    var res sql.Result
    res, err = stmt.Exec(issueTitle, issueBody, issueIsArchived)
    if err != nil {
        panic(err.Error())
    }

    issueId, err := res.LastInsertId()
    fmt.Fprintf(w, "Issue ID %d created", issueId)
}

func getIssue(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)

    result, err := db.Query("SELECT id, title, body, is_archived FROM issues WHERE id = ?;", params["id"])
    if err != nil {
        panic(err.Error())
    }
    defer result.Close()

    var issue models.Issue

    for result.Next() {
        err := result.Scan(&issue.ID, &issue.Title, &issue.Body, &issue.IsArchived)
        if err != nil {
            panic(err.Error())
        }
    }
    json.NewEncoder(w).Encode(issue)
}

func updateIssue(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)

    stmt, err := db.Prepare("UPDATE issues SET title = ?, body = ?, is_archived = ? WHERE id = ?;")
    if err != nil {
        panic(err.Error())
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        panic(err.Error())
    }

    keyVal := make(map[string]string)
    json.Unmarshal(body, &keyVal)
    newTitle := keyVal["title"]
    newBody := keyVal["body"]
    newIsArchived := keyVal["is_archived"]

    _, err = stmt.Exec(newTitle, newBody, newIsArchived, params["id"])
    if err != nil {
        panic(err.Error())
    }

    fmt.Fprintf(w, "Issus %s updated", params["id"])
}
