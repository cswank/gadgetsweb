package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"

	"github.com/cswank/gadgetsweb/models"
)

var (
	testStartTime = time.Date(2014, 0, 0, 0, 0, 0, 0, time.UTC)
	testEndTime   = time.Date(2014, 0, 0, 2, 0, 0, 0, time.UTC)
	testTime      = time.Date(2014, 0, 0, 1, 0, 0, 0, time.UTC)
)

func cleanup(tmp string) {
	if models.DB != nil {
		models.DB.Close()
	}
	os.RemoveAll(tmp)
}

func saveTestNote() string {
	tmp, _ := ioutil.TempDir("", "")
	os.Setenv("GADGETSDB", path.Join(tmp, "db"))
	n := models.Note{
		Text:   "me do this",
		Gadget: "lab",
		Taken:  testTime,
	}
	n.Save()
	return tmp
}

func TestGetNotes(t *testing.T) {
	tmp := saveTestNote()
	defer cleanup(tmp)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", fmt.Sprintf("http://gadgetsweb/api/lab/notes?start=%d&end=%d", testStartTime.Unix(), testEndTime.Unix()), nil)
	vars := map[string]string{"name": "lab"}
	GetNotes(w, r, nil, vars)
	d := w.Body.Bytes()
	notes := []models.Note{}
	err := json.Unmarshal(d, &notes)
	if err != nil {
		t.Fatal(err)
	}
	if len(notes) != 1 {
		t.Fatal(notes)
	}
	n := notes[0]
	if n.Text != "me do this" {
		t.Fatal(n)
	}

	vars = map[string]string{"name": "overthere"}
	GetNotes(w, r, nil, vars)
	d = w.Body.Bytes()
	notes = []models.Note{}
	json.Unmarshal(d, &notes)
	if len(notes) != 0 {
		t.Fatal(notes)
	}
}

func TestSaveNote(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	os.Setenv("GADGETSDB", path.Join(tmp, "db"))
	defer cleanup(nil, tmp)

	n := &models.Note{
		Text:  "hiya",
		Taken: testTime,
	}
	d, _ := json.Marshal(n)
	buf := bytes.NewBuffer(d)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "http://gadgetsweb/api/lab/notes", buf)
	vars := map[string]string{"name": "lab"}

	err := SaveNote(w, r, nil, vars)

	if err != nil {
		t.Fatal(err)
	}

	notes := models.GetNotes("lab", testStartTime, testEndTime)

	if len(notes) != 1 {
		t.Fatal(notes)
	}
	n2 := notes[0]
	if n2.Text != "hiya" {
		t.Fatal(n2)
	}
}
