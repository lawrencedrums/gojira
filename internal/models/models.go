package models


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
