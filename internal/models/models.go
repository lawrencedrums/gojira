package models

import "github.com/lawrencedrums/gojira/internal/database"


type Project struct {
    ID         string
    Title      string
    IsArchived bool
}

type Issue struct {
    ID         string
    ProjectID  string
    Title      string
    Body       string
    IsArchived bool
}

func GetProjects() []Project {
    projectsRes, err := database.DBCon.Query("SELECT id, title FROM projects WHERE is_archived=0")
    if err != nil {
        panic(err.Error())
    }
    defer projectsRes.Close()

    var projects []Project

    for projectsRes.Next() {
        var project Project
        err := projectsRes.Scan(&project.ID, &project.Title)
        if err != nil {
            panic(err.Error())
        }

        projects = append(projects, project)
    }
    return projects
}
