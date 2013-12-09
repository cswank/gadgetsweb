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

func AddMethod(w http.ResponseWriter, r *http.Request) error {
	method := &models.Method{}
	body, err := ioutil.ReadAll(r.Body)
	fmt.Println(string(body))
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, method)
	if err != nil {
		return err
	}
	mmethod.Steps = 
	return method.Save()
}

func UpdateMethod(w http.ResponseWriter, r *http.Request) error {
	return nil
}
