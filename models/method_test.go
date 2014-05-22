package models

import (
	"testing"
	"os"
	"io/ioutil"
	"path"
)


func TestSaveMethod(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	os.Setenv("GADGETSDB", path.Join(tmp, "db"))
	db, _ := GetDB()
	defer db.Close()
	m := Method{
		Name: "flash!",
		Steps: []string{"turn on lab led 2", "wait 2 seconds", "turn off lab led 2"},
	}
	err := m.Save()
	if err != nil {
		t.Fatal(err)
	}
	m2 := Method{}
	var stepsStr string
	err = db.QueryRow("SELECT id, name, steps from methods WHERE name = ?", "flash!").Scan(&m2.Id, &m2.Name, &stepsStr)
	if err != nil {
		t.Fatal(err)
	}
	if m2.Name != "flash!" {
		t.Error(m2)
	}
	os.RemoveAll(tmp)
}

func TestGetMethods(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	os.Setenv("GADGETSDB", path.Join(tmp, "db"))
	db, _ := GetDB()
	defer db.Close()
	m := Method{
		Name: "flash!",
		Gadget: "lab",
		Steps: []string{"turn on lab led 2", "wait 2 seconds", "turn off lab led 2"},
	}
	err := m.Save()
	if err != nil {
		t.Fatal(err)
	}
	methods, err := GetMethods("lab")
	if err != nil {
		t.Error(err)
	}
	if len(methods.Methods) != 1 {
		t.Error(methods)
	}
	m2 := methods.Methods[0]
	if m2.Name != "flash!" {
		t.Error(m2)
	}
	if m2.Steps[2] != "turn off lab led 2" {
		t.Error(m2)
	}
	os.RemoveAll(tmp)
}
