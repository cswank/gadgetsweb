package controllers

import (
	"io/ioutil"
	"github.com/gorilla/mux"
	"bitbucket.org/cswank/gadgetsweb/models"
	"encoding/json"
	"net/http"
)

func GetMethods(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	methods, err := models.GetMethods(vars["name"])
	if err != nil {
		return err
	}
	b, err := json.Marshal(methods)
	if err != nil {
		return err
	}
	w.Write(b)
	return err
}

func SaveMethod(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	method := &models.Method{
		Gadget: vars["name"],
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, method)
	if err != nil {
		return err
	}
	return method.Save()
}

