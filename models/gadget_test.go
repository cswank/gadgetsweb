package models

import (
	"testing"
)



func TestSaveGadget(t *testing.T) {
	DBPath = "/tmp/gadgets.db"
	g := Gadget{
		Name: "brewery",
		Host: "192.168.1.16",
	}
	err := g.Save()
	if err != nil {
		t.Error(err)
	}

	db := getDB()
	if db.Gadgets["brewery"].Host != "192.168.1.16" {
		t.Error(db)
	}
}


func TestGetGadgets(t *testing.T) {
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
	gadgets := GetGadgets()
	if len(gadgets) != 2 {
		t.Error(gadgets)
	}
	g1 := gadgets[0]
	if g1.Name != "brewery" {
		t.Error(gadgets)
	}
}
