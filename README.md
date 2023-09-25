# Gojira

A minimalist, "JIRA" alternative written in Go & HTMX (Gojira is also one of my favourite bands!)

## Requirements
[Go](https://go.dev/dl/)

[mysql](https://dev.mysql.com/downloads/mysql/)

## Setup
Rename `.env.template` to `.env`

Insert your mysql username(usually `root`) and password that you setup during installation

## Running
`go run ./cmd/gojira` to start server
or
`go build ./cmd/gojira && ./gojira` to build and run the executable

## Interacting
Navigate to `localhost:8000` to see all issues created

You can use `curl` or tools like Postman to interact with the DB at `localhost:8000`

Available actions:
1. GET `/issues/` - return all issues created
2. POST `/issues/` - create new issues
3. GET `/issues/{id}` - return issue with the given id
4. PUT `/issues/{id}` - update issue with the given id

