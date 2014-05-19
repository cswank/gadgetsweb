package controllers

import (
	"strconv"
	"io/ioutil"
	"github.com/gorilla/mux"
	"bitbucket.org/cswank/gadgetsweb/models"
	"encoding/json"
	"net/http"
)

func GetNotes(w http.ResponseWriter, r *http.Request, u *models.User) error {
	vars := mux.Vars(r)
	notes := models.GetNotes(vars["name"])
	b, err := json.Marshal(methods)
	if err != nil {
		return err
	}
	w.Write(b)
	return err
}

func SaveNote(w http.ResponseWriter, r *http.Request, u *models.User) error {
	vars := mux.Vars(r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	note := &models.Note{
		Name: vars["name"],
		Taken: time.Now(),
	}
	err = json.Unmarshal(body, notes)
	if err != nil {
		return err
	}
	return notes.Save()
}

