package controllers

import (
	"os"
	"time"
	"testing"
	"net/http"
	"net/http/httptest"
	"database/sql"
)

var (
	testStartTime = time.Time{2014, 0, 0, 0, 0, 0, 0, nil}.Unix()
	testEndTime = time.Time{2014, 0, 0, 2, 0, 0, 0, nil}.Unix()
	testTime = time.Time{2014, 0, 0, 1, 0, 0, 0, nil}
)

func cleanup(db *sql.DB, tmp string) {
	db.Close()
	os.RemoveAll(tmp)
}

func saveTestNote() (*sql.DB, string) {
	tmp, _ := ioutil.TempDir("", "")
	os.Setenv("GADGETSDB", path.Join(tmp, "db"))
	db, _ := getDB()
	n := Note{
		Text: "me do this",
		Gadget: "lab",
		Taken: testTime,
	}
	n.Save()
	return db, tmp
}

func TestGetNotes(t *testing.T) {
	db := saveTestNote()
	defer cleanup(db, tmp)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", fmt.Sprintf("http://gadgetsweb/api/notes?start=%d&end=%d", testStartTime, testEndTime), nil)
	
}

func TestSaveNote(t *testing.T) {
	
}






