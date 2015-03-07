package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/cswank/gadgetsweb/models"
)

func GetMethods(w http.ResponseWriter, r *http.Request, u *models.User, vars map[string]string) error {
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

func SaveMethod(w http.ResponseWriter, r *http.Request, u *models.User, vars map[string]string) error {
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

func DeleteMethod(w http.ResponseWriter, r *http.Request, u *models.User, vars map[string]string) error {
	id, err := strconv.ParseUint(vars["methodId"], 10, 64)
	if err != nil {
		return err
	}
	method := &models.Method{
		Id:     id,
		Gadget: vars["name"],
	}
	return method.Delete()
}
