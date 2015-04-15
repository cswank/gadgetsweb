package models

import (
	"database/sql"
	"fmt"
	"os"

	_ "code.google.com/p/go-sqlite/go1/sqlite3"
	"github.com/cswank/gadgetsweb/utils"
)

var DB *sql.DB

func init() {
	p := os.Getenv("GADGETS_DB")
	if p == "" {
		p = ":memory:"
	}
	var err error
	fmt.Println("db is", p)
	DB, err = sql.Open("sqlite3", p)
	if err != nil {
		panic(err)
	}
	if p == ":memory:" || !utils.FileExists(p) {
		createTables()
	}
}

type Timeseries struct {
	Name string        `json:"name"`
	Data []interface{} `json:"data"`
}

type Summary struct {
	Location  string `json:"location"`
	Name      string `json:"name"`
	Direction string `json:"direction"`
}

func createTables() {
	DB.Exec("CREATE TABLE users (username text PRIMARY KEY, password text, permission text)")
	DB.Exec("CREATE TABLE gadgets (name text PRIMARY KEY, host text)")
	DB.Exec("CREATE TABLE methods (id INTEGER PRIMARY KEY AUTOINCREMENT, gadget TEXT, name TEXT, steps TEXT)")
	DB.Exec("CREATE TABLE notes (id INTEGER PRIMARY KEY AUTOINCREMENT, gadget TEXT, text TEXT, taken INTEGER)")
}
