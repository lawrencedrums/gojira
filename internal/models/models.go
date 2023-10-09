package models


type Project struct {
    ID         string
    Title      string
    IsArchived bool
}

type Issue struct {
    ID         string
    Project    Project
    Title      string
    Body       string
    IsArchived bool
}
