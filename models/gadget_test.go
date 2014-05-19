package models

import (
	"testing"
)


func TestToMap(t *testing.T) {
	g := Gadget{
 		Name: "brewery",
 		Host: "192.168.1.16",
 	}
	m := g.toMap()
	if m["Name"] != "brewery" {
		t.Error(m["Name"])
	}
}

// func _TestSaveGadget(t *testing.T) {
// 	DBPath = "/tmp/gadgetsdb"
// 	os.RemoveAll(DBPath)
// 	defer os.RemoveAll(DBPath)

// 	if err := myDB.Create("gadgets", 2); err != nil {
// 		t.Fatal(err)
// 	}
	
// 	g := Gadget{
// 		Name: "brewery",
// 		Host: "192.168.1.16",
// 	}
// 	err := g.Save()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	// db, _ := getDB()
// 	// if db.Gadgets["brewery"].Host != "192.168.1.16" {
// 	// 	t.Error(db)
// 	// }
// }



// func _TestGetGadgets(t *testing.T) {
// 	g := Gadget{
// 		Name: "brewery",
// 		Host: "192.168.1.16",
// 	}
// 	g.Save()
// 	g = Gadget{
// 		Name: "greenhouse",
// 		Host: "192.168.1.13",
// 	}
// 	g.Save()
// 	gadgets := GetGadgets()
// 	if len(gadgets) != 2 {
// 		t.Error(gadgets)
// 	}
// 	g1 := gadgets[0]
// 	if g1.Name != "brewery" {
// 		t.Error(gadgets)
// 	}
// }
