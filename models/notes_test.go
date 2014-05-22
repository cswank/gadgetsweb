package models

import (
	"os"
	"testing"
	"io/ioutil"
	"path"
	"time"
)


func TestSaveNote(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	os.Setenv("GADGETSDB", path.Join(tmp, "db"))
	db, err := GetDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	n := Note{
		Text: "me do this",
		Gadget: "lab",
	}
	err = n.Save()
	if err != nil {
		t.Fatal(err)
	}

	n2 := &Note{}
	var ts int64
	err = db.QueryRow("SELECT text, gadget, taken from notes").Scan(&n2.Text, &n2.Gadget, &ts)
	n2.Taken = time.Unix(ts, 0)
	if err != nil {
		t.Fatal(err)
	}
	if n2.Text != "me do this" {
		t.Error(n2)
	}

	if n2.Gadget != "lab" {
		t.Error(n2)
	}
	start := time.Unix(ts - 1000, 0)
	end := time.Unix(ts + 1000, 0)
	notes := GetNotes("lab", start, end)
	if len(notes) != 1 {
		t.Fatal(notes)
	}
	n = notes[0]
	if n.Text != "me do this" {
		t.Error(n)
	}

	start = time.Unix(ts + 1000, 0)
	end = time.Unix(ts + 2000, 0)
	notes = GetNotes("lab", start, end)
	if len(notes) != 0 {
		t.Fatal(notes)
	}
	
	os.RemoveAll(tmp)
}
