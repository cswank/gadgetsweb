package models

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestSaveGadget(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	os.Setenv("GADGETSDB", path.Join(tmp, "db"))
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
	os.RemoveAll(tmp)
}

func TestDeleteGadget(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	os.Setenv("GADGETSDB", path.Join(tmp, "db"))
	db, _ := GetDB()
	defer db.Close()
	g := Gadget{
		Name: "brewery",
		Host: "192.168.1.16",
	}
	g.Save()
	g = Gadget{
		Name: "lab",
		Host: "192.168.1.17",
	}
	g.Save()
	gadgets, _ := GetGadgets()
	if len(gadgets.Gadgets) != 2 {
		t.Fatal(gadgets)
	}
	g = gadgets.Gadgets[0]
	g.Delete()
	gadgets, _ = GetGadgets()
	if len(gadgets.Gadgets) != 1 {
		t.Fatal(gadgets)
	}
	os.RemoveAll(tmp)
}

func TestGetGadgets(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	os.Setenv("GADGETSDB", path.Join(tmp, "db"))
	db, _ := GetDB()
	defer db.Close()
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
	os.RemoveAll(tmp)
}
