package controllers

import (
	"time"
	"bitbucket.org/cswank/gadgetsweb/models"
	"bitbucket.org/cswank/gogadgets"
	"encoding/json"
	"net/http"
)

func GetGadgets(w http.ResponseWriter, r *http.Request, u *models.User, vars map[string]string) error {
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

func GetGadgetTypes(w http.ResponseWriter, r *http.Request, u *models.User, vars map[string]string) error {
	types := gogadgets.GetTypes()
	d, _ := json.Marshal(types)
	w.Write(d)
	return nil
}

func AddGadgets(w http.ResponseWriter, r *http.Request, u *models.User, vars map[string]string) error {
	d := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var cfg gogadgets.Config
	err := d.Decode(&cfg)
	if err != nil {
		return err
	}
	s, err := gogadgets.NewClientSockets(cfg.Host)
	if err != nil {
		return err
	}
	time.Sleep(200 * time.Millisecond)
	defer s.Close()
	cfg.Host = "localhost"
	msg := gogadgets.Message{
		Config: cfg,
	}
	s.SendMessage(msg)
	time.Sleep(400 * time.Millisecond)
	return nil
}












