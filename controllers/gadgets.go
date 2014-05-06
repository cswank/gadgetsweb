package controllers

import (
	"bitbucket.org/cswank/gadgetsweb/models"
	"bitbucket.org/cswank/gogadgets"
	"encoding/json"
	"net/http"
)

func GetGadgets(w http.ResponseWriter, r *http.Request) error {
	gadgets := models.GetGadgets()
	b, err := json.Marshal(map[string][]models.Gadget{"gadgets": gadgets})
	if err != nil {
		return err
	}
	w.Write(b)
	return nil
}

func GetGadgetTypes(w http.ResponseWriter, r *http.Request) error {
	types := gogadgets.GetTypes()
	d, _ := json.Marshal(types)
	w.Write(d)
	return nil
}

func AddGadgets(w http.ResponseWriter, r *http.Request) error {
	return nil
}
