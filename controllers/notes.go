package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/cswank/gadgetsweb/models"
)

func GetNotes(w http.ResponseWriter, r *http.Request, u *models.User, vars map[string]string) error {
	start, end, err := getNotesStartandEnd(r)
	if err != nil {
		return err
	}
	notes := models.GetNotes(vars["name"], start.UTC(), end.UTC())
	b, err := json.Marshal(notes)
	if err != nil {
		return err
	}
	w.Write(b)
	return err
}

func SaveNote(w http.ResponseWriter, r *http.Request, u *models.User, vars map[string]string) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	note := &models.Note{}
	err = json.Unmarshal(body, note)
	if err != nil {
		return err
	}
	note.Gadget = vars["name"]
	return note.Save()
}

func getNotesStartandEnd(r *http.Request) (time.Time, time.Time, error) {
	params := r.URL.Query()
	var start, end time.Time
	var err error
	var i int64
	if len(params["start"]) > 0 {
		i, err = strconv.ParseInt(params["start"][0], 10, 64)
		start = time.Unix(i, 0)
	} else {
		start = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	}
	if len(params["end"]) > 0 {
		i, err = strconv.ParseInt(params["end"][0], 10, 64)
		end = time.Unix(i, 0)
	} else {
		end = time.Now()
	}
	return start, end, err
}
