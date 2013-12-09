package models

import (
	"time"
	"os"
	"database/sql"
	_ "code.google.com/p/go-sqlite/go1/sqlite3"
	"bitbucket.com/cswank/gadgetsweb/utils"
)

type Device struct {
	Units string      `json:"units"`
	Value interface{} `json:"value"`
	ID    string      `json:"id"`
}

type Location struct {
	Input  map[string]Device `json:"input"`
	Output map[string]Device `json:"output"`
}

type Gadgets struct {
	Sender      string      `json:"sender"`
	Timestamp   time.Time   `json:"timestamp"`
	Name        string      `json:"name"`
	Locations   map[string]Location    `json:"locations"`
}

type Timeseries struct {
	Name     string              `json:"name"`
	Data     []interface{}       `json:"data"`
}

type Summary struct {
	Location string `json:"location"`
	Name string `json:"name"`
	Direction string `json:"direction"`
}

func createTables(db *sql.DB) {
	db.Query("CREATE TABLE users(username text PRIMARY KEY, password text)")
	db.Query("CREATE TABLE gadgets(name text PRIMARY KEY, host text)")
	db.Query("CREATE TABLE methods(id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, steps TEXT)")
}

func getDB() (*sql.DB, error) {
	p := os.Getenv("GADGETSDB")
	db, err := sql.Open("sqlite3", p)
	if err != nil {
		return db, err
	}
	if p != "" && !utils.FileExists(p) {
		createTables(db)
	}
	return db, err
}
