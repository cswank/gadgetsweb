package controllers

import (
	"fmt"
	"io/ioutil"
	"bitbucket.com/cswank/gadgetsweb/models"
	"encoding/json"
	"net/http"
)

func GetMethods(w http.ResponseWriter, r *http.Request) error {
	methods, err := models.GetMethods()
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
	method := &models.Method{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	err = json.Unmarshal(body, method)
	if err != nil {
		return err
	}
	return method.Save()
}

