package models

import (
	"testing"
)


func TestSaveGadget(t *testing.T) {
	db, _ := getDB()
	defer db.Close()
	db.Query("CREATE TABLE gadgets(name text PRIMARY KEY, host text)")
	db.Query("DELETE FROM gadgets")
	g := Gadget{
		Name: "brewery",
		Host: "192.168.1.16",
	}
	err := g.Save()
	if err != nil {
		t.Error(err)
	}
	g2 := Gadget{}
	err = db.QueryRow("SELECT * from gadgets WHERE name = ?", "brewery").Scan(&g2.Name, &g2.Host)
	if err != nil {
		t.Error(err)
	}
	if g2.Name != "brewery" {
		t.Error(g2)
	}
	if g2.Host != "192.168.1.16" {
		t.Error(g2)
	}
	db.Query("DELETE FROM gadgets")
}


func TestGetGadgets(t *testing.T) {
	db, _ := getDB()
	defer db.Close()
	db.Query("CREATE TABLE gadgets(name text PRIMARY KEY, host text)")
	db.Query("DELETE FROM gadgets")
	g := Gadget{
		Name: "brewery",
		Host: "192.168.1.16",
	}
	g.Save()
	g = Gadget{
		Name: "greenhouse",
		Host: "192.168.1.13",
	}
	g.Save()
	gadgets, err := GetGadgets()
	if err != nil {
		t.Error(err)
	}
	if len(gadgets.Gadgets) != 2 {
		t.Error(gadgets)
	}
	g1 := gadgets.Gadgets[0]
	if g1.Name != "brewery" {
		t.Error(gadgets)
	}
	db.Query("DELETE FROM gadgets")
}
