package models

import (
    "database/sql"
    "log"

    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
    var err error
    DB, err = sql.Open("sqlite3", "./db/plants.db")
    if err != nil {
        log.Fatal(err)
    }

    createTables()
}

func createTables() {
    _, err := DB.Exec(`CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL
    );`)
    if err != nil {
        log.Fatal(err)
    }

    _, err = DB.Exec(`CREATE TABLE IF NOT EXISTS plants (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER,
        plant_name TEXT NOT NULL,
        trefle_id INTEGER NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id)
    );`)
    if err != nil {
        log.Fatal(err)
    }

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS roots (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`)
    if err != nil {
        log.Fatal(err)
    }

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		root_id INTEGER,
		user_id INTEGER,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (root_id) REFERENCES roots(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`)
    if err != nil {
        log.Fatal(err)
    }
}