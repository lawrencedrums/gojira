package main

import (
	"database/sql"
	"fmt"
    "net/http"
    "os"

	_ "github.com/go-sql-driver/mysql"
    "github.com/joho/godotenv"

    "github.com/lawrencedrums/gojira/api/v1/router"
    "github.com/lawrencedrums/gojira/internal/database"
)

var (
    DBUser string
    DBPass string
    err error
)

func init() {
    err = godotenv.Load(".env")
    if err != nil {
        panic(err.Error())
    }
    DBUser = os.Getenv("DBUser")
    DBPass = os.Getenv("DBPass")
}

func main() {
    dbSource := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/", DBUser, DBPass)
    database.DBCon, err = sql.Open("mysql", dbSource)
    if err != nil {
        panic(err.Error())
    }
    defer database.DBCon.Close()

    if err = database.DBCon.Ping(); err != nil {
        panic(err.Error())
    }

    ensureDBExists()
    ensureTablesExists()

    fmt.Println("Server running at localhost:8000")
    router := router.NewRouter()
    http.ListenAndServe(":8000", router)
}

func ensureDBExists() {
    _, err = database.DBCon.Exec("CREATE DATABASE IF NOT EXISTS gojira")
    if err != nil {
        panic(err.Error())
    }

    _, err = database.DBCon.Exec("USE gojira")
    if err != nil {
        panic(err.Error())
    }
}

func ensureTablesExists() {
    var stmt *sql.Stmt
    stmt, err =database.DBCon.Prepare(`
        CREATE TABLE IF NOT EXISTS projects(
            id INT NOT NULL AUTO_INCREMENT,
            title VARCHAR(255),
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

    stmt, err = database.DBCon.Prepare(`
        CREATE TABLE IF NOT EXISTS issues(
            id INT NOT NULL AUTO_INCREMENT,
            project_id INT,
            title VARCHAR(255),
            body VARCHAR(1020),
            is_archived BOOLEAN,
            PRIMARY KEY (id),
            FOREIGN KEY (project_id) REFERENCES projects(id)
        );
    `)
    if err != nil {
        panic(err.Error())
    }

    _, err = stmt.Exec()
    if err != nil {
        panic(err.Error())
    }
}
