# Gojira

A minimalist, "JIRA" alternative written in Go (Gojira is also one of my favourite bands!)

## Requirements
[Go](https://go.dev/dl/)
[mysql](https://dev.mysql.com/downloads/mysql/)

## Setup
Rename `.env.template` to `.env`
Insert your mysql username(usually `root`) and password that you setup during installation

## Running
`go run .` to start server
or
`go build && ./gojira` to build and run the server

## Interacting
Right now it is only a simple CRUD endpoint
You can use `curl` or tools like Postman to interact with the DB at `localhost:8000`

Available actions:
1. GET `/issues/` - return all issues created
2. POST `/issues/` - create new issues
3. GET `/issues/{id}` - return issue with the given id
4. PUT `/issues/{id}` - update issue with the given id

