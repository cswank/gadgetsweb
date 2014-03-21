package controllers

import (
	"bitbucket.org/cswank/gadgetsweb/models"
	"encoding/json"
	"net/http"
)

func GetGadgets(w http.ResponseWriter, r *http.Request) error {
	gadgets, err := models.GetGadgets()
	if err != nil {
		return err
	}
	b, err := json.Marshal(gadgets)
	if err != nil {
		return err
	}
	w.Write(b)
	return nil
}
