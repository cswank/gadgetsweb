package models

import (
        "os"
        "database/sql"
        _ "code.google.com/p/go-sqlite/go1/sqlite3"
        "bitbucket.org/cswank/gadgetsweb/utils"
)

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
        db.Exec("CREATE TABLE users (username text PRIMARY KEY, password text, permission text)")
        db.Exec("CREATE TABLE gadgets (name text PRIMARY KEY, host text)")
        db.Exec("CREATE TABLE methods (id INTEGER PRIMARY KEY AUTOINCREMENT, gadget TEXT, name TEXT, steps TEXT)")
	db.Exec("CREATE TABLE notes (id INTEGER PRIMARY KEY AUTOINCREMENT, gadget TEXT, text TEXT, taken INTEGER)")
}

func GetDB() (*sql.DB, error) {
        p := os.Getenv("GADGETSDB")
        if p == "" {
                p = ":memory:"
        }
        db, err := sql.Open("sqlite3", p)
        if err != nil {
                return db, err
        }
        if p == ":memory:" || !utils.FileExists(p) {
                createTables(db)
        }
        return db, err
}

