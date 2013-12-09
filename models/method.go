package models

import (
	"encoding/json"
	"database/sql"
)

var (
	getMethodsQuery = "SELECT id, name, steps FROM methods"
	addMethodQuery = "INSERT INTO methods (name, steps) VALUES (?, ?)"
	updateMethodQuery = "UPDATE methods set name = ?, steps = ? WHERE id = ?"
)

type Methods struct {
	Methods []*Method `json:"methods"`
}

type Method struct {
	Id uint64 `json:"id"`
	Name string `json:"name"`
	Steps []string `json:"steps"`
}

func GetMethods() (*Methods, error) {
	db, err := getDB()
	defer db.Close()
	methods := &Methods{}
	rows, err := db.Query(getMethodsQuery)
	if err != nil {
		return methods, err
	}
	for rows.Next() {
		m, err := GetMethod(rows)
		if err != nil {
			return methods, err
		}
		methods.Methods = append(methods.Methods, m)
	}
	return methods, err
}

func GetMethod(rows *sql.Rows) (*Method, error) {
	m := &Method{}
	var methodData []byte
	if err := rows.Scan(&m.Id, &m.Name, &methodData); err != nil {
		return m, err
	}
	var steps []string
	err := json.Unmarshal(methodData, &steps)
	if err != nil {
		return m, err
	}
	m.Steps = steps
	return m, err
}

func (m *Method)Save() error {
	db, err := getDB()
	defer db.Close()
	if err != nil {
		return err
	}
	b, err := json.Marshal(m.Steps)
	if err != nil {
		return err
	}
	steps := string(b)
	if m.Id > 0 {
		_, err = db.Query(updateMethodQuery, m.Name, steps, m.Id)
	} else {
		_, err = db.Query(addMethodQuery, m.Name, steps)
	}
	return err
}
