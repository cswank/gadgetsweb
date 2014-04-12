package models

import (
	"testing"
	"path"
	"io/ioutil"
)


func TestSaveMethod(t *testing.T) {
	tmp, _ := ioutil.TempDir("/tmp", "")
	DBPath = path.Join(tmp, "gadgets.db")
	m := Method{
		Name: "flash!",
		Steps: []string{"turn on lab led 2", "wait 2 seconds", "turn off lab led 2"},
	}
	err := m.Save()
	if err != nil {
		t.Error(err)
	}
	db := getDB()
	m2 := db.Methods["flash!"]
	if m2.Name != "flash!" {
		t.Error(m2)
	}
}

// func _TestGetMethods(t *testing.T) {
// 	db, _ := getDB()
// 	defer db.Close()
// 	db.Query("CREATE TABLE methods(id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, steps TEXT)")
// 	db.Query("DELETE FROM methods")
// 	m := Method{
// 		Name: "flash!",
// 		Steps: []string{"turn on lab led 2", "wait 2 seconds", "turn off lab led 2"},
// 	}
// 	err := m.Save()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	methods, err := GetMethods("lab")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if len(methods.Methods) != 1 {
// 		t.Error(methods)
// 	}
// 	m2 := methods.Methods[0]
// 	if m2.Name != "flash!" {
// 		t.Error(m2)
// 	}
// 	if m2.Steps[2] != "turn off lab led 2" {
// 		t.Error(m2)
// 	}
// 	db.Query("DELETE FROM methods")
// }
