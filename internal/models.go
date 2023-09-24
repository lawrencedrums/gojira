package models

type Issue struct {
    ID         string `json:"id"`
    Title      string `json:"title"`
    Body       string `json:"body"`
    IsArchived bool   `json:"IsArchived"`
}
