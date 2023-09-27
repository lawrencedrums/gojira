package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
    "html/template"
	"io"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

    "github.com/lawrencedrums/gojira/internal/models"
    "github.com/lawrencedrums/gojira/internal/database"
)

var tplDir = "cmd/gojira/templates"

func BaseHandler(w http.ResponseWriter, r *http.Request) {
    baseTpl := fmt.Sprintf("%s/base.html", tplDir)
    t := template.Must(template.ParseFiles(baseTpl))
    t.Execute(w, nil)
}

func GetIssues(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    result, err := database.DBCon.Query("SELECT id, title, body FROM issues WHERE is_archived=0")
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
    issuesTpl := fmt.Sprintf("%s/issue_board.html", tplDir)
    t := template.Must(template.ParseFiles(issuesTpl))
    t.Execute(w, issues)
}

func CreateIssue(w http.ResponseWriter, r *http.Request) {
    stmt, err := database.DBCon.Prepare("INSERT INTO issues(title, body, is_archived) VALUES(?,?,?)")
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

func GetIssue(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)

    result, err := database.DBCon.Query("SELECT id, title, body, is_archived FROM issues WHERE id = ?;", params["id"])
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

    issueDetailsTpl := fmt.Sprintf("%s/issue_details.html", tplDir)
    t := template.Must(template.ParseFiles(issueDetailsTpl))
    t.Execute(w, issue)
}

func EditIssue(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)

    result, err := database.DBCon.Query("SELECT id, title, body, is_archived FROM issues WHERE id = ?;", params["id"])
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

    editIssueTpl := fmt.Sprintf("%s/issue_edit.html", tplDir)
    t := template.Must(template.ParseFiles(editIssueTpl))
    t.Execute(w, issue)
}

func UpdateIssue(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)

    stmt, err := database.DBCon.Prepare("UPDATE issues SET title = ?, body = ?, is_archived = ? WHERE id = ?;")
    if err != nil {
        panic(err.Error())
    }

    r.ParseForm()
    newTitle := r.Form["title"][0]
    newBody := r.Form["body"][0]
    newIsArchived := "0"

    _, err = stmt.Exec(newTitle, newBody, newIsArchived, params["id"])
    if err != nil {
        panic(err.Error())
    }

    fmt.Printf("Issus %s updated", params["id"])

    issue := models.Issue{
        ID: params["id"],
        Title: newTitle,
        Body: newBody,
        IsArchived: true,
    }

    baseTpl := fmt.Sprintf("%s/base.html", tplDir)
    t := template.Must(template.ParseFiles(baseTpl))
    t.Execute(w, issue)
}

func Reset(w http.ResponseWriter, r *http.Request) {
    database.DBCon.Exec("TRUNCATE TABLE issues")
}
