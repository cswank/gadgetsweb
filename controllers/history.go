package controllers

import (
	"time"
	"strconv"
	"github.com/gorilla/mux"
	"bitbucket.com/cswank/gadgetsweb/models"
	"net/http"
	"encoding/json"
)


func GetTimeseries(w http.ResponseWriter, r *http.Request) error {
	hq, err := getQuery(r)
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

func getQuery(r *http.Request) (hq *models.HistoryQuery, err error) {
	vars := mux.Vars(r)
	start, end, err := getStartandEnd(r)
	if err != nil {
		return hq, err
	}
	hq = &models.HistoryQuery{
		Host: "localhost",
		DBName: "brewery",
		Collection: "updates",
		Location: vars["location"],
		Name: vars["device"],
		start: start,
		end: end,
	}
	return hq, err
}

func getStartandEnd(r *http.Request) (time.Time, time.Time, error) {
	params := r.URL.Query()
	startStr := params["start"][0]
	endStr := params["end"][0]
	start, err := strconv.ParseInt(startStr, 10, 64)
	end, err := strconv.ParseInt(endStr, 10, 64)
	return time.Unix(start, 0), time.Unix(end, 0), err
}

