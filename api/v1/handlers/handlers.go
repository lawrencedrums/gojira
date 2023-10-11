package handlers

import (
	"database/sql"
	"fmt"
    "html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

    "github.com/lawrencedrums/gojira/internal/models"
    "github.com/lawrencedrums/gojira/internal/database"
)

var tplDir = "cmd/gojira/templates"

func BaseHandler(w http.ResponseWriter, r *http.Request) {
    allTpl := fmt.Sprintf("%s/*.html", tplDir)
    t := template.Must(template.ParseGlob(allTpl))
    t.Execute(w, nil)
}

func GetIssues(w http.ResponseWriter, r *http.Request) {
    projects := models.GetProjects()

    if len(projects) == 0 {
        issuesTpl := fmt.Sprintf("%s/project_new.html", tplDir)
        t := template.Must(template.ParseFiles(issuesTpl))
        t.Execute(w, nil)
        return
    }

    projectID := projects[0].ID
    issuesRes, err := database.DBCon.Query(`
        SELECT i.id, i.title, i.body
        FROM issues i
        INNER JOIN projects p
        ON i.id = p.id
        WHERE i.is_archived=0 AND p.id=?;`,
        projectID,
    )
    if err != nil {
        panic(err.Error())
    }
    defer issuesRes.Close()

    var issues []models.Issue

    for issuesRes.Next() {
        var issue models.Issue
        err := issuesRes.Scan(&issue.ID, &issue.Title, &issue.Body)
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

    r.ParseForm()
    issueTitle := r.Form["title"][0]
    issueBody := r.Form["body"][0]
    issueIsArchived := "0"

    var res sql.Result
    res, err = stmt.Exec(issueTitle, issueBody, issueIsArchived)
    if err != nil {
        panic(err.Error())
    }

    issueId, err := res.LastInsertId()
    fmt.Printf("Issue ID %d created\n", issueId)

    indexTpl := fmt.Sprintf("%s/index.html", tplDir)
    t := template.Must(template.ParseFiles(indexTpl))
    t.Execute(w, nil)
}

func NewIssueForm(w http.ResponseWriter, r *http.Request) {
    newIssueTpl := fmt.Sprintf("%s/issue_new.html", tplDir)
    t := template.Must(template.ParseFiles(newIssueTpl))
    t.Execute(w, nil)
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
    if val, ok := r.Form["isArchived"]; ok {
        newIsArchived = val[0]
    }

    _, err = stmt.Exec(newTitle, newBody, newIsArchived, params["id"])
    if err != nil {
        panic(err.Error())
    }

    issue := models.Issue{
        ID: params["id"],
        Title: newTitle,
        Body: newBody,
        IsArchived: newIsArchived != "0",
    }

    fmt.Printf("Issus %s updated\n", params["id"])

    indexTpl := fmt.Sprintf("%s/index.html", tplDir)
    t := template.Must(template.ParseFiles(indexTpl))
    t.Execute(w, issue)
}

func CreateProject(w http.ResponseWriter, r *http.Request) {
    stmt, err := database.DBCon.Prepare("INSERT INTO projects(title, is_archived) VALUES(?,?)")
    if err != nil {
        panic(err.Error())
    }

    r.ParseForm()
    projectTitle := r.Form["title"][0]

    var res sql.Result
    res, err = stmt.Exec(projectTitle, "0")
    if err != nil {
        panic(err.Error())
    }

    projectId, err := res.LastInsertId()
    fmt.Printf("Project ID %d created\n", projectId)

    indexTpl := fmt.Sprintf("%s/index.html", tplDir)
    t := template.Must(template.ParseFiles(indexTpl))
    t.Execute(w, nil)
}

func Reset(w http.ResponseWriter, r *http.Request) {
    database.DBCon.Exec("TRUNCATE TABLE issues")
    database.DBCon.Exec("TRUNCATE TABLE projects")
}
