package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cswank/gadgetsweb/models"
)

var (
	mongoHost = "localhost"
)

func init() {
	h := os.Getenv("MONGOHOST")
	if len(h) > 0 {
		mongoHost = h
	}
}

func GetDevices(w http.ResponseWriter, r *http.Request, u *models.User, vars map[string]string) error {
	hq, err := getQuery(r, vars)
	if err != nil {
		return err
	}
	d, err := models.GetDevices(hq)
	if err != nil {
		return err
	}
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	w.Write(b)
	return nil
}

func GetTimeseries(w http.ResponseWriter, r *http.Request, u *models.User, vars map[string]string) error {
	hq, err := getQuery(r, vars)
	if err != nil {
		return err
	}
	h, err := models.GetHistory(hq)
	if err != nil {
		return err
	}
	b, err := json.Marshal(h)
	if err != nil {
		return err
	}
	w.Write(b)
	return nil
}

func getQuery(r *http.Request, vars map[string]string) (hq *models.HistoryQuery, err error) {
	start, end, err := getStartandEnd(r)
	if err != nil {
		return hq, err
	}
	hq = &models.HistoryQuery{
		Host:       mongoHost,
		DBName:     vars["gadget"],
		Collection: "updates",
		Location:   vars["location"],
		Name:       vars["device"],
		Start:      start,
		End:        end,
	}
	return hq, err
}

func getStartandEnd(r *http.Request) (time.Time, time.Time, error) {
	params := r.URL.Query()
	var startStr, endStr string
	if len(params["start"]) > 0 {
		startStr = params["start"][0]
		endStr = params["end"][0]
	} else {
		startStr = "0"
		endStr = "0"
	}
	start, err := strconv.ParseInt(startStr, 10, 64)
	end, err := strconv.ParseInt(endStr, 10, 64)
	return time.Unix(start, 0), time.Unix(end, 0), err
}
