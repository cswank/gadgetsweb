package models

import (
	"fmt"
	"time"
	"labix.org/v2/mgo"
	"net/url"
        "labix.org/v2/mgo/bson"
	"bitbucket.org/cswank/gogadgets/models"	
)


type HistoryQuery struct {
	Host string
	DBName string
	Collection string
	Location string
	Name string
	Start time.Time
	End time.Time
}

type Series struct {
	Name string `json:"name"`
	Data []interface{} `json:"data"`
}

type Link struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type Devices struct {
	Links []Link `json:"links"`
}

func NewLink(gadget, location, name string) Link {
	u, _ := url.Parse(fmt.Sprintf("/api/history/gadgets/%s/locations/%s/devices/%s", gadget, location, name))
	return Link{
		Name: fmt.Sprintf("%s %s", location, name),
		Path: u.String(),
	}
}

func GetDevices(hq *HistoryQuery) (*Devices, error) {
	d := &Devices{}
	c, session, err := getCollection(hq)
	defer session.Close()
	if err != nil {
		return d, err
	}
	var locations []string
	c.Find(bson.M{}).Distinct("location", &locations)
	links := []Link{}
	for _, l := range locations {
		var devices []string
		c.Find(bson.M{"location": l}).Distinct("name", &devices)
		for _, d := range devices {			
			links = append(links, NewLink(hq.DBName, l, d))
		}
	}
	d.Links = links
	return d, nil
	
}

func GetHistory(hq *HistoryQuery) (Series, error) {
	c, session, err := getCollection(hq)
	defer session.Close()
	if err != nil {
		return Series{}, err
	}
	var results []models.Message
	err = c.Find(
		bson.M{
			"location": hq.Location,
			"name": hq.Name,
			"timestamp": bson.M{
				"$gte": hq.Start,
				"$lte": hq.End,
			},
		},
	).Sort("timestamp").All(&results)
	s := Series{
		Name: fmt.Sprintf("%s %s", hq.Location, hq.Name),
		Data: make([]interface{}, len(results)),
	}
	for i, r := range results {
		f, ok := r.Value.ToFloat()
		if ! ok {
			return Series{}, err
		}
		s.Data[i] = []interface{}{r.Timestamp.Unix() * 1000, f}
	}
	return s, nil
}


func getCollection(hq *HistoryQuery) (*mgo.Collection, *mgo.Session, error) {
	session, err := mgo.Dial(hq.Host)
	c := &mgo.Collection{}
	if err != nil {
		return c, session, err
        }
	c = session.DB(hq.DBName).C(hq.Collection)
	return c, session, nil
}
