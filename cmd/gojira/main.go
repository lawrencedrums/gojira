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
    database.DBCon, err = sql.Open("mysql", dbSource)
    if err != nil {
        panic(err.Error())
    }
    defer database.DBCon.Close()

    if err = database.DBCon.Ping(); err != nil {
        panic(err.Error())
    }
    fmt.Println("Connection to DB established")

    ensureDBExists(database.DBCon)
    ensureTablesExists(database.DBCon)

    router := router.NewRouter()
    http.ListenAndServe(":8000", router)
}

func ensureDBExists(DB *sql.DB) {
    _, err = DB.Exec("CREATE DATABASE IF NOT EXISTS gojira")
    if err != nil {
        panic(err.Error())
    }

    _, err = DB.Exec("USE gojira")
    if err != nil {
        panic(err.Error())
    }
}

func ensureTablesExists(DB *sql.DB) {
    var stmt *sql.Stmt
    stmt, err = DB.Prepare(`
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
}
