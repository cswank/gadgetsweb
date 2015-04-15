package models

import (
	"database/sql"
	"encoding/json"
)

var (
	getMethodsQuery   = "SELECT id, name, steps FROM methods where gadget = ?"
	addMethodQuery    = "INSERT INTO methods (name, gadget, steps) VALUES (?, ?, ?)"
	updateMethodQuery = "UPDATE methods set name = ?, steps = ? WHERE id = ?"
	deleteMethodQuery = "DELETE FROM methods WHERE id = ?"
)

type Methods struct {
	Methods []*Method `json:"methods"`
}

type Method struct {
	Id     uint64   `json:"id"`
	Name   string   `json:"name"`
	Gadget string   `json:"gadget"`
	Steps  []string `json:"steps"`
}

func GetMethods(gadget string) (*Methods, error) {
	methods := &Methods{}
	rows, err := DB.Query(getMethodsQuery, gadget)
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

func (m *Method) Delete() error {
	_, err := DB.Exec(deleteMethodQuery, m.Id)
	return err
}

func (m *Method) Save() error {
	b, err := json.Marshal(m.Steps)
	if err != nil {
		return err
	}
	steps := string(b)
	if m.Id > 0 {
		_, err = DB.Query(updateMethodQuery, m.Name, steps, m.Id)
	} else {
		_, err = DB.Query(addMethodQuery, m.Name, m.Gadget, steps)
	}
	return err
}
