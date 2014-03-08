package controllers

import (
	"log"
	"time"
	"fmt"
	"strconv"
	"github.com/gorilla/mux"
	"labix.org/v2/mgo"
        "labix.org/v2/mgo/bson"
	"bitbucket.org/cswank/gadgetsweb/models"
	"net/http"
	"encoding/json"
)


type timeseriesVars struct {
	location string
	direction string
	device string
	start time.Time
	end time.Time
}

type Device struct {
	Units string      `json:"units"`
	Value interface{} `json:"value"`
	ID    string      `json:"id"`
}

type Location struct {
	Input  map[string]Device `json:"input"`
	Output map[string]Device `json:"output"`
}

type Message struct {
	Sender      string      `json:"sender"`
	Type        string      `json:"type"`
	Body        string      `json:"body"`
	Timestamp   time.Time   `json:"timestamp"`
	Name        string      `json:"name"`
	Locations   map[string]Location `json:"locations"`
}


func GetSummary(w http.ResponseWriter, r *http.Request) error {
	summary := []models.Summary{}
	c, session, err := getCollection("updates")
	defer session.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	iter := c.Find(nil).Sort("-timestamp").Limit(1).Iter()
	gadgets := &Message{}
	iter.Next(gadgets)
	for location, value := range(gadgets.Locations) {
		for key, _ := range value.Input {
			summary = append(summary, models.Summary{Name:key, Location:location, Direction:"input"})
		}
		for key, _ := range value.Output {
			summary = append(summary, models.Summary{Name:key, Location:location, Direction:"output"})
		}
	}
	b, err := json.Marshal(summary)
	if err != nil {
		return err
	}
	w.Write(b)
	return nil
}

func Btoi(b bool) int {
    if b {
        return 1
    }
    return 0
 }

func GetTimeseries(w http.ResponseWriter, r *http.Request) error {
	//location, direction, device string, start, end *time.Time
	vars, err := getTimeseriesVars(r)
	t, err := getTimeseries(vars)
	if err != nil {
		return err
	}
	b, err := json.Marshal(t)
	if err != nil {
		return err
	}
	w.Write(b)
	return err
}

func getTimeseries(vars * timeseriesVars) (*models.Timeseries, error) {
	t := &models.Timeseries{Name:fmt.Sprintf("%s %s", vars.location, vars.device)}
	gadget := &Message{}
	c, session, err := getCollection("updates")
	defer session.Close()
	if err != nil {
		return t, err
	}
	query := getQuery(vars)
	iter := c.Find(query).Sort("timestamp").Iter()
	for iter.Next(gadget) {
		appendData(t, gadget, vars)
	}
	return t, nil
}


func appendData(t *models.Timeseries, gadget *Message, vars *timeseriesVars) {
	location := gadget.Locations[vars.location]
	device := getDevice(vars.direction, &location, vars)
	switch v := device.Value.(type) {
	default:
		t.Data = append(t.Data, []interface{}{gadget.Timestamp.Unix() * 1000, device.Value})
	case bool:
		i := Btoi(v)
		t.Data = append(t.Data, []interface{}{gadget.Timestamp.Unix() * 1000, i})
	}
}

func getDevice(direction string, location *Location, vars *timeseriesVars) *Device {
	var device Device
	if direction == "input" {
		device = location.Input[vars.device]
	} else {
		device = location.Output[vars.device]
	}
	return &device
}

func getQuery(vars *timeseriesVars) bson.M {
	gt := bson.M{"timestamp": bson.M{"$gte": vars.start}}
	lt := bson.M{"timestamp": bson.M{"$lte": vars.end}}
	return bson.M{"$and": []bson.M{gt, lt}}
}

func getTimeseriesVars(r *http.Request) (tv *timeseriesVars, err error) {
	vars := mux.Vars(r)
	start, end, err := getStartandEnd(r)
	if err != nil {
		return tv, err
	}
	tv = &timeseriesVars{
		location: vars["location"],
		direction: vars["direction"],
		device: vars["device"],
		start: start,
		end: end,
	}
	return tv, err
}

func getStartandEnd(r *http.Request) (time.Time, time.Time, error) {
	params := r.URL.Query()
	startStr := params["start"][0]
	endStr := params["end"][0]
	start, err := strconv.ParseInt(startStr, 10, 64)
	end, err := strconv.ParseInt(endStr, 10, 64)
	return time.Unix(start / 1000, 0), time.Unix(end / 1000, 0), err
}

func getCollection(collection string) (*mgo.Collection, *mgo.Session, error) {
	session, err := mgo.Dial("localhost")
	c := &mgo.Collection{}
	if err != nil {
		log.Println("couldn't open mongo connection", err)
		return c, session, err
        }
	c = session.DB("greenhouse").C(collection)
	return c, session, nil
}
