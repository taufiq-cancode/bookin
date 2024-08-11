package db

import (
    "database/sql"
    "log"
    _ "github.com/glebarez/go-sqlite" // Using the pure Go SQLite driver
)

var DB *sql.DB

func InitDB() {
    var err error
    DB, err = sql.Open("sqlite", "api.db")
    if err != nil {
        log.Fatalf("Could not connect to DB: %v", err)
    }

    DB.SetMaxOpenConns(10)
    DB.SetMaxIdleConns(5)

    err = createTables()
    if err != nil {
        log.Fatalf("Could not create tables: %v", err)
    }
}

func createTables() error {
    createUsersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        email TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL
    )
    `
    _, err := DB.Exec(createUsersTable)
    if err != nil {
        log.Printf("SQL Error: %v", err)
        return err
    }

    createEventsTable := `
    CREATE TABLE IF NOT EXISTS events (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        description TEXT NOT NULL, 
        location TEXT NOT NULL,
        dateTime DATETIME NOT NULL,
        user_id INTEGER,
        FOREIGN KEY(user_id) REFERENCES users(id)
    )
    `
    _, err = DB.Exec(createEventsTable)
    if err != nil {
        log.Printf("SQL Error: %v", err)
        return err
    }

    return nil
}
