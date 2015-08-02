package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/cswank/gadgetsweb/models"
	"github.com/cswank/gogadgets"
)

func GetStatus(w http.ResponseWriter, r *http.Request, u *models.User, vars map[string]string) error {
	host, ok := vars["name"]
	if !ok {
		return errors.New("you must supply a host arg")
	}
	cfg := gogadgets.SocketsConfig{
		Host:   host,
		Master: false,
	}
	s := gogadgets.NewClientSockets(cfg)
	err := s.Connect()
	if err != nil {
		return err
	}
	time.Sleep(200 * time.Millisecond)
	defer s.Close()
	status, err := s.SendStatusRequest()
	if err == nil {
		d, _ := json.Marshal(status)
		w.Write(d)
	}
	return err
}

func SendCommand(w http.ResponseWriter, r *http.Request, u *models.User, vars map[string]string) error {
	host, ok := vars["name"]
	if !ok {
		return errors.New("you must supply a host arg")
	}
	cfg := gogadgets.SocketsConfig{
		Host:   host,
		Master: false,
	}
	s := gogadgets.NewClientSockets(cfg)
	err := s.Connect()
	if err != nil {
		return err
	}
	time.Sleep(200 * time.Millisecond)
	defer s.Close()

	d := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var cmd gogadgets.Message
	if err = d.Decode(&cmd); err != nil {
		return err
	}
	s.SendMessage(cmd)
	time.Sleep(200 * time.Millisecond)
	return nil
}

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

// func AddGadgets(w http.ResponseWriter, r *http.Request, u *models.User, vars map[string]string) error {
// 	d := json.NewDecoder(r.Body)
// 	defer r.Body.Close()
// 	var cfg gogadgets.Config
// 	err := d.Decode(&cfg)
// 	if err != nil {
// 		return err
// 	}
// 	s := gogadgets.NewClientSockets(cfg)
// 	err := s.Connect()
// 	if err != nil {
// 		return err
// 	}
// 	time.Sleep(200 * time.Millisecond)
// 	defer s.Close()
// 	cfg.Host = "localhost"
// 	msg := gogadgets.Message{
// 		Config: cfg,
// 	}
// 	s.SendMessage(msg)
// 	time.Sleep(400 * time.Millisecond)
// 	return nil
// }
