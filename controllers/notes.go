package controllers

import (
	"time"
	"io/ioutil"
	"bitbucket.org/cswank/gadgetsweb/models"
	"encoding/json"
	"net/http"
)

func GetNotes(w http.ResponseWriter, r *http.Request, u *models.User, vars map[string]string) error {
	
	start, end, err := getStartandEnd(r)
	if err != nil {
		return err
	}
	notes := models.GetNotes(vars["name"], start, end)
	b, err := json.Marshal(notes)
	if err != nil {
		return err
	}
	w.Write(b)
	return err
}

func SaveNote(w http.ResponseWriter, r *http.Request, u *models.User, vars map[string]string) error {
	vars := mux.Vars(r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	note := &models.Note{
		Gadget: vars["name"],
	}
	err = json.Unmarshal(body, note)
	if err != nil {
		return err
	}
	return note.Save()
}

